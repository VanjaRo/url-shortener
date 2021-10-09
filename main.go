package main

import (
	"log"
	"url-shortener/config"
	"url-shortener/handler"
	"url-shortener/storage/redis"

	"github.com/valyala/fasthttp"
)

func main() {
	conf, err := config.FromFile("./conf.json")
	if err != nil {
		log.Fatal(err)
	}

	service, err := redis.New(conf.Redis.Host, conf.Redis.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	router := handler.New(conf.Options.Schema, conf.Options.Prefix, service)

	log.Fatal(fasthttp.ListenAndServe(":" + conf.Server.Port, router.Handler))
}