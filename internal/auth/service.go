package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	ErrNoMatchPassword          = errors.New("no match password")
	ErrAuthorizationCodeExpired = errors.New("authorization code is expired")
	ErrRefreshTokenExpired      = errors.New("refresh token is expired")
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

type NewAuthorizationCodeConfig struct {
	UserID, ServiceClientID string
	// Scope string
}

// 認可コードを発行する
func (s *Service) NewAuthorizationCode(ctx context.Context, config NewAuthorizationCodeConfig) (*database.AuthorizationCode, error) {
	row := database.AuthorizationCode{
		UserID:          config.UserID,
		ServiceClientID: config.ServiceClientID,
		Expires:         time.Now().Add(time.Duration(10) * time.Minute),
		Scope:           "profile:view",
		Code:            uuid.NewString(),
	}
	if err := s.Database.CreateAuthorizationCode(ctx, row); err != nil {
		return nil, err
	}
	return &row, nil
}

// 認可コード[code]を検証しアクセストークンを発行する
func (s *Service) NewAccessToken(ctx context.Context, code string) (
	*database.AccessToken,
	*database.RefreshToken,
	error,
) {
	authorization, err := s.Database.GetAuthorizationCodeByCode(ctx, code)
	if err != nil {
		return nil, nil, err
	}
	if time.Now().After(authorization.Expires) {
		return nil, nil, ErrAuthorizationCodeExpired
	}
	token := database.AccessToken{
		Token:           uuid.NewString(),
		UserID:          authorization.UserID,
		ServiceClientID: authorization.ServiceClientID,
		Scope:           authorization.Scope,
		Expires:         time.Now().AddDate(0, 0, 3),
	}
	refresh := database.RefreshToken{
		Token:           uuid.NewString(),
		UserID:          authorization.UserID,
		ServiceClientID: authorization.ServiceClientID,
		Scope:           authorization.Scope,
		Expires:         time.Now().AddDate(0, 1, 0),
	}
	if err := s.Database.CreateAccessToken(ctx, token); err != nil {
		return nil, nil, err
	}
	if err := s.Database.CreateRefreshToken(ctx, refresh); err != nil {
		return nil, nil, err
	}

	return &token, &refresh, nil
}

// [refreshToken]から新しくアクセストークンを発行する
func (s *Service) UpdateAccessToken(ctx context.Context, refreshToken string) (
	*database.AccessToken,
	*database.RefreshToken,
	error,
) {
	refresh, err := s.Database.GetRefreshTokenByToken(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}
	if time.Now().After(refresh.Expires) {
		return nil, nil, ErrRefreshTokenExpired
	}
	updateToken := database.AccessToken{
		Token:           uuid.NewString(),
		UserID:          refresh.UserID,
		ServiceClientID: refresh.ServiceClientID,
		Scope:           refresh.Scope,
		Expires:         time.Now().AddDate(0, 0, 3),
	}
	updateRefresh := database.RefreshToken{
		Token:           uuid.NewString(),
		UserID:          refresh.UserID,
		ServiceClientID: refresh.ServiceClientID,
		Scope:           refresh.Scope,
		Expires:         time.Now().AddDate(0, 1, 0),
	}
	if err := s.Database.CreateAccessToken(ctx, updateToken); err != nil {
		return nil, nil, err
	}
	if err := s.Database.CreateRefreshToken(ctx, updateRefresh); err != nil {
		return nil, nil, err
	}

	return &updateToken, &updateRefresh, nil
}
