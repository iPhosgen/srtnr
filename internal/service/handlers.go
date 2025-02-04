package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/iphosgen/srtnr/internal/storage"
	"github.com/iphosgen/srtnr/pkg/shortener"
	"github.com/valyala/fasthttp"
)

const (
	minUrlLen = 11
	userIdLen = 36
)

var (
	invalidUrlError    = errors.New("given URL is too short")
	invalidUserIdError = errors.New("given user ID is incorrect")
)

type ShortenDTO struct {
	Url   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

type Handler struct {
	shortener shortener.Shortener
	storage   storage.Storage
}

func NewHandler(shortener shortener.Shortener, storage storage.Storage) *Handler {
	return &Handler{shortener: shortener, storage: storage}
}

func (h *Handler) Health(ctx *fasthttp.RequestCtx) {
}

func (h *Handler) Readiness(ctx *fasthttp.RequestCtx) {
}

func fakeURLResolver(_ string) (string, error) {
	return "https://stackoverflow.com/questions/5885486/get-current-time-as-formatted-string-in-go", nil
}

func (h *Handler) DecodeURL(ctx *fasthttp.RequestCtx) {
	shortened := fmt.Sprintf("%v", ctx.UserValue("shortened"))

	if originalURL, err := h.storage.Lookup(shortened); err != nil {
		writeJSONErrorResponse(ctx, ShortenDTO{Error: "URL not found"}, fasthttp.StatusNotFound)
		return
	} else {
		ctx.SetStatusCode(fasthttp.StatusMovedPermanently)
		ctx.Response.Header.Set("Location", originalURL)
	}
}

func (h *Handler) EncodeURL(ctx *fasthttp.RequestCtx) {
	sr := &ShortenDTO{}
	if err := json.Unmarshal(ctx.PostBody(), sr); err != nil {
		writeJSONErrorResponse(ctx, ShortenDTO{Error: err.Error()}, fasthttp.StatusBadRequest)
		return
	}

	userId := string(ctx.Request.Header.Peek("X-User-Id"))

	if len(sr.Url) < minUrlLen {
		writeJSONErrorResponse(ctx, ShortenDTO{Error: invalidUrlError.Error()}, fasthttp.StatusBadRequest)
		return
	}

	if len(userId) != userIdLen {
		writeJSONErrorResponse(ctx, ShortenDTO{Error: invalidUserIdError.Error()}, fasthttp.StatusBadRequest)
		return
	}

	shortened, err := h.shortener.Shorten(sr.Url, userId)
	if err != nil {
		writeJSONErrorResponse(ctx, ShortenDTO{Error: err.Error()}, fasthttp.StatusInternalServerError)
		return
	}

	if err := h.storage.Save(sr.Url, userId, shortened); err != nil {
		writeJSONErrorResponse(ctx, ShortenDTO{Error: err.Error()}, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	writeJSONResponse(ctx, ShortenDTO{Url: shortened})
}

func writeJSONErrorResponse(ctx *fasthttp.RequestCtx, response ShortenDTO, statusCode int) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Error serializing response: %s", err.Error())
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(responseJSON)
}

func writeJSONResponse(ctx *fasthttp.RequestCtx, response ShortenDTO) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Error serializing response: %s", err.Error())
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetBody(responseJSON)
}
