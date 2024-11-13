package main

import (
	"log"

	"github.com/Israel-Ferreira/redis-cache-plugin/pkg/plugins"
	"github.com/Kong/go-pdk/server"
)

func main() {
	log.Println("Iniciando o Plugin")

	var (
		version  string = "0.0.1"
		priority int    = 1000
	)

	server.StartServer(plugins.New, version, priority)

	if err := recover(); err != nil {
		log.Printf("Erro ao executar o plugin: %T \n", err)
	}
}
