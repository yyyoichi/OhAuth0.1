package serviceclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/yyyoichi/OhAuth0.1/internal/auth"
)

type CodeReceiver struct {
	HostURI string
	codeCh  chan string
	closed  bool
}

func NewCodeReceiver(uri string) CodeReceiver {
	b := CodeReceiver{}
	b.Init()
	return b
}

func (b *CodeReceiver) Start(ctx context.Context) {
	if b.closed {
		panic("please Init()")
	}
	http.ListenAndServe(b.HostURI, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		code := r.URL.Query().Get("code")
		if code != "" {
			slog.ErrorContext(ctx, "code is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			b.closed = true
			close(b.codeCh)
		}()
		select {
		case <-ctx.Done():
			return
		default:
			b.codeCh <- code
		}
	}))
}

func (b *CodeReceiver) Receive(cxt context.Context) string {
	select {
	case <-cxt.Done():
		return ""
	case code, ok := <-b.codeCh:
		if !ok {
			return ""
		}
		return code
	}
}

func (b *CodeReceiver) Init() {
	b.codeCh = make(chan string)
	b.closed = false
}

type (
	AccessTokenClient struct {
		AuthServerURI string
	}
	AccessTokenRequestParam struct {
		ClientId     string
		ClientSecret string
	}
)

func (c *AccessTokenClient) GetByCode(code string, param AccessTokenRequestParam) (
	*auth.AccessTokenResponse, error,
) {
	var req auth.AccessTokenRequest
	req.GrantType = "authorization_code"
	req.ClientId = param.ClientId
	req.ClientSecret = param.ClientSecret
	req.Code = code
	return c.get(req)
}

func (c *AccessTokenClient) GetByRefreshToken(token string, param AccessTokenRequestParam) (
	*auth.AccessTokenResponse, error,
) {
	var req auth.AccessTokenRequest
	req.GrantType = "authorization_code"
	req.ClientId = param.ClientId
	req.ClientSecret = param.ClientSecret
	req.RefreshToken = token
	return c.get(req)
}

func (c *AccessTokenClient) get(req auth.AccessTokenRequest) (*auth.AccessTokenResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Post(
		c.AuthServerURI+"/api/v1/accesstoken",
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		var body struct {
			Status string
		}
		if err := json.Unmarshal(data, &body); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code is %d: %s", resp.StatusCode, body.Status)
	}
	var body auth.AccessTokenResponse
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, err
	}
	return &body, nil
}
