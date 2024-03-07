package auth

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
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
}
