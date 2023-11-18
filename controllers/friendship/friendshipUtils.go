package friendship

import "github.com/simhonchourasia/betfr-be/models"

func createFriendshipReq(friendReq models.FriendReq) models.Friendship {
	var fship models.Friendship
	fship.InitiatorName = friendReq.InitiatorName
	fship.ReceiverName = friendReq.ReceiverName
	fship.FriendStatus = models.Pending
	fship.Balance = 0
	return fship
}
