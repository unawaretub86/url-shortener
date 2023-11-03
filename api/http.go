package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	errs "github.com/pkg/errors"

	"github.com/unawaretub86/url-shortener/serialize/json"
	msgpack "github.com/unawaretub86/url-shortener/serialize/msgPack"
	"github.com/unawaretub86/url-shortener/shortener"
)

type RedirectHandler interface {
	Get(*gin.Context)
	Post(*gin.Context)
}

type handler struct {
	redirectService shortener.RedirectService
}

func NewHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService: redirectService}
}

func (h *handler) serializer(contentType string) shortener.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &msgpack.Redirect{}
	}

	return &json.Redirect{}
}

func (h *handler) Get(c *gin.Context) {
	code := c.Param("code")

	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errs.Cause(err) == shortener.ErrRedirectNotFound {
			c.JSON(http.StatusCreated, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusMovedPermanently, redirect.URL)
}

func (h *handler) Post(c *gin.Context) {
	contentType := c.Request.Header.Get(c.ContentType())

	req, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	redirect, err := h.serializer(contentType).Decode(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = h.redirectService.Store(redirect)
	if err != nil {
		if errs.Cause(err) == shortener.ErrRedirectInvalid {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	responseBody, err := h.serializer(contentType).Encode(redirect)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusCreated, contentType, []byte(responseBody))
}
