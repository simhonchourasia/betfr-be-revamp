package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/controllers/friendship"
)

func UnprotectedFriendshipRoutes(incomingRoutes *gin.Engine, handler friendship.FriendshipHandler) {
	incomingRoutes.GET("/friends/all", handler.GetAllFriendsFunc)
}

func ProtectedFriendshipRoutes(incomingRoutes *gin.Engine, handler friendship.FriendshipHandler) {
	incomingRoutes.POST("/friends/sendfriendreq", handler.SendFriendReqFunc)
	incomingRoutes.POST("/friends/handlefriendreq", handler.ResolveFriendReqFunc)
}
