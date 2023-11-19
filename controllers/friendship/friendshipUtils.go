package friendship

import (
	"errors"

	"github.com/simhonchourasia/betfr-be/dbinterface"
	"github.com/simhonchourasia/betfr-be/models"
	"gorm.io/gorm"
)

func createFriendshipReq(friendReq models.FriendReq) models.Friendship {
	var fship models.Friendship
	fship.InitiatorName = friendReq.InitiatorName
	fship.ReceiverName = friendReq.ReceiverName
	fship.FriendStatus = models.Pending
	fship.Balance = 0
	return fship
}

func getFriendshipBetweenUsers(db dbinterface.DBInterface, username1 string, username2 string) (models.Friendship, error) {
	var fship models.Friendship
	err := db.Where(
		"(initiator_name = ? AND receiver_name  = ?) OR (initiator_name = ? AND receiver_name  = ?)",
		username1, username2, username2, username1,
	).First(&fship).Error

	return fship, err
}

func getFriendStatusBetweenUsers(db dbinterface.DBInterface, username1 string, username2 string) (models.FriendshipStatus, error) {
	fship, err := getFriendshipBetweenUsers(db, username1, username2)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NotConnected, nil
		} else {
			return models.NotConnected, err
		}
	}
	return fship.FriendStatus, nil
}
