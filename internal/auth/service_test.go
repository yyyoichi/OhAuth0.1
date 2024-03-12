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
	newLocalService := func() *Service {
		db, _ := database.NewDatabase()
		return &Service{
			client: db,
		}
	}
	t.Run("Authentication", func(t *testing.T) {
		test := []struct {
			id, password string
			expErr       error
		}{
			{"1", "password", nil},
			{"1", "invalidpass", ErrNoMatchPassword},
			{"99", "password", database.ErrNotFound},
		}
		ctx := context.Background()
		for _, tt := range test {
			tservice := newLocalService()
			_, err := tservice.Authentication(ctx, tt.id, tt.password)
			assert.ErrorIs(t, err, tt.expErr)
		}
	})
	t.Run("ParseMyClaims", func(t *testing.T) {
		token := func() string {
			ctx := context.Background()
			tservice := newLocalService()
			claims, err := tservice.Authentication(ctx, "1", "password")
			assert.NoError(t, err)
			claims.ClientId = "hoge"
			assert.NoError(t, err)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			ss, err := token.SignedString(JWT_SECRET)
			assert.NoError(t, err)
			return ss
		}()
		test := []struct {
			ss     string
			secret []byte
			isNil  bool
		}{
			{token, JWT_SECRET, true},
			{token, []byte("secret"), false},
		}
		ctx := context.Background()
		for _, tt := range test {
			tservice := newLocalService()
			_, err := tservice.ParseMyClaims(ctx, tt.ss, tt.secret)
			if tt.isNil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		}
	})
	t.Run("JWT", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tservice := newLocalService()
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
		tservice := newLocalService()
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
		tservice := newLocalService()

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
}
