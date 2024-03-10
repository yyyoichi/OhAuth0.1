package database

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
	"github.com/yyyoichi/OhAuth0.1/api/v1/apiv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type (
	ServerConfig struct {
		Port string
	}
	handler struct {
		*Database
		// apiv1connect.UnimplementedDatabaseServiceHandler
	}
)

func NewDatabaseServer(config ServerConfig) error {
	if config.Port == "" {
		config.Port = "3306"
	}
	addr := fmt.Sprintf(":%s", config.Port)
	db, err := NewDatabase()
	if err != nil {
		return err
	}

	rpc := http.NewServeMux()
	rpc.Handle(apiv1connect.NewDatabaseServiceHandler(&handler{
		Database: db,
	}))
	server := &http.Server{
		Addr:         addr,
		Handler:      h2c.NewHandler(rpc, &http2.Server{}),
		ReadTimeout:  time.Duration(5) * time.Second, // クライアントからのリクエスト読み取りタイムアウト
		WriteTimeout: time.Duration(5) * time.Second, // レスポンス書き込みタイムアウト
		IdleTimeout:  0,
	}
	go func() {
		log.Println("Starting HTTP/2 server on", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()
	return nil
}

// CreateAccessToken implements apiv1connect.DatabaseServiceHandler.
func (h *handler) CreateAccessToken(ctx context.Context, stream *connect.BidiStream[apiv1.CreateAccessTokenRequest, apiv1.CreateAccessTokenResponse]) error {
	for {
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		if err := h.Database.CreateAccessToken(ctx, msg.GetToken()); err != nil {
			return h.newConnectError(err)
		}
		if err := stream.Send(&apiv1.CreateAccessTokenResponse{}); err != nil {
			return err
		}
		continue
	}
}

// CreateAuthorizationCode implements apiv1connect.DatabaseServiceHandler.
func (h *handler) CreateAuthorizationCode(ctx context.Context, stream *connect.BidiStream[apiv1.CreateAuthorizationCodeRequest, apiv1.CreateAuthorizationCodeResponse]) error {
	for {
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		if err := h.Database.CreateAuthorizationCode(ctx, msg.GetCode()); err != nil {
			return h.newConnectError(err)
		}
		if err := stream.Send(&apiv1.CreateAuthorizationCodeResponse{}); err != nil {
			return err
		}
		continue
	}
}

// CreateRefreshToken implements apiv1connect.DatabaseServiceHandler.
func (h *handler) CreateRefreshToken(ctx context.Context, stream *connect.BidiStream[apiv1.CreateRefreshTokenRequest, apiv1.CreateRefreshTokenResponse]) error {
	for {
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		if err := h.Database.CreateRefreshToken(ctx, msg.GetToken()); err != nil {
			return h.newConnectError(err)
		}
		if err := stream.Send(&apiv1.CreateRefreshTokenResponse{}); err != nil {
			return err
		}
		continue
	}
}

// GetAccessToken implements apiv1connect.DatabaseServiceHandler.
// Subtle: this method shadows the method (DatabaseServiceHandler).GetAccessToken of handler.DatabaseServiceHandler.
func (h *handler) GetAccessToken(ctx context.Context, stream *connect.BidiStream[apiv1.GetAccessTokenRequest, apiv1.GetAccessTokenResponse]) error {
	for {
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		token, err := h.Database.GetAccessTokenByToken(ctx, msg.GetToken())
		if err != nil {
			return h.newConnectError(err)
		}
		if err := stream.Send(&apiv1.GetAccessTokenResponse{
			Token: token,
		}); err != nil {
			return err
		}
		continue
	}
}

// GetAuthorizationCode implements apiv1connect.DatabaseServiceHandler.
// Subtle: this method shadows the method (DatabaseServiceHandler).GetAuthorizationCode of handler.DatabaseServiceHandler.
func (h *handler) GetAuthorizationCode(ctx context.Context, stream *connect.BidiStream[apiv1.GetAuthorizationCodeRequest, apiv1.GetAuthorizationCodeResponse]) error {
	for {
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		code, err := h.Database.GetAuthorizationCodeByCode(ctx, msg.GetCode())
		if err != nil {
			return h.newConnectError(err)
		}
		if err := stream.Send(&apiv1.GetAuthorizationCodeResponse{
			Code: code,
		}); err != nil {
			return err
		}
		continue
	}
}

// GetRefreshToken implements apiv1connect.DatabaseServiceHandler.
// Subtle: this method shadows the method (DatabaseServiceHandler).GetRefreshToken of handler.DatabaseServiceHandler.
func (h *handler) GetRefreshToken(ctx context.Context, stream *connect.BidiStream[apiv1.GetRefreshTokenRequest, apiv1.GetRefreshTokenResponse]) error {
	for {
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		token, err := h.Database.GetRefreshTokenByToken(ctx, msg.GetToken())
		if err != nil {
			return h.newConnectError(err)
		}
		if err := stream.Send(&apiv1.GetRefreshTokenResponse{
			Token: token,
		}); err != nil {
			return err
		}
		continue
	}
}

// GetServiceClient implements apiv1connect.DatabaseServiceHandler.
// Subtle: this method shadows the method (DatabaseServiceHandler).GetServiceClient of handler.DatabaseServiceHandler.
func (h *handler) GetServiceClient(ctx context.Context, stream *connect.BidiStream[apiv1.GetServiceClientRequest, apiv1.GetServiceClientResponse]) error {
	for {
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		client, err := h.Database.GetServieClientByID(ctx, msg.GetId())
		if err != nil {
			return h.newConnectError(err)
		}
		if err := stream.Send(&apiv1.GetServiceClientResponse{
			Client: client,
		}); err != nil {
			return err
		}
		continue
	}
}

// GetUser implements apiv1connect.DatabaseServiceHandler.
// Subtle: this method shadows the method (DatabaseServiceHandler).GetUser of handler.DatabaseServiceHandler.
func (h *handler) GetUser(ctx context.Context, stream *connect.BidiStream[apiv1.GetUserRequest, apiv1.GetUserResponse]) error {
	for {
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		user, err := h.Database.GetUserByID(ctx, msg.GetId())
		if err != nil {
			return h.newConnectError(err)
		}
		if err := stream.Send(&apiv1.GetUserResponse{
			User: user,
		}); err != nil {
			return err
		}
		continue
	}
}

// Ping implements apiv1connect.DatabaseServiceHandler.
func (h *handler) Ping(context.Context, *connect.Request[apiv1.PingRequest]) (*connect.Response[apiv1.PingResponse], error) {
	return &connect.Response[apiv1.PingResponse]{}, nil
}

func (c *handler) newConnectError(err error) error {
	if errors.Is(ErrNotFound, err) {
		return connect.NewError(connect.CodeNotFound, err)
	}
	if errors.Is(ErrAlreadyExists, err) {
		return connect.NewError(connect.CodeAlreadyExists, err)
	}
	return connect.NewError(connect.CodeInternal, err)
}
