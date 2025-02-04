package service

import (
	"github.com/fasthttp/router"
	"github.com/iphosgen/srtnr/internal/storage"
	"github.com/iphosgen/srtnr/pkg/shortener"
	"github.com/valyala/fasthttp"
)

func NewRouter(shortener shortener.Shortener, storage storage.Storage) fasthttp.RequestHandler {
	r := router.New()
	h := NewHandler(shortener, storage)

	r.POST("/shorten", h.EncodeURL)
	r.GET("/s/{shortened}", h.DecodeURL)

	r.GET("/health", h.Health)
	r.GET("/readiness", h.Readiness)

	return r.Handler
}
