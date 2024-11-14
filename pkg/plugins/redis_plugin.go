package plugins

import (
	"encoding/json"
	"fmt"

	"github.com/Kong/go-pdk"
	"github.com/go-redis/redis"
)

type RedisCachePlugin struct {
	RedisAddr string `json:"redis_addr"`
	RedisPssw string `json:"redis_pssw"`
	RedisUser string `json:"redis_user"`

	DB            int    `json:"redis_db"`
	RedisCacheKey string `json:"redis_cache_key"`
}

func (rcp RedisCachePlugin) Access(kong *pdk.PDK) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     rcp.RedisAddr,
		Password: rcp.RedisPssw, // no password set
		DB:       rcp.DB,        // use default DB
	})

	result, err := rdb.Get(rcp.RedisCacheKey).Result()

	if err != nil && err != redis.Nil {
		rcp.handleError(fmt.Sprintf("Erro ao conectar com o redis: %v", err), kong)
		return
	}

	if result != "" {

		var jsonResult []map[string]any

		json.Unmarshal([]byte(result), &jsonResult)

		response, _ := json.Marshal(jsonResult)

		kong.Response.SetHeader("X-Cache-Redis-Hit", "true")
		kong.Response.Exit(200, response, nil)

		return
	}

}

func (rcp RedisCachePlugin) Response(kong *pdk.PDK) {
	fmt.Println(rcp.RedisPssw)

	rdb := redis.NewClient(&redis.Options{
		Addr:     rcp.RedisAddr,
		Password: rcp.RedisPssw, // no password set
		DB:       rcp.DB,        // use default DB
	})

	responseKong, _ := kong.ServiceResponse.GetRawBody()

	_, err := rdb.Get(rcp.RedisCacheKey).Result()

	if err == redis.Nil {
		err := rdb.Set(rcp.RedisCacheKey, string(responseKong), 0).Err()

		if err != nil {
			rcp.handleError(fmt.Sprintf("Erro ao conectar com o redis: %v", err), kong)
			return
		}
	}

}

func (rcp RedisCachePlugin) handleError(msg string, kong *pdk.PDK) {
	jsonMsg := map[string]string{
		"msg": msg,
	}

	response, _ := json.Marshal(&jsonMsg)

	kong.Response.Exit(500, response, nil)
}

func New() any {
	return &RedisCachePlugin{}
}
