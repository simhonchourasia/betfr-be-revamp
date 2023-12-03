package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/controllers/stake"
)

func UnprotectedStakeRoutes(incomingRoutes *gin.Engine, handler stake.StakeHandler) {

}

func ProtectedStakeRoutes(incomingRoutes *gin.Engine, handler stake.StakeHandler) {
	incomingRoutes.POST("/stakes/createstake", handler.CreateStakeFunc)
}
