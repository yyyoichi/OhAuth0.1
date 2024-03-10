package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDatabase(t *testing.T) {
	test := map[string]databaseInterface{
		"local": func() *Database {
			db, _ := NewDatabase()
			return db
		}(),
		"remote": func() *Client {
			port := "3366"
			err := NewDatabaseServer(ServerConfig{
				Port: port,
			})
			assert.NoError(t, err)
			var client *Client
			connectedCh := make(chan struct{})
			go func() {
				defer close(connectedCh)
				for {
					ctx := context.Background()
					var err error
					client, err = NewDatabaseClient(ctx, ClientConfig{
						URL: "http://localhost:" + port,
					})
					if err != nil {
						continue
					}
					if err := client.Ping(ctx); err != nil {
						continue
					}
					return
				}
			}()
			<-connectedCh
			return client
		}(),
	}
	for scenario, tt := range test {
		t.Run(scenario, func(t *testing.T) {
			testdb(t, tt)
		})
	}
}

func testdb(t *testing.T, db databaseInterface) {
	ctx := context.Background()
	user, err := db.GetUserByID(ctx, "1")
	assert.NoError(t, err)
	assert.NotZero(t, user)
	client, err := db.GetServieClientByID(ctx, "500")
	assert.NoError(t, err)
	assert.NotZero(t, client)
	NOW := timestamppb.Now()
	expcode := &apiv1.AuthorizationCode{
		Code:            "code",
		UserId:          "11",
		ServiceClientId: "222",
		Expires:         NOW,
		Scope:           "hoge",
	}
	err = db.CreateAuthorizationCode(ctx, expcode)
	assert.NoError(t, err)
	code, err := db.GetAuthorizationCodeByCode(ctx, expcode.Code)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, expcode, code)
	exptoken := &apiv1.AccessToken{
		Token:           "token",
		UserId:          "11",
		ServiceClientId: "222",
		Expires:         NOW,
		Scope:           "hoge",
	}
	err = db.CreateAccessToken(ctx, exptoken)
	assert.NoError(t, err)
	token, err := db.GetAccessTokenByToken(ctx, exptoken.Token)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, exptoken, token)
	exprefresh := &apiv1.RefreshToken{
		Token:           "token",
		UserId:          "11",
		ServiceClientId: "222",
		Expires:         NOW,
		Scope:           "hoge",
	}
	err = db.CreateRefreshToken(ctx, exprefresh)
	assert.NoError(t, err)
	refresh, err := db.GetRefreshTokenByToken(ctx, exprefresh.Token)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, exprefresh, refresh)
}

type databaseInterface interface {
	GetUserByID(ctx context.Context, id string) (*apiv1.UserProfile, error)
	GetServieClientByID(ctx context.Context, id string) (*apiv1.ServiceClient, error)
	GetAuthorizationCodeByCode(ctx context.Context, code string) (*apiv1.AuthorizationCode, error)
	CreateAuthorizationCode(ctx context.Context, row *apiv1.AuthorizationCode) error
	GetAccessTokenByToken(ctx context.Context, token string) (*apiv1.AccessToken, error)
	CreateAccessToken(ctx context.Context, row *apiv1.AccessToken) error
	GetRefreshTokenByToken(ctx context.Context, token string) (*apiv1.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, row *apiv1.RefreshToken) error
}
