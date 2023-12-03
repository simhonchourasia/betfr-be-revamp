package stake

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/controllers"
	"github.com/simhonchourasia/betfr-be/controllers/bet"
	"github.com/simhonchourasia/betfr-be/models"
)

type StakeHandler controllers.Handler

func (stakeHandler *StakeHandler) CreateStakeFunc(c *gin.Context) {
	var stake models.Stake
	if err := c.BindJSON(&stake); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	underlyingBet, err := bet.GetBet(stakeHandler.Db, stake.UnderlyingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error getting underlying for stake: %s", err.Error())})
		return
	}

	stake.SharesFilled = 0
	stake.TimeCreated = time.Now()

	creatorUnfilledStakes, receiverUnfilledStakes, err := getAllStakes(stakeHandler.Db, &underlyingBet, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	modifiedStakes := handleStakes(creatorUnfilledStakes, receiverUnfilledStakes, &stake, &underlyingBet)
	// The way handleStakes is implemented will not return the new stake in modifiedStakes, so we need to save it separately

	for _, modifiedStake := range modifiedStakes {
		if err := stakeHandler.Db.Save(&modifiedStake).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if err := stakeHandler.Db.Save(&stake).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := stakeHandler.Db.Save(&underlyingBet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("Successfully created stake %s", stake.ID)})
}
