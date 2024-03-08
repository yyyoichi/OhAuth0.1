package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

type options struct {
	body  io.Reader
	query *url.Values
}
type Option func(options *options) error

func WithBody(body io.Reader) Option {
	return func(options *options) error {
		options.body = body
		return nil
	}
}

func WithQuery(query url.Values) Option {
	return func(options *options) error {
		options.query = &query
		return nil
	}
}

type Config struct {
	Router *gin.Engine
	Method string
	Path   string
}

func Serve(t *testing.T, config Config, opts ...Option) (*http.Request, *httptest.ResponseRecorder) {
	var body io.Reader = nil
	var path = config.Path

	var options options
	for _, opt := range opts {
		_ = opt(&options)
	}
	if options.body != nil {
		body = options.body
	}
	if options.query != nil {
		path += "?" + options.query.Encode()
	}
	req := httptest.NewRequest(config.Method, path, body)
	w := httptest.NewRecorder()
	config.Router.ServeHTTP(w, req)
	return req, w
}
