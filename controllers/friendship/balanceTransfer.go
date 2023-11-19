package friendship

import (
	"github.com/simhonchourasia/betfr-be/dbinterface"
	"github.com/simhonchourasia/betfr-be/models"
)

func TransferBalance(db dbinterface.DBInterface, balanceTransfer models.BalanceTransfer) error {
	fship, err := getFriendshipBetweenUsers(db, balanceTransfer.LoserName, balanceTransfer.WinnerName)
	if err != nil {
		// Note that this will also error if no friendship exists between these users
		return err
	}

	// Note that balance is how much initiator owes receiver
	// so balance will increase if initiator loses
	if fship.InitiatorName == balanceTransfer.LoserName {
		fship.Balance += balanceTransfer.Amount
	} else if fship.InitiatorName == balanceTransfer.WinnerName {
		fship.Balance -= balanceTransfer.Amount
	}

	return db.Save(&fship).Error
}
