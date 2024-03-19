package database

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect"
	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
	"github.com/yyyoichi/OhAuth0.1/api/v1/apiv1connect"
	"golang.org/x/net/http2"
)

type (
	ClientConfig struct {
		URL string
	}
	Client struct {
		client apiv1connect.DatabaseServiceClient
	}
)

func NewDatabaseClient(ctx context.Context, config ClientConfig) (*Client, error) {
	httpClient := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
			ReadIdleTimeout: time.Duration(10) * time.Second, //10s接続がなかったらpingを開始
			PingTimeout:     time.Duration(15) * time.Second, //15sのpingに対する応答を待機
		},
	}
	var client Client
	client.client = apiv1connect.NewDatabaseServiceClient(
		httpClient,
		config.URL,
	)
	return &client, nil
}

func (c *Client) GetUserById(ctx context.Context, id string) (*apiv1.UserProfile, error) {
	cc := c.client.GetUser(ctx)
	if err := cc.Send(&apiv1.GetUserRequest{
		Id: id,
	}); err != nil {
		return nil, err
	}
	resp, err := cc.Receive()
	if err != nil {
		err := c.parseConnectError(err)
		return nil, err
	}
	return resp.GetUser(), nil
}

func (c *Client) GetServieClientById(ctx context.Context, id string) (*apiv1.ServiceClient, error) {
	cc := c.client.GetServiceClient(ctx)
	if err := cc.Send(&apiv1.GetServiceClientRequest{
		Id: id,
	}); err != nil {
		return nil, err
	}
	resp, err := cc.Receive()
	if err != nil {
		err := c.parseConnectError(err)
		return nil, err
	}
	return resp.GetClient(), nil
}

func (c *Client) GetAuthorizationCodeByCode(ctx context.Context, code string) (*apiv1.AuthorizationCode, error) {
	cc := c.client.GetAuthorizationCode(ctx)
	if err := cc.Send(&apiv1.GetAuthorizationCodeRequest{
		Code: code,
	}); err != nil {
		return nil, err
	}
	resp, err := cc.Receive()
	if err != nil {
		err := c.parseConnectError(err)
		return nil, err
	}
	return resp.GetCode(), nil
}

func (c *Client) CreateAuthorizationCode(ctx context.Context, row *apiv1.AuthorizationCode) error {
	cc := c.client.CreateAuthorizationCode(ctx)
	if err := cc.Send(&apiv1.CreateAuthorizationCodeRequest{
		Code: row,
	}); err != nil {
		return err
	}
	_, err := cc.Receive()
	if err != nil {
		err := c.parseConnectError(err)
		return err
	}
	return nil
}

func (c *Client) GetAccessTokenByToken(ctx context.Context, token string) (*apiv1.AccessToken, error) {
	cc := c.client.GetAccessToken(ctx)
	if err := cc.Send(&apiv1.GetAccessTokenRequest{
		Token: token,
	}); err != nil {
		return nil, err
	}
	resp, err := cc.Receive()
	if err != nil {
		err := c.parseConnectError(err)
		return nil, err
	}
	return resp.GetToken(), nil
}

func (c *Client) CreateAccessToken(ctx context.Context, row *apiv1.AccessToken) error {
	cc := c.client.CreateAccessToken(ctx)
	if err := cc.Send(&apiv1.CreateAccessTokenRequest{
		Token: row,
	}); err != nil {
		return err
	}
	_, err := cc.Receive()
	if err != nil {
		err := c.parseConnectError(err)
		return err
	}
	return nil
}

func (c *Client) GetRefreshTokenByToken(ctx context.Context, token string) (*apiv1.RefreshToken, error) {
	cc := c.client.GetRefreshToken(ctx)
	if err := cc.Send(&apiv1.GetRefreshTokenRequest{
		Token: token,
	}); err != nil {
		return nil, err
	}
	resp, err := cc.Receive()
	if err != nil {
		err := c.parseConnectError(err)
		return nil, err
	}
	return resp.GetToken(), nil
}

func (c *Client) CreateRefreshToken(ctx context.Context, row *apiv1.RefreshToken) error {
	cc := c.client.CreateRefreshToken(ctx)
	if err := cc.Send(&apiv1.CreateRefreshTokenRequest{
		Token: row,
	}); err != nil {
		return err
	}
	_, err := cc.Receive()
	if err != nil {
		err := c.parseConnectError(err)
		return err
	}
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &connect.Request[apiv1.PingRequest]{})
	return err
}

func (c *Client) parseConnectError(err error) error {
	connectErr, ok := err.(*connect.Error)
	if !ok {
		return err
	}
	switch connectErr.Code() {
	case connect.CodeAlreadyExists:
		return ErrAlreadyExists
	case connect.CodeNotFound:
		return ErrNotFound
	}
	return err
}
