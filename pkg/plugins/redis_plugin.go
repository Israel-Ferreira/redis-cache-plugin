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

	DB int `json:"redis_db"`
}

func (rcp RedisCachePlugin) Response(kong *pdk.PDK) {
	fmt.Println(rcp.RedisPssw)

	rdb := redis.NewClient(&redis.Options{
		Addr:     rcp.RedisAddr,
		Password: rcp.RedisPssw, // no password set
		DB:       rcp.DB,        // use default DB
	})

	responseKong, _ := kong.ServiceResponse.GetRawBody()

	err := rdb.Set("banks-api-response", string(responseKong), 0).Err()

	if err != nil {
		msg := map[string]string{
			"msg": fmt.Sprintf("Erro ao conectar com o redis: %v", err),
		}

		json_msg, _ := json.Marshal(&msg)

		kong.Response.Exit(500, json_msg, nil)
		return
	}

}

func New() any {
	return &RedisCachePlugin{}
}