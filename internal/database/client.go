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
		guser          *connect.BidiStreamForClient[apiv1.GetUserRequest, apiv1.GetUserResponse]
		gserviceClient *connect.BidiStreamForClient[apiv1.GetServiceClientRequest, apiv1.GetServiceClientResponse]
		gcode          *connect.BidiStreamForClient[apiv1.GetAuthorizationCodeRequest, apiv1.GetAuthorizationCodeResponse]
		gtoken         *connect.BidiStreamForClient[apiv1.GetAccessTokenRequest, apiv1.GetAccessTokenResponse]
		grefresh       *connect.BidiStreamForClient[apiv1.GetRefreshTokenRequest, apiv1.GetRefreshTokenResponse]
		ccode          *connect.BidiStreamForClient[apiv1.CreateAuthorizationCodeRequest, apiv1.CreateAuthorizationCodeResponse]
		ctoken         *connect.BidiStreamForClient[apiv1.CreateAccessTokenRequest, apiv1.CreateAccessTokenResponse]
		crefresh       *connect.BidiStreamForClient[apiv1.CreateRefreshTokenRequest, apiv1.CreateRefreshTokenResponse]
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
	c := apiv1connect.NewDatabaseServiceClient(
		httpClient,
		config.URL,
	)
	var client Client
	client.guser = c.GetUser(ctx)
	client.gserviceClient = c.GetServiceClient(ctx)
	client.gcode = c.GetAuthorizationCode(ctx)
	client.ccode = c.CreateAuthorizationCode(ctx)
	client.gtoken = c.GetAccessToken(ctx)
	client.ctoken = c.CreateAccessToken(ctx)
	client.grefresh = c.GetRefreshToken(ctx)
	client.crefresh = c.CreateRefreshToken(ctx)
	return &client, nil
}

func (c *Client) GetUserByID(ctx context.Context, id string) (*apiv1.UserProfile, error) {
	if err := c.guser.Send(&apiv1.GetUserRequest{
		Id: id,
	}); err != nil {
		return nil, err
	}
	resp, err := c.guser.Receive()
	if err != nil {
		return nil, err
	}
	return resp.GetUser(), nil
}

func (c *Client) GetServieClientByID(ctx context.Context, id string) (*apiv1.ServiceClient, error) {
	if err := c.gserviceClient.Send(&apiv1.GetServiceClientRequest{
		Id: id,
	}); err != nil {
		return nil, err
	}
	resp, err := c.gserviceClient.Receive()
	if err != nil {
		return nil, err
	}
	return resp.GetClient(), nil
}

func (c *Client) GetAuthorizationCodeByCode(ctx context.Context, code string) (*apiv1.AuthorizationCode, error) {
	if err := c.gcode.Send(&apiv1.GetAuthorizationCodeRequest{
		Code: code,
	}); err != nil {
		return nil, err
	}
	resp, err := c.gcode.Receive()
	if err != nil {
		return nil, err
	}
	return resp.GetCode(), nil
}

func (c *Client) CreateAuthorizationCode(ctx context.Context, row *apiv1.AuthorizationCode) error {
	if err := c.ccode.Send(&apiv1.CreateAuthorizationCodeRequest{
		Code: row,
	}); err != nil {
		return err
	}
	_, err := c.ccode.Receive()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAccessTokenByToken(ctx context.Context, token string) (*apiv1.AccessToken, error) {
	if err := c.gtoken.Send(&apiv1.GetAccessTokenRequest{
		Token: token,
	}); err != nil {
		return nil, err
	}
	resp, err := c.gtoken.Receive()
	if err != nil {
		return nil, err
	}
	return resp.GetToken(), nil
}

func (c *Client) CreateAccessToken(ctx context.Context, row *apiv1.AccessToken) error {
	if err := c.ctoken.Send(&apiv1.CreateAccessTokenRequest{
		Token: row,
	}); err != nil {
		return err
	}
	_, err := c.ctoken.Receive()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetRefreshTokenByToken(ctx context.Context, token string) (*apiv1.RefreshToken, error) {
	if err := c.grefresh.Send(&apiv1.GetRefreshTokenRequest{
		Token: token,
	}); err != nil {
		return nil, err
	}
	resp, err := c.grefresh.Receive()
	if err != nil {
		return nil, err
	}
	return resp.GetToken(), nil
}

func (c *Client) CreateRefreshToken(ctx context.Context, row *apiv1.RefreshToken) error {
	if err := c.crefresh.Send(&apiv1.CreateRefreshTokenRequest{
		Token: row,
	}); err != nil {
		return err
	}
	_, err := c.crefresh.Receive()
	if err != nil {
		return err
	}
	return nil
}
