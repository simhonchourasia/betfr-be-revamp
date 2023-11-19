package bet

import (
	"fmt"

	"github.com/simhonchourasia/betfr-be/controllers/friendship"
	"github.com/simhonchourasia/betfr-be/dbinterface"
	"github.com/simhonchourasia/betfr-be/models"
)

// TODO: add stakes handling
func PayoutBet(db dbinterface.DBInterface, bet *models.Bet) error {
	if bet.OverallStatus != models.Resolved {
		return fmt.Errorf("cannot pay out bet that is not resolved")
	}
	var loserName string
	var winnerName string
	var amount int64
	if bet.Outcome == models.CreatorWon {
		loserName = bet.ReceiverName
		winnerName = bet.CreatorName
		amount = bet.CreatorAmount * bet.NumShares
	} else if bet.Outcome == models.ReceiverWon {
		loserName = bet.CreatorName
		winnerName = bet.ReceiverName
		amount = bet.ReceiverAmount * bet.NumShares
	}

	var balanceTransfer = models.BalanceTransfer{LoserName: loserName, WinnerName: winnerName, Amount: amount}

	return friendship.TransferBalance(db, balanceTransfer)
}
