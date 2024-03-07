package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
)

var JWT_SECRET = []byte("JWT_SECRET")

func SetupRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	v1 := api.Group("/v1")

	service := NewService(Config{})

	v1.GET("/clients/:client_id", func(ctx *gin.Context) {
		var req ServiceClientGetRequest
		if err := ctx.Bind(&req); err != nil {
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
		if err := ctx.Bind(&req); err != nil {
			ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
			return
		}
		slog.InfoContext(ctx, "recieve", "body", req)
		claims, err := service.Authentication(ctx, req.ClientID, req.Password)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot authenticate: %v", err))
			if errors.Is(err, database.ErrNotFound) || errors.Is(err, ErrNoMatchPassword) {
				ctx.SecureJSON(http.StatusBadRequest, BadRequestMessage)
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
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, BadRequestMessage)
			return
		}
		claims, err := service.ParseMyClaims(ctx, req.JWT, JWT_SECRET)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot parse jwt: %v", err))
			ctx.JSON(http.StatusBadRequest, BadRequestMessage)
			return
		}
		if claims.ClientID != req.ClientID {
			slog.ErrorContext(ctx, fmt.Sprintf("cannot match clientID jwt:%s, req:%s", claims.ClientID, req.ClientID))
			ctx.JSON(http.StatusBadRequest, BadRequestMessage)
			return
		}
		var resp AuthorizationResponse
		ctx.SecureJSON(http.StatusOK, resp)
	})

	v1.POST("/accesstoken", func(ctx *gin.Context) {
		var req AccessTokenRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		var resp AccessTokenResponse
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
