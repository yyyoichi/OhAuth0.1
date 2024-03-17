package auth

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
	enging "github.com/yyyoichi/OhAuth0.1/internal/engine"
)

var JWT_SECRET = []byte("JWT_SECRET")

func SetupRouter(service *Service, allowOrigins ...string) *gin.Engine {
	router := gin.Default()
	log.Printf("allow origins: %v", allowOrigins)
	// cross origin
	router.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin",
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		AllowCredentials: false,
	}))

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/clients/:client_id", func(ctx *gin.Context) {
		var req ServiceClientGetRequest
		if err := ctx.BindUri(&req); err != nil {
			ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
			return
		}
		slog.InfoContext(ctx, "recieve", "body", req)
		client, err := service.client.GetServieClientById(ctx, req.ClientId)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot get client: %v", err))
			if errors.Is(err, database.ErrNotFound) {
				ctx.SecureJSON(http.StatusNotFound, enging.NotFoundMessage)
				return
			}
			ctx.SecureJSON(http.StatusInternalServerError, enging.InternalServerErrorMessage)
			return
		}

		var resp ServiceClientGetResponse
		resp.ClientId = client.GetId()
		resp.Name = client.GetName()
		resp.Scope = client.GetScope()
		resp.RedirectUri = client.GetRedirectUri()
		ctx.SecureJSON(http.StatusOK, resp)
	})

	v1.POST("/authentication", func(ctx *gin.Context) {
		var req AuthenticationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
			return
		}
		slog.InfoContext(ctx, "recieve", "body", req)
		claims, err := service.Authentication(ctx, req.UserId, req.Password)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot authenticate: %v", err))
			if errors.Is(err, database.ErrNotFound) || errors.Is(err, ErrNoMatchPassword) {
				ctx.SecureJSON(http.StatusBadRequest, gin.H{"status": "Invalid Id or Password"})
				return
			}
			ctx.SecureJSON(http.StatusInternalServerError, enging.InternalServerErrorMessage)
			return
		}
		claims.ClientId = req.ClientId // !
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(JWT_SECRET)
		if err != nil {
			ctx.SecureJSON(http.StatusInternalServerError, enging.InternalServerErrorMessage)
			return
		}

		var resp AuthenticationResponse
		resp.JWT = ss
		ctx.SecureJSON(http.StatusOK, resp)
	})

	v1.POST("/authorization", func(ctx *gin.Context) {
		var req AuthorizationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
			return
		}
		claims, err := service.ParseMyClaims(ctx, req.JWT, JWT_SECRET)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot parse jwt: %v", err))
			ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
			return
		}
		if claims.ClientId != req.ClientId {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot match clientId jwt:%s, req:%s", claims.ClientId, req.ClientId))
			ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
			return
		}

		authorization, err := service.NewAuthorizationCode(ctx, NewAuthorizationCodeConfig{
			UserId:          claims.Subject,
			ServiceClientId: claims.ClientId,
		})
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot get authorization code: %v", err))
			ctx.SecureJSON(http.StatusInternalServerError, enging.InternalServerErrorMessage)
			return
		}

		var resp AuthorizationResponse
		resp.Code = authorization.GetCode()
		ctx.SecureJSON(http.StatusOK, resp)
	})

	v1.POST("/accesstoken", func(ctx *gin.Context) {
		var req AccessTokenRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
			return
		}
		// client secret
		if client, err := service.client.GetServieClientById(ctx, req.ClientId); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot get service client[%s]: %v", req.ClientId, err))
			ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
		} else if req.ClientSecret != client.Secret {
			slog.ErrorContext(ctx, fmt.Sprintf("client secret[%s] is invalid: %v", req.ClientSecret, err))
			ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
			return
		}

		var token *apiv1.AccessToken
		var refresh *apiv1.RefreshToken
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
			ctx.SecureJSON(http.StatusInternalServerError, enging.InternalServerErrorMessage)
			return
		}
		var resp AccessTokenResponse
		resp.AccessToken = token.GetToken()
		resp.RefreshToken = refresh.GetToken()
		resp.ExpiresIn = uint(time.Until(token.Expires.AsTime()).Seconds())
		ctx.SecureJSON(http.StatusOK, resp)
	})
	return router
}
