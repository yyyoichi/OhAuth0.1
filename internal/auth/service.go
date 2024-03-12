package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	Service struct {
		client clientInterface
	}
	clientInterface interface {
		GetUserById(ctx context.Context, id string) (*apiv1.UserProfile, error)
		GetServieClientById(ctx context.Context, id string) (*apiv1.ServiceClient, error)
		CreateAuthorizationCode(ctx context.Context, row *apiv1.AuthorizationCode) error
		GetAuthorizationCodeByCode(ctx context.Context, code string) (*apiv1.AuthorizationCode, error)
		CreateAccessToken(ctx context.Context, row *apiv1.AccessToken) error
		GetRefreshTokenByToken(ctx context.Context, token string) (*apiv1.RefreshToken, error)
		CreateRefreshToken(ctx context.Context, row *apiv1.RefreshToken) error
	}
	Config struct {
		DatabaseServerURL string
	}
	MyClaims struct {
		ClientId string `json:"client_id"`
		jwt.RegisteredClaims
	}
)

var (
	ErrNoMatchPassword          = errors.New("no match password")
	ErrAuthorizationCodeExpired = errors.New("authorization code is expired")
	ErrRefreshTokenExpired      = errors.New("refresh token is expired")
)

func NewService(ctx context.Context, config Config) (*Service, error) {
	client, err := database.NewDatabaseClient(ctx, database.ClientConfig{
		URL: config.DatabaseServerURL,
	})
	if err != nil {
		return nil, err
	}
	return &Service{
		client: client,
	}, nil
}

// UserIdと有効期限を詰めたClaimsを返す
func (s *Service) Authentication(ctx context.Context, id, password string) (*MyClaims, error) {
	u, err := s.client.GetUserById(ctx, id)
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
			Subject:   u.GetId(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(10) * time.Minute).In(tz)),
		},
	}, nil
}

func (s *Service) ParseMyClaims(ctx context.Context, ss string, secret []byte) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(ss, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
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
	UserId, ServiceClientId string
	// Scope string
}

// 認可コードを発行する
func (s *Service) NewAuthorizationCode(ctx context.Context, config NewAuthorizationCodeConfig) (*apiv1.AuthorizationCode, error) {
	row := apiv1.AuthorizationCode{
		UserId:          config.UserId,
		ServiceClientId: config.ServiceClientId,
		Expires:         timestamppb.New(time.Now().Add(time.Duration(10) * time.Minute)),
		Scope:           "profile:view",
		Code:            uuid.NewString(),
	}
	if err := s.client.CreateAuthorizationCode(ctx, &row); err != nil {
		return nil, err
	}
	return &row, nil
}

// 認可コード[code]を検証しアクセストークンを発行する
func (s *Service) NewAccessToken(ctx context.Context, code string) (
	*apiv1.AccessToken,
	*apiv1.RefreshToken,
	error,
) {
	authorization, err := s.client.GetAuthorizationCodeByCode(ctx, code)
	if err != nil {
		return nil, nil, err
	}
	if time.Now().After(authorization.Expires.AsTime()) {
		return nil, nil, ErrAuthorizationCodeExpired
	}
	token := apiv1.AccessToken{
		Token:           uuid.NewString(),
		UserId:          authorization.UserId,
		ServiceClientId: authorization.ServiceClientId,
		Scope:           authorization.Scope,
		Expires:         timestamppb.New(time.Now().AddDate(0, 0, 3)),
	}
	refresh := apiv1.RefreshToken{
		Token:           uuid.NewString(),
		UserId:          authorization.UserId,
		ServiceClientId: authorization.ServiceClientId,
		Scope:           authorization.Scope,
		Expires:         timestamppb.New(time.Now().AddDate(0, 1, 0)),
	}
	if err := s.client.CreateAccessToken(ctx, &token); err != nil {
		return nil, nil, err
	}
	if err := s.client.CreateRefreshToken(ctx, &refresh); err != nil {
		return nil, nil, err
	}

	return &token, &refresh, nil
}

// [refreshToken]から新しくアクセストークンを発行する
func (s *Service) UpdateAccessToken(ctx context.Context, refreshToken string) (
	*apiv1.AccessToken,
	*apiv1.RefreshToken,
	error,
) {
	refresh, err := s.client.GetRefreshTokenByToken(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}
	if time.Now().After(refresh.Expires.AsTime()) {
		return nil, nil, ErrRefreshTokenExpired
	}
	updateToken := apiv1.AccessToken{
		Token:           uuid.NewString(),
		UserId:          refresh.UserId,
		ServiceClientId: refresh.ServiceClientId,
		Scope:           refresh.Scope,
		Expires:         timestamppb.New(time.Now().AddDate(0, 0, 3)),
	}
	updateRefresh := apiv1.RefreshToken{
		Token:           uuid.NewString(),
		UserId:          refresh.UserId,
		ServiceClientId: refresh.ServiceClientId,
		Scope:           refresh.Scope,
		Expires:         timestamppb.New(time.Now().AddDate(0, 1, 0)),
	}
	if err := s.client.CreateAccessToken(ctx, &updateToken); err != nil {
		return nil, nil, err
	}
	if err := s.client.CreateRefreshToken(ctx, &updateRefresh); err != nil {
		return nil, nil, err
	}

	return &updateToken, &updateRefresh, nil
}
