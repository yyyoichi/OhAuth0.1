package serviceclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/yyyoichi/OhAuth0.1/internal/auth"
)

type (
	CodeReceiver struct {
		HostURI string
		codeCh  chan codeResult
		closed  bool
	}
	codeResult struct {
		code string
		err  error
	}
)

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
		defer func() {
			b.closed = true
			close(b.codeCh)
		}()
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		code := r.URL.Query().Get("code")
		if code != "" {
			b.codeCh <- codeResult{
				err: errors.New("code is empty"),
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		select {
		case <-ctx.Done():
			b.codeCh <- codeResult{
				err: context.Cause(ctx),
			}
			w.WriteHeader(http.StatusRequestTimeout)
			return
		default:
			b.codeCh <- codeResult{
				code: code,
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}))
}

func (b *CodeReceiver) Receive() codeResult {
	return <-b.codeCh
}

func (b *CodeReceiver) Init() {
	b.codeCh = make(chan codeResult)
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
