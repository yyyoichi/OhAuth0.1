package serviceclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/yyyoichi/OhAuth0.1/internal/auth"
	"github.com/yyyoichi/OhAuth0.1/internal/resource"
)

type (
	CodeReceiver struct {
		Port   string
		codeCh chan codeResult
		closed bool
	}
	codeResult struct {
		code string
		err  error
	}
)

func NewCodeReceiver(port int) CodeReceiver {
	b := CodeReceiver{}
	b.Init()
	b.Port = ":" + strconv.Itoa(port)
	return b
}

func (b *CodeReceiver) Start(ctx context.Context) {
	if b.closed {
		panic("please Init()")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		defer func() {
			b.closed = true
			close(b.codeCh)
		}()
		code := r.URL.Query().Get("code")
		if code == "" {
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
	go func() {
		if err := http.ListenAndServe(b.Port, mux); err != nil {
			panic(err)
		}
	}()
	connected := make(chan struct{})
	go func() {
		defer close(connected)
		for {
			resp, err := http.DefaultClient.Get("http://localhost" + b.Port + "/status")
			if err != nil {
				continue
			}
			if resp.StatusCode != http.StatusOK {
				continue
			}
			return
		}
	}()
	<-connected
}

func (b *CodeReceiver) Receive() <-chan codeResult {
	return b.codeCh
}

func (b *CodeReceiver) Init() {
	b.codeCh = make(chan codeResult)
	b.closed = false
}

type (
	AccessTokenClient struct {
		post func(ctx context.Context, path string, body io.Reader) (resp *http.Response, err error)
	}
	AccessTokenRequestParam struct {
		ClientId     string
		ClientSecret string
	}
)

func NewAccessTokenClient(authServerURI string) AccessTokenClient {
	return AccessTokenClient{
		post: func(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, authServerURI+path, body)
			if err != nil {
				return nil, err
			}
			return http.DefaultClient.Do(req)
		},
	}
}

func (c *AccessTokenClient) GetByCode(ctx context.Context, code string, param AccessTokenRequestParam) (
	*auth.AccessTokenResponse, error,
) {
	var req auth.AccessTokenRequest
	req.GrantType = "authorization_code"
	req.ClientId = param.ClientId
	req.ClientSecret = param.ClientSecret
	req.Code = code
	return c.get(ctx, req)
}

func (c *AccessTokenClient) GetByRefreshToken(ctx context.Context, token string, param AccessTokenRequestParam) (
	*auth.AccessTokenResponse, error,
) {
	var req auth.AccessTokenRequest
	req.GrantType = "authorization_code"
	req.ClientId = param.ClientId
	req.ClientSecret = param.ClientSecret
	req.RefreshToken = token
	return c.get(ctx, req)
}

func (c *AccessTokenClient) get(ctx context.Context, req auth.AccessTokenRequest) (*auth.AccessTokenResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.post(ctx, "/api/v1/accesstoken", bytes.NewReader(b))
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

type ResourceClient struct {
	get func(ctx context.Context, path, token string) (*http.Response, error)
}

func NewResourceClient(resourceServerURI string) ResourceClient {
	return ResourceClient{
		get: func(ctx context.Context, path, token string) (*http.Response, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, resourceServerURI+path, nil)
			if err != nil {
				return nil, err
			}
			req.Header.Add("Authorization", "Bearer "+token)
			return http.DefaultClient.Do(req)
		},
	}
}

func (c *ResourceClient) ViewProfile(ctx context.Context, token string) (*resource.ProfileGetResponse, error) {
	resp, err := c.get(ctx, "/api/v1/profile", token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, resource.ErrAccessTokenExpired
		}
		var body struct {
			Status string
		}
		if err := json.Unmarshal(data, &body); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code is %d: %s", resp.StatusCode, body.Status)
	}
	var body resource.ProfileGetResponse
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, err
	}
	return &body, nil
}
