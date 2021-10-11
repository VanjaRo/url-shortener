package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"url-shortener/storage"

	routerClient "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"shortUrl"`
}

type handler struct {
	schema  string
	host    string
	storage storage.Service
}

func New(schema, host string, storage storage.Service) *routerClient.Router {
	router := routerClient.New()

	h := handler{
		schema,
		host,
		storage,
	}

	router.POST("/encode/", responseHandler(h.encode))
	router.GET("/{shortLink}", h.redirect)
	router.GET("/{shortLink}/info", responseHandler(h.decode))
	return router
}

func responseHandler(h func(ctx *fasthttp.RequestCtx) (interface{}, int, error)) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		data, status, err := h(ctx)
		if err != nil {
			data = err.Error()
		}
		ctx.Response.Header.Set("Context-Type", "application/json")
		ctx.Response.SetStatusCode(status)
		err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response{
			Data:    data,
			Success: err == nil,
		})
		if err != nil {
			log.Printf("could not encode response to output: %v", err)
		}
	}
}

func (h handler) encode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var input struct {
		URL     string `json:"url"`
		Expires string `json:"expires"`
	}

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unable to decode JSON request body: %v", err)
	}

	uri, err := url.ParseRequestURI(input.URL)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid url")
	}

	layoutISO := "2006-01-02 15:04:05"
	expires, err := time.Parse(layoutISO, input.Expires)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid expiration date")
	}
	// expires = "2022-01-01 11:11:11"

	c, err := h.storage.Save(uri.String(), expires)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not store in database: %v", err)
	}

	u := url.URL{
		Scheme: h.schema,
		Host:   h.host,
		Path:   c,
	}
	fmt.Printf("Generated link: %v \n", u.String())

	return u.String(), http.StatusCreated, nil
}

func (h handler) decode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	code := ctx.UserValue("shortLink").(string)

	model, err := h.storage.LoadInfo(code)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("URL not found")
	}

	return model, http.StatusOK, nil
}

func (h handler) redirect(ctx *fasthttp.RequestCtx) {
	code := ctx.UserValue("shortLink").(string)

	uri, err := h.storage.Load(code)
	if err != nil {
		ctx.Response.Header.Set("Context-Type", "application/json")
		ctx.Response.SetStatusCode(http.StatusNotFound)
		return
	}
	ctx.Redirect(uri, http.StatusTemporaryRedirect)
}
