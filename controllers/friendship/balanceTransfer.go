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

func TransferBalancesBulk(db dbinterface.DBInterface, balanceTransfers []models.BalanceTransfer) error {
	// This needs to match up by index with the given balance transfer array
	friendships := make([]models.Friendship, len(balanceTransfers))
	for _, transfer := range balanceTransfers {
		fship, err := getFriendshipBetweenUsers(db, transfer.LoserName, transfer.WinnerName)
		if err != nil {
			// Note that this will also error if no friendship exists between these users
			return err
		}
		friendships = append(friendships, fship)
	}

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for i, fship := range friendships {
		balanceTransfer := balanceTransfers[i]
		// Note that balance is how much initiator owes receiver
		// so balance will increase if initiator loses
		if fship.InitiatorName == balanceTransfer.LoserName {
			fship.Balance += balanceTransfer.Amount
		} else if fship.InitiatorName == balanceTransfer.WinnerName {
			fship.Balance -= balanceTransfer.Amount
		}

		if err := tx.Save(&fship).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
