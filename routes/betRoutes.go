package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/controllers/bet"
)

func UnprotectedBetRoutes(incomingRoutes *gin.Engine, handler bet.BetHandler) {
	incomingRoutes.POST("/bets/userbets", handler.GetAllBetsFunc)
}

func ProtectedBetRoutes(incomingRoutes *gin.Engine, handler bet.BetHandler) {
	incomingRoutes.POST("/bets/createbetreq", handler.CreateBetReqFunc)
	incomingRoutes.POST("/bets/handlebetreq", handler.HandleBetReqFunc)
	incomingRoutes.POST("/bets/resolvebet", handler.ResolveBetFunc)
}
