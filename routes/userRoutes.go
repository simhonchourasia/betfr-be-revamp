package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/controllers/user"
)

func UnprotectedUserRoutes(incomingRoutes *gin.Engine, handler user.UserHandler) {
	incomingRoutes.POST("/users/signup", handler.SignUpFunc)
	incomingRoutes.POST("/users/login", handler.LoginFunc)
	incomingRoutes.GET("/users/get", handler.GetUserFunc)
	incomingRoutes.POST("/users/logout", handler.LogoutFunc)
}

func ProtectedUserRoutes(incomingRoutes *gin.Engine) {
}
