package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	server_test "github.com/yyyoichi/OhAuth0.1/internal/test"
)

func TestHandlerOK(t *testing.T) {
	router := SetupRouter()
	test := map[string]struct {
		config  server_test.Config
		options []server_test.Option
		body    interface{}
	}{
		"GET:/client/:client_id": {
			config: server_test.Config{
				Router: router,
				Method: http.MethodGet,
				Path:   "/api/v1/clients/500",
			},
			options: []server_test.Option{},
			body:    ServiceClientGetResponse{},
		},
		"POST:/authentication": {
			config: server_test.Config{
				Router: router,
				Method: http.MethodPost,
				Path:   "/api/v1/authentication",
			},
			options: []server_test.Option{
				server_test.WithBody(func() io.Reader {
					req := AuthenticationRequest{
						UserID:   "1",
						ClientID: "501",
						Password: "password",
					}
					b, _ := json.Marshal(req)
					return bytes.NewBuffer(b)
				}()),
			},
			body: AuthenticationResponse{},
		},
	}
	for scenario, tt := range test {
		t.Run(scenario, func(t *testing.T) {
			_, resp := server_test.Serve(t, tt.config, tt.options...)
			assert.Equalf(t, http.StatusOK, resp.Code, resp.Body.String())
			err := json.Unmarshal(resp.Body.Bytes(), &tt.body)
			assert.NoError(t, err)
			assert.NotZero(t, tt.body)
		})
	}
}
