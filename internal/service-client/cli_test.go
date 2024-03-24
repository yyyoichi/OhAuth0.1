package serviceclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yyyoichi/OhAuth0.1/internal/auth"
	"github.com/yyyoichi/OhAuth0.1/internal/resource"
)

func TestBrawser(t *testing.T) {
	t.Run("login", func(t *testing.T) {
		brawser := newBrawserMock(t)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
		defer cancel()
		assert.ErrorIs(t, brawser.login(ctx), ErrNoSite)
		_ = brawser.moveToServiceClient("TEST_ID")
		go func() {
			for {
				resp, err := http.DefaultClient.Get("http://localhost:9010/status")
				if err != nil {
					continue
				}
				if resp.StatusCode != http.StatusOK {
					continue
				}
				resp, err = http.DefaultClient.Get("http://localhost:9010?code=12345")
				if err != nil {
					continue
				}
				if resp.StatusCode != http.StatusOK {
					continue
				}
				return
			}
		}()
		assert.Nil(t, brawser.login(ctx)) // exp send code from goroutin
		assert.Equal(t, "accesstoken", brawser.accessTokens["TEST_ID"])
		assert.Equal(t, "refreshtoken", brawser.refreshTokens["TEST_ID"])
		assert.ErrorIs(t, brawser.login(ctx), ErrAlreadyLogin)

		assert.Nil(t, context.Cause(ctx))
	})

	t.Run("view profile", func(t *testing.T) {
		brawser := newBrawserMock(t)
		ctx := context.Background()

		// no switced
		_, err := brawser.viewProfile(ctx)
		assert.Error(t, err)

		// regular
		_ = brawser.moveToServiceClient("TEST_ID")
		brawser.accessTokens["TEST_ID"] = "accesstoken"
		brawser.refreshTokens["TEST_ID"] = "refreshtoken"
		p, err := brawser.viewProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "1", p["id"])
		assert.Equal(t, "name", p["name"])
		assert.EqualValues(t, 20, p["age"])
		assert.Equal(t, "my profile", p["profile"])

		// no switched
		brawser.logout()
		_, err = brawser.viewProfile(ctx)
		assert.Error(t, err)
	})

	t.Run("refreshtoken", func(t *testing.T) {
		brawser := newBrawserMock(t)
		ctx := context.Background()

		assert.Error(t, brawser.refreshToken(ctx))

		_ = brawser.moveToServiceClient("TEST_ID")
		brawser.accessTokens["TEST_ID"] = ""
		brawser.refreshTokens["TEST_ID"] = "token"
		assert.Nil(t, brawser.refreshToken(ctx))
		assert.NotEmpty(t, brawser.accessTokens["TEST_ID"])
	})

}

type resourceClientMock struct {
	count int
}

// ViewProfile implements resourceClientInterface.
func (r *resourceClientMock) ViewProfile(ctx context.Context, token string) (*resource.ProfileGetResponse, error) {
	if r.count == 0 {
		r.count++
		return nil, resource.ErrAccessTokenExpired
	}
	return &resource.ProfileGetResponse{}, nil
}

func TestRetryVeiwProfile(t *testing.T) {
	brawser := newBrawserMock(t)
	mock := resourceClientMock{}
	brawser.resourceClient = &mock
	_ = brawser.moveToServiceClient("TEST_ID")
	brawser.accessTokens["TEST_ID"] = "token"
	brawser.refreshTokens["TEST_ID"] = "token"
	_, err := brawser.viewProfile(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, mock.count) // !
}

func newBrawserMock(t *testing.T) Brawser {
	var brawser Brawser
	brawser.codeReceiver = NewCodeReceiver(9010)
	brawser.accessTokenClient = func() AccessTokenClient {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		a := auth.AccessTokenResponse{
			AccessToken:  "accesstoken",
			RefreshToken: "refreshtoken",
			ExpiresIn:    500,
		}
		b, err := json.Marshal(a)
		assert.NoError(t, err)
		resp.Write([]byte(b))
		return AccessTokenClient{
			post: func(_ context.Context, _ string, _ io.Reader) (*http.Response, error) {
				return resp.Result(), nil
			},
		}
	}()
	brawser.resourceClient = func() *ResourceClient {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		p := resource.ProfileGetResponse{
			UserId:  "1",
			Name:    "name",
			Age:     20,
			Profile: "my profile",
		}
		b, err := json.Marshal(p)
		assert.NoError(t, err)
		resp.Write([]byte(b))
		return &ResourceClient{
			get: func(_ context.Context, _0, _1 string) (*http.Response, error) {
				return resp.Result(), nil
			},
		}
	}()
	brawser.mu = &sync.Mutex{}
	brawser.accessTokens = map[string]string{}
	brawser.refreshTokens = map[string]string{}
	return brawser
}
