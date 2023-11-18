package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/controllers/friendship"
)

func UnprotectedFriendshipRoutes(incomingRoutes *gin.Engine, handler friendship.FriendshipHandler) {

}

func ProtectedFriendshipRoutes(incomingRoutes *gin.Engine, handler friendship.FriendshipHandler) {
	incomingRoutes.POST("/users/sendfriendreq", handler.SendFriendReqFunc)
	incomingRoutes.POST("/users/handlefriendreq", handler.ResolveFriendReqFunc)
}
