package enging

import "github.com/gin-gonic/gin"

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
	ForbiddenErrorMessage = gin.H{
		"status": "Forbidden",
	}
	StatusUnauthorizedErrorMessage = gin.H{
		"status": "unauthorized",
	}
)
