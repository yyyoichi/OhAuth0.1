package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	server_test "github.com/yyyoichi/OhAuth0.1/internal/test"
)

func TestHandlerOK(t *testing.T) {
	service := NewService(Config{})
	router := SetupRouter(service)
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
		"POST:/authorization": {
			config: server_test.Config{
				Router: router,
				Method: http.MethodPost,
				Path:   "/api/v1/authorization",
			},
			options: []server_test.Option{
				server_test.WithBody(func() io.Reader {
					claims, err := service.Authentication(context.Background(), "1", "password")
					assert.NoError(t, err)
					claims.ClientID = "501" // !
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					ss, err := token.SignedString(JWT_SECRET)
					assert.NoError(t, err)
					var req AuthorizationRequest
					req.JWT = ss
					req.ClientID = "501"
					req.ResponseType = "code"
					req.Scope = "profile:view"
					b, err := json.Marshal(req)
					assert.NoError(t, err)
					return bytes.NewBuffer(b)
				}()),
			},
			body: AuthorizationResponse{},
		},
		"POST:/accesstoken?code": {
			config: server_test.Config{
				Router: router,
				Method: http.MethodPost,
				Path:   "/api/v1/accesstoken",
			},
			options: []server_test.Option{
				server_test.WithBody(func() io.Reader {
					authorization, err := service.NewAuthorizationCode(context.Background(), NewAuthorizationCodeConfig{
						UserID:          "1",
						ServiceClientID: "501",
					})
					assert.NoError(t, err)
					var req AccessTokenRequest
					req.ClientID = authorization.ServiceClientID
					req.ClientSecret = "secret"
					req.Code = authorization.Code // !
					req.GrantType = "authorization_code"
					b, err := json.Marshal(req)
					assert.NoError(t, err)
					return bytes.NewBuffer(b)
				}()),
			},
			body: AccessTokenResponse{},
		},
		"POST:/accesstoken?refreshtoken": {
			config: server_test.Config{
				Router: router,
				Method: http.MethodPost,
				Path:   "/api/v1/accesstoken",
			},
			options: []server_test.Option{
				server_test.WithBody(func() io.Reader {
					authorization, err := service.NewAuthorizationCode(context.Background(), NewAuthorizationCodeConfig{
						UserID:          "1",
						ServiceClientID: "501",
					})
					assert.NoError(t, err)
					_, refresh, err := service.NewAccessToken(context.Background(), authorization.Code)
					assert.NoError(t, err)
					var req AccessTokenRequest
					req.ClientID = authorization.ServiceClientID
					req.ClientSecret = "secret"
					req.RefreshToken = refresh.Token
					req.GrantType = "authorization_code"
					b, err := json.Marshal(req)
					assert.NoError(t, err)
					return bytes.NewBuffer(b)
				}()),
			},
			body: AccessTokenResponse{},
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
