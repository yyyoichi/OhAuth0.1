package resource

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService(t *testing.T) {
	t.Run("VerifyAccessToken", func(t *testing.T) {
		test := []struct {
			token    string
			expires  time.Time
			argToken string
			expErr   error
		}{
			{"token", time.Now().AddDate(1, 0, 0), "token", nil},
			{"token", time.Now().AddDate(1, 0, 0), "not found token", database.ErrNotFound},
			{"token", time.Now().AddDate(-1, 0, 0), "token", ErrAccessTokenExpired},
		}
		ctx := context.Background()
		for _, tt := range test {
			db, _ := database.NewDatabase()
			tservice := &Service{
				client: db,
			}
			err := db.CreateAccessToken(ctx, &apiv1.AccessToken{
				UserId:          "1",
				ServiceClientId: "501",
				Scope:           "profile:view",
				Token:           tt.token,
				Expires:         timestamppb.New(tt.expires),
			})
			assert.NoError(t, err)
			token, err := tservice.VerifyAccessToken(ctx, tt.argToken)
			if tt.expErr == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			} else {
				assert.ErrorIs(t, err, tt.expErr)
			}
		}
	})
	t.Run("ViewUserProfile", func(t *testing.T) {
		test := []struct {
			userId string
			expErr error
		}{
			{"1", nil},
			{"99", database.ErrNotFound},
		}
		ctx := context.Background()
		for _, tt := range test {
			db, _ := database.NewDatabase()
			tservice := &Service{
				client: db,
			}
			user, err := tservice.ViewUserProfile(ctx, tt.userId)
			if tt.expErr == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, user)
			} else {
				assert.ErrorIs(t, err, tt.expErr)
			}
		}
	})
}
