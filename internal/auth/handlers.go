package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/clients/:client_id", func(ctx *gin.Context) {
		var req ServiceClientGetRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		var resp ServiceClientGetResponse
		ctx.SecureJSON(http.StatusOK, resp)
	})

	v1.POST("/authentication", func(ctx *gin.Context) {
		var req AuthenticationRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		var resp AuthenticationResponse
		ctx.SecureJSON(http.StatusOK, resp)
	})

	v1.POST("/authorization", func(ctx *gin.Context) {
		var req AuthorizationRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
