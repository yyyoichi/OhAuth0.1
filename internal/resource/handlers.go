package resource

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/yyyoichi/OhAuth0.1/api/v1"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
	enging "github.com/yyyoichi/OhAuth0.1/internal/engine"
)

const (
	USER_CONTEXT = "USER_CONTEXT"
)

func SetupRouter(service *Service) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	v1 := api.Group("/v1")
	v1.Use(func(ctx *gin.Context) {
		var h HeaderRequest
		if err := ctx.ShouldBindHeader(&h); err != nil {
			ctx.SecureJSON(http.StatusForbidden, enging.ForbiddenErrorMssage)
			return
		}
		accesstoken, err := h.FilterToken()
		if err != nil {
			slog.InfoContext(ctx, fmt.Sprintf("has not header: %v", err), slog.String("error", err.Error()), slog.String("Authorization", h.Authorization))
			ctx.SecureJSON(http.StatusForbidden, enging.ForbiddenErrorMssage)
			return
		}
		token, err := service.VerifyAccessToken(ctx, accesstoken)
		if err != nil {
			slog.InfoContext(ctx, "cannot varify accesstoken", slog.String("error", err.Error()), slog.String("access token", accesstoken))
			if errors.Is(err, database.ErrNotFound) {
				ctx.SecureJSON(http.StatusForbidden, enging.ForbiddenErrorMssage)
				return
			}
			if errors.Is(err, ErrAccessTokenExpired) {
				ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
				return
			}
			ctx.SecureJSON(http.StatusInternalServerError, enging.InternalServerErrorMessage)
			return
		}
		ctx.Set(USER_CONTEXT, token)
		ctx.Next()
	})
	getUser := func(ctx *gin.Context) (*apiv1.AccessToken, error) {
		output, ok := ctx.Get(USER_CONTEXT)
		if !ok {
			return nil, errors.New("unexpected")
		}
		user, ok := output.(*apiv1.AccessToken)
		if !ok {
			return nil, errors.New("unexpected")
		}
		return user, nil
	}
	v1.GET("/status", func(ctx *gin.Context) {
		ctx.SecureJSON(http.StatusNoContent, struct{}{})
	})
	v1.GET("/profile", func(ctx *gin.Context) {
		user, err := getUser(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "cannot get user in context", slog.String("error", err.Error()))
			ctx.SecureJSON(http.StatusInternalServerError, enging.InternalServerErrorMessage)
			return
		}
		// check scope.
		// Maybe this code should be in the 'scope' package, for example.
		if user.Scope != "profile:view" {
			slog.InfoContext(ctx, "user doesnot have 'profile:view' scope", slog.String("scope", user.Scope))
			ctx.SecureJSON(http.StatusBadRequest, enging.BadRequestMessage)
			return
		}
		profile, err := service.ViewUserProfile(ctx, user.UserId)
		if err != nil {
			slog.ErrorContext(ctx, "cannot view user profile", slog.String("error", err.Error()))
			ctx.SecureJSON(http.StatusInternalServerError, enging.InternalServerErrorMessage)
			return
		}
		var resp ProfileGetResponse
		resp.UserId = profile.Id
		resp.Name = profile.Name
		resp.Age = profile.Age
		resp.Profile = profile.Profile
		ctx.SecureJSON(http.StatusOK, resp)
	})
	return router
}
