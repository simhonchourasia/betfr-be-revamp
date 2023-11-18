package friendship

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/simhonchourasia/betfr-be/authentication"
	"github.com/simhonchourasia/betfr-be/controllers"
	"github.com/simhonchourasia/betfr-be/models"
)

type FriendshipHandler controllers.Handler

// TODO: this could probably be refactored to a utils file and take in a DB interface for mocking purposes
// func (friendshipHandler *FriendshipHandler)

func (friendshipHandler *FriendshipHandler) SendFriendReqFunc(c *gin.Context) {
	// We only need/use the fields InitiatorID, InitiatorName and ReceiverName
	var friendReq models.Friendship

	if err := c.BindJSON(&friendReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if permissionErr := authentication.CheckUserPermissions(c, friendReq.InitiatorName); permissionErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": permissionErr.Error()})
		return
	}

}

func (friendshipHandler *FriendshipHandler) ResolveFriendReqFunc(c *gin.Context) {

}
