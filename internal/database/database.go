package database

import (
	"context"
	"errors"
	"fmt"
)

// Define a simple to understand the structure of OAuth2.0
type (
	Database struct {
		userByID          map[string]UserProfile
		serviceClientByID map[string]ServiceClient
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
		Scope       string // いつもプロフィールの閲覧
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
	return &db, nil
}

func (db *Database) Login(ctx context.Context, id, pass string) (*UserProfile, error) {
	u, found := db.userByID[id]
	if !found {
		return nil, ErrNotFound
	}
	if u.password != pass {
		return nil, fmt.Errorf("password is invalid")
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

var (
	ErrNotFound = errors.New("not found")
)
