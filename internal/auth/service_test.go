package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService(t *testing.T) {
	test := map[string]struct {
		tservice *Service
	}{
		"local": {func() *Service {
			db, _ := database.NewDatabase()
			return &Service{
				client: db,
			}
		}()},
		"remote": {func() *Service {
			port := "3336"
			err := database.NewDatabaseServer(database.ServerConfig{
				Port: port,
			})
			assert.NoError(t, err)
			ctx := context.Background()
			var service *Service
			connectedCh := make(chan interface{})
			go func() {
				defer close(connectedCh)
				for {
					time.Sleep(time.Duration(10) * time.Millisecond)
					var err error
					service, err = NewService(ctx, Config{
						DatabaseServerURL: "http://localhost:" + port,
					})
					if err != nil {
						continue
					}
					if _, err := service.client.GetUserByID(context.Background(), "1"); err != nil {
						continue
					}
					return
				}
			}()
			<-connectedCh
			return service
		}()},
	}
	for scenario, tt := range test {
		t.Run(scenario, func(t *testing.T) {
			t.Run("JWT", func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				tservice := tt.tservice
				claims, err := tservice.Authentication(ctx, "1", "password")
				assert.NoError(t, err)
				claims.ClientId = "hoge"
				assert.NoError(t, err)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				ss, err := token.SignedString([]byte("secret"))
				assert.NoError(t, err)

				// parse
				act, err := tservice.ParseMyClaims(ctx, ss, []byte("secret"))
				assert.NoError(t, err)
				sub, err := act.GetSubject()
				assert.NoError(t, err)
				assert.Equal(t, "1", sub)
				assert.Equal(t, "hoge", act.ClientId)
			})
			t.Run("code to accesstoken refresh", func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				tservice := tt.tservice
				const (
					USER_ID   = "TESTING_USER"
					CLIENT_ID = "TESTING_CLIENT"
				)
				code, err := tservice.NewAuthorizationCode(ctx, NewAuthorizationCodeConfig{
					UserId:          USER_ID,
					ServiceClientId: CLIENT_ID,
				})
				assert.NoError(t, err)
				assert.Equal(t, USER_ID, code.UserId)
				assert.Equal(t, CLIENT_ID, code.ServiceClientId)
				assert.NotEmpty(t, code.Code)
				assert.False(t, code.Expires.AsTime().IsZero())

				testTokens := func(token *apiv1.AccessToken, refresh *apiv1.RefreshToken) {
					assert.NoError(t, err)
					assert.Equal(t, USER_ID, token.UserId)
					assert.Equal(t, CLIENT_ID, token.ServiceClientId)
					assert.NotEmpty(t, token.Token)
					assert.False(t, code.Expires.AsTime().IsZero())

					assert.Equal(t, USER_ID, refresh.UserId)
					assert.Equal(t, CLIENT_ID, refresh.ServiceClientId)
					assert.NotEmpty(t, refresh.Token)
					assert.False(t, refresh.Expires.AsTime().IsZero())
				}
				token, refresh, err := tservice.NewAccessToken(ctx, code.Code)
				testTokens(token, refresh)
				token, refresh, err = tservice.UpdateAccessToken(ctx, refresh.Token)
				testTokens(token, refresh)
			})

			t.Run("expired", func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				tservice := tt.tservice

				err := tservice.client.CreateAuthorizationCode(ctx, &apiv1.AuthorizationCode{
					Code:    "example",
					Expires: timestamppb.New(time.Now().Add(time.Duration(-1) * time.Minute)),
				})
				assert.NoError(t, err)

				_, _, err = tservice.NewAccessToken(ctx, "example")
				assert.ErrorIs(t, ErrAuthorizationCodeExpired, err)

				err = tservice.client.CreateRefreshToken(ctx, &apiv1.RefreshToken{
					Token:   "example",
					Expires: timestamppb.New(time.Now().Add(time.Duration(-1) * time.Minute)),
				})
				assert.NoError(t, err)
				_, _, err = tservice.UpdateAccessToken(ctx, "example")
				assert.ErrorIs(t, ErrRefreshTokenExpired, err)

			})
		})
	}
}
