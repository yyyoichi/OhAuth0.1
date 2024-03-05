package database

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Define a simple to understand the structure of OAuth2.0
type (
	Database struct {
		userByID                map[string]UserProfile
		serviceClientByID       map[string]ServiceClient
		authorizationCodeByCode map[string]AuthorizationCode
		accessTokenByToken      map[string]AccessToken
		refreshTokenByToken     map[string]RefreshToken
		mu                      sync.Mutex
	}
	// OhAuth0.1を利用しているユーザ情報
	UserProfile struct {
		ID string
		// store no-hashed passwards
		password string
		Name     string
		Age      uint8
		Profile  string
	}
	// OhAuth0.1を利用するクライアント情報
	ServiceClient struct {
		ID          string
		Secret      string
		Name        string
		RedirectURI string
		Scope       string // profile:view
	}

	// 認可コード
	AuthorizationCode struct {
		Code                    string
		UserID, ServiceClientID string
		Expires                 time.Time
		Scope                   string // profile:view
	}

	// アクセストークン
	AccessToken struct {
		Token                   string
		UserID, ServiceClientID string
		Expires                 time.Time
		Scope                   string // profile:view
	}

	// リフレッシュトークン
	RefreshToken struct {
		Token                   string
		UserID, ServiceClientID string
		Expires                 time.Time
		Scope                   string // profile:view
	}
)

func NewDatabase() (*Database, error) {
	var db Database
	db.userByID = map[string]UserProfile{
		"0": {"0", "password", "Taro", 20, "Hello🎈"},
		"1": {"1", "password", "Hanako", 20, "Hello🌸"},
	}
	db.serviceClientByID = map[string]ServiceClient{
		// TODO set RedirectURI
		"500": {"500", "secret", "ABC-App", "", "profile:view"},
		"501": {"501", "secret", "ZZZ-App", "", "profile:view"},
	}
	db.authorizationCodeByCode = make(map[string]AuthorizationCode)
	db.accessTokenByToken = make(map[string]AccessToken)
	db.refreshTokenByToken = make(map[string]RefreshToken)
	return &db, nil
}

func (db *Database) GetUserByID(ctx context.Context, id string) (*UserProfile, error) {
	u, found := db.userByID[id]
	if !found {
		return nil, ErrNotFound
	}
	return &u, nil
}

func (db *Database) GetServieClientByID(ctx context.Context, id string) (*ServiceClient, error) {
	c, found := db.serviceClientByID[id]
	if !found {
		return nil, ErrNotFound
	}
	return &c, nil
}

func (db *Database) GetAuthorizationCodeByCode(ctx context.Context, code string) (*AuthorizationCode, error) {
	c, found := db.authorizationCodeByCode[code]
	if !found {
		return nil, ErrNotFound
	}
	return &c, nil
}

func (db *Database) CreateAuthorizationCode(ctx context.Context, row AuthorizationCode) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, found := db.authorizationCodeByCode[row.Code]; found {
		return ErrAlreadyExists
	}
	db.authorizationCodeByCode[row.Code] = row
	return nil
}

func (db *Database) GetAccessTokenByToken(ctx context.Context, token string) (*AccessToken, error) {
	t, found := db.accessTokenByToken[token]
	if !found {
		return nil, ErrNotFound
	}
	return &t, nil
}

func (db *Database) CreateAccessToken(ctx context.Context, row AccessToken) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, found := db.accessTokenByToken[row.Token]; found {
		return ErrAlreadyExists
	}
	db.accessTokenByToken[row.Token] = row
	return nil
}

func (db *Database) GetRefreshTokenByToken(ctx context.Context, token string) (*RefreshToken, error) {
	t, found := db.refreshTokenByToken[token]
	if !found {
		return nil, ErrNotFound
	}
	return &t, nil
}

func (db *Database) CreateRefreshToken(ctx context.Context, row RefreshToken) error {
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
