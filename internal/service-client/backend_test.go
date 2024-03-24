package serviceclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yyyoichi/OhAuth0.1/internal/auth"
	"github.com/yyyoichi/OhAuth0.1/internal/resource"
)

func TestCodeReciever(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		t.Parallel()
		turi := "http://localhost:9001"
		tserver := NewCodeReceiver(9001)
		ctx := context.Background()
		tserver.Start(ctx)
		go func() {
			resp, err := http.DefaultClient.Get(turi + "?code=12345")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}()
		result := <-tserver.Receive()
		assert.NoError(t, result.err)
		assert.Equal(t, "12345", result.code)
	})
	t.Run("post method and context cancel", func(t *testing.T) {
		t.Parallel()
		turi := "http://localhost:9002"
		tserver := NewCodeReceiver(9002)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		tserver.Start(ctx)
		go func() {
			resp, err := http.DefaultClient.Post(turi+"?code=12345", "application/json", nil)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

			cancel() // !

			resp, err = http.DefaultClient.Get(turi + "?code=12345")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusRequestTimeout, resp.StatusCode)
		}()
		result, ok := <-tserver.Receive()
		assert.True(t, ok)
		assert.Equal(t, context.Canceled, result.err)

		_, ok = <-tserver.Receive()
		assert.False(t, ok)
	})
}

func TestAccessTokenClient(t *testing.T) {
	test := []struct {
		statusCode int
	}{
		{http.StatusOK},
		{http.StatusBadRequest},
	}
	for _, tt := range test {
		var client = func() AccessTokenClient {
			resp := httptest.NewRecorder()
			resp.WriteHeader(tt.statusCode)
			resp.Write([]byte("{}"))
			return AccessTokenClient{
				post: func(_ context.Context, _ string, _ io.Reader) (*http.Response, error) {
					return resp.Result(), nil
				},
			}
		}()
		_, err := client.get(context.Background(), auth.AccessTokenRequest{})
		if tt.statusCode != http.StatusOK {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestResourceClient(t *testing.T) {
	test := []struct {
		statusCode int
	}{
		{http.StatusOK},
		{http.StatusBadRequest},
		{http.StatusUnauthorized},
	}
	for _, tt := range test {
		var client = func() ResourceClient {
			resp := httptest.NewRecorder()
			resp.WriteHeader(tt.statusCode)
			resp.Write([]byte("{}"))
			return ResourceClient{
				get: func(_ context.Context, _0, _1 string) (*http.Response, error) {
					return resp.Result(), nil
				},
			}
		}()
		_, err := client.ViewProfile(context.Background(), "")
		if tt.statusCode == http.StatusUnauthorized {
			assert.ErrorIs(t, err, resource.ErrAccessTokenExpired)
		}
		if tt.statusCode != http.StatusOK {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
