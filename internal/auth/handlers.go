package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
)

var JWT_SECRET = []byte("JWT_SECRET")

func SetupRouter(service *Service) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/clients/:client_id", func(ctx *gin.Context) {
		var req ServiceClientGetRequest
		if err := ctx.BindUri(&req); err != nil {
			ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
			return
		}
		slog.InfoContext(ctx, "recieve", "body", req)
		client, err := service.GetServieClientByID(ctx, req.ClientID)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot get client: %v", err))
			if errors.Is(err, database.ErrNotFound) {
				ctx.SecureJSON(http.StatusNotFound, NotFoundMessage)
				return
			}
			ctx.SecureJSON(http.StatusInternalServerError, InternalServerErrorMessage)
			return
		}

		var resp ServiceClientGetResponse
		resp.ClientID = client.ID
		resp.Name = client.Name
		resp.Scope = client.Scope
		ctx.SecureJSON(http.StatusOK, resp)
	})

	v1.POST("/authentication", func(ctx *gin.Context) {
		var req AuthenticationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
			return
		}
		slog.InfoContext(ctx, "recieve", "body", req)
		claims, err := service.Authentication(ctx, req.UserID, req.Password)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot authenticate: %v", err))
			if errors.Is(err, database.ErrNotFound) || errors.Is(err, ErrNoMatchPassword) {
				ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
				return
			}
			if errors.Is(err, database.ErrNotFound) {
				ctx.SecureJSON(http.StatusNotFound, NotFoundMessage)
				return
			}
			ctx.SecureJSON(http.StatusInternalServerError, InternalServerErrorMessage)
			return
		}
		claims.ClientID = req.ClientID // !
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(JWT_SECRET)
		if err != nil {
			ctx.SecureJSON(http.StatusInternalServerError, InternalServerErrorMessage)
			return
		}

		var resp AuthenticationResponse
		resp.JWT = ss
		ctx.SecureJSON(http.StatusOK, resp)
	})

	v1.POST("/authorization", func(ctx *gin.Context) {
		var req AuthorizationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
			return
		}
		claims, err := service.ParseMyClaims(ctx, req.JWT, JWT_SECRET)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot parse jwt: %v", err))
			ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
			return
		}
		if claims.ClientID != req.ClientID {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot match clientID jwt:%s, req:%s", claims.ClientID, req.ClientID))
			ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
			return
		}

		authorization, err := service.NewAuthorizationCode(ctx, NewAuthorizationCodeConfig{
			UserID:          claims.ID,
			ServiceClientID: claims.ClientID,
		})
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot get authorization code: %v", err))
			ctx.SecureJSON(http.StatusInternalServerError, InternalServerErrorMessage)
			return
		}

		var resp AuthorizationResponse
		resp.Code = authorization.Code
		ctx.SecureJSON(http.StatusOK, resp)
	})

	v1.POST("/accesstoken", func(ctx *gin.Context) {
		var req AccessTokenRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
			return
		}
		// client secret
		if client, err := service.GetServieClientByID(ctx, req.ClientID); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot get service client[%s]: %v", req.ClientID, err))
			ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
		} else if req.ClientSecret != client.Secret {
			slog.ErrorContext(ctx, fmt.Sprintf("client secret[%s] is invalid: %v", req.ClientSecret, err))
			ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
			return
		}

		var token *database.AccessToken
		var refresh *database.RefreshToken
		var err error
		switch {
		case req.Code != "":
			token, refresh, err = service.NewAccessToken(ctx, req.Code)
		case req.RefreshToken != "":
			token, refresh, err = service.UpdateAccessToken(ctx, req.RefreshToken)
		}
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot get tokens: %v", err))
			if errors.Is(err, ErrAuthorizationCodeExpired) || errors.Is(err, ErrRefreshTokenExpired) {
				ctx.SecureJSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
				return
			}
			ctx.SecureJSON(http.StatusInternalServerError, InternalServerErrorMessage)
			return
		}
		var resp AccessTokenResponse
		resp.AccessToken = token.Token
		resp.RefreshToken = refresh.Token
		resp.ExpiresIn = uint(time.Until(token.Expires).Seconds())
		ctx.SecureJSON(http.StatusOK, resp)
	})
	return router
}

var (
	BadRequestMessage = gin.H{
		"status": "Bad Request",
	}
	NotFoundMessage = gin.H{
		"status": "Not Found",
	}
	InternalServerErrorMessage = gin.H{
		"status": "Internal Server Error",
	}
)
