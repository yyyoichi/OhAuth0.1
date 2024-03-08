package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
)

func TestService(t *testing.T) {
	t.Run("JWT", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tservice := NewService(Config{})
		claims, err := tservice.Authentication(ctx, "0", "password")
		claims.ClientID = "hoge"
		assert.NoError(t, err)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString([]byte("secret"))
		assert.NoError(t, err)

		// parse
		act, err := tservice.ParseMyClaims(ctx, ss, []byte("secret"))
		assert.NoError(t, err)
		sub, err := act.GetSubject()
		assert.NoError(t, err)
		assert.Equal(t, "0", sub)
		assert.Equal(t, "hoge", act.ClientID)
	})
	t.Run("code to accesstoken refresh", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tservice := NewService(Config{})
		const (
			USER_ID   = "TESTING_USER"
			CLIENT_ID = "TESTING_CLIENT"
		)
		code, err := tservice.NewAuthorizationCode(ctx, NewAuthorizationCodeConfig{
			UserID:          USER_ID,
			ServiceClientID: CLIENT_ID,
		})
		assert.NoError(t, err)
		assert.Equal(t, USER_ID, code.UserID)
		assert.Equal(t, CLIENT_ID, code.ServiceClientID)
		assert.NotEmpty(t, code.Code)
		assert.False(t, code.Expires.IsZero())

		testTokens := func(token *database.AccessToken, refresh *database.RefreshToken) {
			assert.NoError(t, err)
			assert.Equal(t, USER_ID, token.UserID)
			assert.Equal(t, CLIENT_ID, token.ServiceClientID)
			assert.NotEmpty(t, token.Token)
			assert.False(t, code.Expires.IsZero())

			assert.Equal(t, USER_ID, refresh.UserID)
			assert.Equal(t, CLIENT_ID, refresh.ServiceClientID)
			assert.NotEmpty(t, refresh.Token)
			assert.False(t, refresh.Expires.IsZero())
		}
		token, refresh, err := tservice.NewAccessToken(ctx, code.Code)
		testTokens(token, refresh)
		token, refresh, err = tservice.UpdateAccessToken(ctx, refresh.Token)
		testTokens(token, refresh)
	})

	t.Run("expired", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tservice := NewService(Config{})

		err := tservice.Database.CreateAuthorizationCode(ctx, database.AuthorizationCode{
			Code:    "example",
			Expires: time.Now().Add(time.Duration(-1) * time.Minute),
		})
		assert.NoError(t, err)

		_, _, err = tservice.NewAccessToken(ctx, "example")
		assert.ErrorIs(t, ErrAuthorizationCodeExpired, err)

		err = tservice.Database.CreateRefreshToken(ctx, database.RefreshToken{
			Token:   "example",
			Expires: time.Now().Add(time.Duration(-1) * time.Minute),
		})
		assert.NoError(t, err)
		_, _, err = tservice.UpdateAccessToken(ctx, "example")
		assert.ErrorIs(t, ErrRefreshTokenExpired, err)

	})
}
