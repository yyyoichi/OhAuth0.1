package serviceclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
