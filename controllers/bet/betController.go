package bet

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/authentication"
	"github.com/simhonchourasia/betfr-be/controllers"
	"github.com/simhonchourasia/betfr-be/models"
	"gorm.io/gorm"
)

type BetHandler controllers.Handler

// Pass in creator name, receiver name, creator amount, receiver amount, shares, underlying, title, description, expiry date
// TODO: put an example request here
func (betHandler *BetHandler) CreateBetReqFunc(c *gin.Context) {
	var bet models.Bet
	if err := c.BindJSON(&bet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check that the user sending the bet request is the one logged in
	if permissionErr := authentication.CheckUserPermissions(c, bet.CreatorName); permissionErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": permissionErr.Error()})
		return
	}

	// Initialize the rest of the fields
	bet.OverallStatus = models.PendingBet
	bet.Outcome = models.Undecided
	bet.CreatorOutcome = models.Undecided
	bet.ReceiverOutcome = models.Undecided
	bet.CreatedTime = time.Now()
	bet.NumStakesFilled = 0
	bet.CreatorStakedUnfilled = 0
	bet.ReceiverStakedUnfilled = 0

	// TODO: this code could probably be moved elsewhere
	if bet.ExpiryTime.Before(bet.CreatedTime.Add(5 * time.Minute)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bets cannot be created with less than 5 minutes to expiry upon creation"})
		return
	}

	if err := betHandler.Db.Create(&bet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Something went wrong in creating bet: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Created bet: '%s'", bet.ID))
}

func (betHandler *BetHandler) HandleBetReqFunc(c *gin.Context) {
	var betReqHandle models.BetReqHandle
	if err := c.BindJSON(&betReqHandle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get bet from database (and also check that it exists)
	var bet models.Bet
	if err := betHandler.Db.First(&bet, betReqHandle.BetID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bet does not exist"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	// Check that the user accepting/declining the bet request is the one logged in
	if permissionErr := authentication.CheckUserPermissions(c, bet.ReceiverName); permissionErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": permissionErr.Error()})
		log.Printf("Bet receiver %s is not the one logged in", bet.ReceiverName)
		return
	}

	if bet.OverallStatus != models.PendingBet {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bet exists but is not pending"})
		return
	}

	if betReqHandle.ReqStatus == models.Accepted {
		bet.OverallStatus = models.OngoingBet
		c.JSON(http.StatusOK, gin.H{"msg": "Bet request accepted!"})
	} else if betReqHandle.ReqStatus == models.Declined {
		bet.OverallStatus = models.DeclinedBet
		c.JSON(http.StatusOK, gin.H{"msg": "Bet request declined..."})
	}

	if err := betHandler.Db.Save(&bet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (betHandler *BetHandler) ResolveBetFunc(c *gin.Context) {
	var betResolve models.BetResolve
	if err := c.BindJSON(&betResolve); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check that the user modifying the bet request is the one logged in
	if permissionErr := authentication.CheckUserPermissions(c, betResolve.Username); permissionErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": permissionErr.Error()})
		log.Printf("Bet request modifier %s is not the one logged in", betResolve.Username)
		return
	}

	// TODO: this can be factored out
	var bet models.Bet
	if err := betHandler.Db.First(&bet, betResolve.BetID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bet does not exist"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	creatorPermissible := authentication.CheckUserPermissions(c, bet.CreatorName)
	receiverPermissible := authentication.CheckUserPermissions(c, bet.ReceiverName)
	if creatorPermissible != nil && receiverPermissible != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": creatorPermissible.Error()})
		fmt.Printf("logged in user is not one of the bettors\n")
		return
	}

	if bet.CreatorName != betResolve.Username && bet.ReceiverName != betResolve.Username {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only creator or receiver can provide a resolve update"})
		return
	}

	if bet.OverallStatus == models.DeclinedBet {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bet is not active (has been declined)"})
		return
	} else if bet.OverallStatus == models.Resolved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bet is not active (has been resolved)"})
		return
	}

	// Assume then that the bet is ongoing or conflicted
	// In any case, update the CreatorOutcome/ReceiverOutcome,
	// and check if both sides have provided an outcome
	bothStatusDecided := false
	if betResolve.Username == bet.CreatorName {
		bet.CreatorOutcome = betResolve.StatedOutcome
		if bet.ReceiverOutcome != models.Undecided {
			bothStatusDecided = true
		}
	} else if betResolve.Username == bet.ReceiverName {
		bet.ReceiverOutcome = betResolve.StatedOutcome
		if bet.CreatorOutcome != models.Undecided {
			bothStatusDecided = true
		}
	}

	// Update final outcome if both sides provided an outcome
	if bothStatusDecided {
		if bet.CreatorOutcome == bet.ReceiverOutcome {
			bet.Outcome = bet.CreatorOutcome
			bet.OverallStatus = models.Resolved
			if err := PayoutBet(betHandler.Db, &bet); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			bet.Outcome = models.HasConflict
			bet.OverallStatus = models.Conflicted
		}
	}

	betHandler.Db.Save(&bet)
	c.JSON(http.StatusOK, gin.H{"msg": "Updated bet!"})
}

func (betHandler *BetHandler) GetAllBetsFunc(c *gin.Context) {
	claims, statusCode, err := authentication.GetClaimsFromCookie(c)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	var getBetsForUser models.GetBetsForUser
	if err := c.BindJSON(&getBetsForUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if claims.Issuer != getBetsForUser.Username {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cookie issuer does not match request"})
		return
	}

	bets, err := GetBetsForUserWithStatus(
		betHandler.Db,
		getBetsForUser.Username,
		getBetsForUser.DesiredStatus,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bets)
}
