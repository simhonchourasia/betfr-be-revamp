package friendship

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/authentication"
	"github.com/simhonchourasia/betfr-be/controllers"
	"github.com/simhonchourasia/betfr-be/models"
)

type FriendshipHandler controllers.Handler

func (friendshipHandler *FriendshipHandler) sendFriendReqHelper(friendReq models.FriendReq) error {
	fship := createFriendshipReq(friendReq)
	return friendshipHandler.Db.Create(&fship).Error
}

// TODO: this could probably be refactored to a utils file and take in a DB interface for mocking purposes
func (friendshipHandler *FriendshipHandler) SendFriendReqFunc(c *gin.Context) {
	var friendReq models.FriendReq
	if err := c.BindJSON(&friendReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if permissionErr := authentication.CheckUserPermissions(c, friendReq.InitiatorName); permissionErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": permissionErr.Error()})
		return
	}

	// Check if there is already a friend request between the two
	currentStatus, err := getFriendStatusBetweenUsers(friendshipHandler.Db, friendReq.InitiatorName, friendReq.ReceiverName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if currentStatus == models.Pending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "There is already an friend request with this user"})
		return
	} else if currentStatus == models.Friends {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already friends with this user!"})
		return
	} else if currentStatus == models.NotConnected {
		// This is the interesting case
		if err := friendshipHandler.sendFriendReqHelper(friendReq); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Successfully sent friend request!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Bad friend request mode %d selected", currentStatus)})
	}

}

func (friendshipHandler *FriendshipHandler) ResolveFriendReqFunc(c *gin.Context) {
	var friendReqResolution models.FriendReqResolution
	if err := c.BindJSON(&friendReqResolution); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if permissionErr := authentication.CheckUserPermissions(c, friendReqResolution.ReceiverName); permissionErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": permissionErr.Error()})
		return
	}

	var fship models.Friendship
	err := friendshipHandler.Db.Where(
		"initiator_name = ? AND receiver_name  = ?",
		friendReqResolution.InitiatorName,
		friendReqResolution.ReceiverName,
	).First(&fship).Error

	if err != nil {
		// This should handle the case of a friendship relation not existing
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if friendReqResolution.ReqStatus == models.Accepted {
		fship.FriendStatus = models.Friends
		if err := friendshipHandler.Db.Save(&fship).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Friend request accepted!"})
	} else if friendReqResolution.ReqStatus == models.Declined {
		// Remove from table so they can try again later
		friendshipHandler.Db.Delete(&fship)
		c.JSON(http.StatusOK, gin.H{"msg": "Friend request declined..."})
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "Friend request untouched."})
	}
}

// TODO: make this get 50 by 50 or something like that
// Gets all friends of given user
// This will just return an empty list if a user doesn't exist
func (friendshipHandler *FriendshipHandler) GetAllFriendsFunc(c *gin.Context) {
	// TODO: surely there is a cleaner way to do this
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	currentUser := user.Username

	log.Printf("getting all friends for user %s", currentUser)
	var friendships []models.Friendship
	err := friendshipHandler.Db.Model(&models.Friendship{}).
		Select("DISTINCT friendships.*").
		Joins("JOIN users ON (friendships.receiver_name = users.username OR friendships.initiator_name = users.username)").
		Where(
			"(friendships.initiator_name = ? OR friendships.receiver_name = ?) AND friendships.friend_status = ?",
			currentUser, currentUser, models.Friends).
		Find(&friendships).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, friendships)
}
