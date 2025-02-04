package main

import (
	"fmt"

	"github.com/iphosgen/srtnr/config"
	"github.com/iphosgen/srtnr/internal/service"
	"github.com/iphosgen/srtnr/internal/storage"
	"github.com/iphosgen/srtnr/pkg/shortener"
	"github.com/valyala/fasthttp"
)

func main() {
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		panic(err)
	}

	shortenerService := shortener.NewUrlShortener()
	storage, err := storage.NewPostgresStorage(&cfg.Database)
	if err != nil {
		panic(err)
	}
	router := service.NewRouter(shortenerService, storage)

	server := &fasthttp.Server{
		Handler: router,
	}

	if err := server.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)); err != nil {
		panic(err)
	}
}
