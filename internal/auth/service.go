package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
)

type (
	Service struct {
		*database.Database
	}
	Config   struct{}
	MyClaims struct {
		ClientID string `json:"client_id"`
		jwt.RegisteredClaims
	}
)

var (
	ErrNoMatchPassword = errors.New("no match password")
)

func NewService(config Config) *Service {
	db, _ := database.NewDatabase()
	return &Service{
		Database: db,
	}
}

// UserIDと有効期限を詰めたClaimsを返す
func (s *Service) Authentication(ctx context.Context, id, password string) (*MyClaims, error) {
	u, err := s.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get user: %w", err)
	}
	if u.Password != password {
		return nil, ErrNoMatchPassword
	}
	tz, _ := time.LoadLocation("Asia/Tokyo")
	return &MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "OhAuth0.1",
			Subject:   u.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(10) * time.Minute).In(tz)),
		},
	}, nil
}

func (s *Service) ParseMyClaims(ctx context.Context, ss string, secret []byte) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(ss, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if !token.Valid {
		// switch {
		// case token.Valid:
		// case errors.Is(err, jwt.ErrTokenMalformed):
		// case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		// case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		// default:
		// }
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*MyClaims)
	if !ok {
		return nil, errors.New("unexpected")
	}
	return claims, nil
}
