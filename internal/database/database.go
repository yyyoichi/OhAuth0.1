package database

import (
	"context"
	"errors"
	"sync"

	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
)

// Define a simple to understand the structure of OAuth2.0
type Database struct {
	userById                map[string]*apiv1.UserProfile
	serviceClientById       map[string]*apiv1.ServiceClient
	authorizationCodeByCode map[string]*apiv1.AuthorizationCode
	accessTokenByToken      map[string]*apiv1.AccessToken
	refreshTokenByToken     map[string]*apiv1.RefreshToken
	mu                      sync.Mutex
}

const (
	CLIENT_SECRET = "secret"
	REDIRECT_URI  = "http://localhost:7777"
)

var (
	MockServiceClient500 = apiv1.ServiceClient{
		Id:          "500",
		Name:        "Professional Q&A",
		Secret:      CLIENT_SECRET,
		RedirectUri: REDIRECT_URI,
		Scope:       "profile:view",
	}
	MockServiceClient501 = apiv1.ServiceClient{
		Id:          "501",
		Name:        "Complete Offece",
		Secret:      CLIENT_SECRET,
		RedirectUri: REDIRECT_URI,
		Scope:       "profile:view",
	}
)

func NewDatabase() (*Database, error) {
	var db Database
	db.userById = map[string]*apiv1.UserProfile{
		"1": {
			Id:       "1",
			Password: "password",
			Name:     "Taro",
			Age:      20,
			Profile:  "Hello🎈",
		},
		"2": {
			Id:       "2",
			Password: "password",
			Name:     "Hanako",
			Age:      20,
			Profile:  "Hello🌸",
		},
	}
	db.serviceClientById = map[string]*apiv1.ServiceClient{
		"500": {
			Id:          MockServiceClient500.Id,
			Name:        MockServiceClient500.Name,
			Secret:      MockServiceClient500.Secret,
			RedirectUri: MockServiceClient500.RedirectUri,
			Scope:       MockServiceClient500.Scope,
		},
		"501": {
			Id:          MockServiceClient501.Id,
			Name:        MockServiceClient501.Name,
			Secret:      MockServiceClient501.Secret,
			RedirectUri: MockServiceClient501.RedirectUri,
			Scope:       MockServiceClient501.Scope,
		},
	}
	db.authorizationCodeByCode = make(map[string]*apiv1.AuthorizationCode)
	db.accessTokenByToken = make(map[string]*apiv1.AccessToken)
	db.refreshTokenByToken = make(map[string]*apiv1.RefreshToken)
	return &db, nil
}

func (db *Database) GetUserById(ctx context.Context, id string) (*apiv1.UserProfile, error) {
	u, found := db.userById[id]
	if !found {
		return nil, ErrNotFound
	}
	return u, nil
}

func (db *Database) GetServieClientById(ctx context.Context, id string) (*apiv1.ServiceClient, error) {
	c, found := db.serviceClientById[id]
	if !found {
		return nil, ErrNotFound
	}
	return c, nil
}

func (db *Database) GetAuthorizationCodeByCode(ctx context.Context, code string) (*apiv1.AuthorizationCode, error) {
	c, found := db.authorizationCodeByCode[code]
	if !found {
		return nil, ErrNotFound
	}
	return c, nil
}

func (db *Database) CreateAuthorizationCode(ctx context.Context, row *apiv1.AuthorizationCode) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, found := db.authorizationCodeByCode[row.Code]; found {
		return ErrAlreadyExists
	}
	db.authorizationCodeByCode[row.Code] = row
	return nil
}

func (db *Database) GetAccessTokenByToken(ctx context.Context, token string) (*apiv1.AccessToken, error) {
	t, found := db.accessTokenByToken[token]
	if !found {
		return nil, ErrNotFound
	}
	return t, nil
}

func (db *Database) CreateAccessToken(ctx context.Context, row *apiv1.AccessToken) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, found := db.accessTokenByToken[row.Token]; found {
		return ErrAlreadyExists
	}
	db.accessTokenByToken[row.Token] = row
	return nil
}

func (db *Database) GetRefreshTokenByToken(ctx context.Context, token string) (*apiv1.RefreshToken, error) {
	t, found := db.refreshTokenByToken[token]
	if !found {
		return nil, ErrNotFound
	}
	return t, nil
}

func (db *Database) CreateRefreshToken(ctx context.Context, row *apiv1.RefreshToken) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, found := db.refreshTokenByToken[row.Token]; found {
		return ErrAlreadyExists
	}
	db.refreshTokenByToken[row.Token] = row
	return nil
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)
