package resource

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
	server_test "github.com/yyyoichi/OhAuth0.1/internal/test"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHandlerOK(t *testing.T) {
	db, _ := database.NewDatabase()
	service := &Service{
		client: db,
	}
	headerOption := func() server_test.Option {
		accesstoken := "token"
		err := db.CreateAccessToken(context.Background(), &apiv1.AccessToken{
			Token:           accesstoken,
			UserId:          "1",
			ServiceClientId: "501",
			Expires:         timestamppb.New(time.Now().AddDate(0, 0, 1)),
			Scope:           "profile:view",
		})
		assert.NoError(t, err)
		return server_test.WithHeader("Authorization", "Bearer "+accesstoken)
	}()
	router := SetupRouter(service)
	test := map[string]struct {
		config  server_test.Config
		options []server_test.Option
		body    interface{}
	}{
		"GET:/status": {
			config: server_test.Config{
				Router: router,
				Method: http.MethodGet,
				Path:   "/api/v1/profile",
			},
			options: []server_test.Option{headerOption},
			body:    struct{}{},
		},
		"GET:/profile": {
			config: server_test.Config{
				Router: router,
				Method: http.MethodGet,
				Path:   "/api/v1/profile",
			},
			options: []server_test.Option{headerOption},
			body:    ProfileGetResponse{},
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

func TestHandlerAuthorizationStatus(t *testing.T) {
	db, _ := database.NewDatabase()
	service := &Service{
		client: db,
	}
	router := SetupRouter(service)
	test := map[string]struct {
		options []server_test.Option
		expCode int
	}{
		"ok": {
			options: []server_test.Option{
				func() server_test.Option {
					accesstoken := "token"
					err := db.CreateAccessToken(context.Background(), &apiv1.AccessToken{
						Token:           accesstoken,
						UserId:          "1",
						ServiceClientId: "501",
						Expires:         timestamppb.New(time.Now().AddDate(0, 0, 1)),
						Scope:           "profile:view",
					})
					assert.NoError(t, err)
					return server_test.WithHeader("Authorization", "Bearer "+accesstoken)
				}()},
			expCode: http.StatusNoContent,
		},
		"empty": {
			options: []server_test.Option{},
			expCode: http.StatusForbidden,
		},
		"Bearer empty": {
			options: []server_test.Option{server_test.WithHeader("Authorization", "Bearer ")},
			expCode: http.StatusForbidden,
		},
		"invalid": {
			options: []server_test.Option{server_test.WithHeader("Authorization", "Bearer hogehoge")},
			expCode: http.StatusForbidden,
		},
		"expired": {
			options: []server_test.Option{
				func() server_test.Option {
					accesstoken := "expired-token"
					err := db.CreateAccessToken(context.Background(), &apiv1.AccessToken{
						Token:           accesstoken,
						UserId:          "1",
						ServiceClientId: "501",
						Expires:         timestamppb.New(time.Now().AddDate(0, 0, -1)), // !
						Scope:           "profile:view",
					})
					assert.NoError(t, err)
					return server_test.WithHeader("Authorization", "Bearer "+accesstoken)
				}()},
			expCode: http.StatusBadRequest,
		},
	}
	for scenario, tt := range test {
		t.Run(scenario, func(t *testing.T) {
			config := server_test.Config{
				Router: router,
				Method: http.MethodGet,
				Path:   "/api/v1/status",
			}
			_, resp := server_test.Serve(t, config, tt.options...)
			assert.Equal(t, tt.expCode, resp.Code)
		})
	}
}

func TestProfileHandler(t *testing.T) {
	db, _ := database.NewDatabase()
	service := &Service{
		client: db,
	}
	router := SetupRouter(service)
	test := map[string]struct {
		options []server_test.Option
		expCode int
	}{
		"ok": {
			options: []server_test.Option{func() server_test.Option {
				accesstoken := "token"
				err := db.CreateAccessToken(context.Background(), &apiv1.AccessToken{
					Token:           accesstoken,
					UserId:          "1",
					ServiceClientId: "501",
					Expires:         timestamppb.New(time.Now().AddDate(0, 0, 1)),
					Scope:           "profile:view",
				})
				assert.NoError(t, err)
				return server_test.WithHeader("Authorization", "Bearer "+accesstoken)
			}()},
			expCode: http.StatusOK,
		},
		// Maybe this test should be in the 'scope' package, for example.
		"has no-scope": {
			options: []server_test.Option{func() server_test.Option {
				accesstoken := "limitedtoken"
				err := db.CreateAccessToken(context.Background(), &apiv1.AccessToken{
					Token:           accesstoken,
					UserId:          "1",
					ServiceClientId: "501",
					Expires:         timestamppb.New(time.Now().AddDate(0, 0, 1)),
					Scope:           "unknwon",
				})
				assert.NoError(t, err)
				return server_test.WithHeader("Authorization", "Bearer "+accesstoken)
			}()},
			expCode: http.StatusBadRequest,
		},
	}
	for scenario, tt := range test {
		t.Run(scenario, func(t *testing.T) {
			config := server_test.Config{
				Router: router,
				Method: http.MethodGet,
				Path:   "/api/v1/profile",
			}
			_, resp := server_test.Serve(t, config, tt.options...)
			assert.Equalf(t, tt.expCode, resp.Code, resp.Body.String())
		})
	}
}
