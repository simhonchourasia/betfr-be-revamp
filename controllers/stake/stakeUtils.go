package stake

import (
	"log"
	"sort"

	"github.com/simhonchourasia/betfr-be/controllers/friendship"
	"github.com/simhonchourasia/betfr-be/dbinterface"
	"github.com/simhonchourasia/betfr-be/models"
)

// Provides a list of all stakes for a bet in ascending order of creation
// i.e. stakes created earlier are earlier in the list
// Returns two lists, one for stakes backing creator and one for stakes backing receiver
// If option `onlyUnfilled` is true, only return unfilled stakes
func getAllStakes(
	db dbinterface.DBInterface,
	bet *models.Bet,
	onlyUnfilled bool) (creatorStakes []models.Stake, receiverStakes []models.Stake, dbErr error) {

	onlyGetUnfilledStr := ""
	if onlyUnfilled {
		onlyGetUnfilledStr = " AND stakes.shares_staked > stakes.shares_filled"
	}

	var allBets []models.Stake
	dbErr = db.Model(&models.Stake{}).
		Select("DISTINCT stakes.*").
		Where("stakes.underlying_id = ?"+onlyGetUnfilledStr, bet.ID).
		Order("stakes.time_created ASC").
		Find(&allBets).Error

	for _, bet := range allBets {
		if bet.BackingCreator {
			creatorStakes = append(creatorStakes, bet)
		} else {
			receiverStakes = append(receiverStakes, bet)
		}
	}

	return
}

// Given a new stake, a bet, and a list of existingstakes against it, try to fill as much as possible
// and return the resulting modified stakes, including the new stake
// Note that this modifies the passed in bet and new stake (as they are passed in by reference)
func fillStakes(
	newStake *models.Stake,
	stakesAgainstNew []models.Stake,
	underlyingBet *models.Bet) (modifiedStakes []models.Stake) {
	// stakesAgainstNew should already be sorted by time, but we sort again just in case
	sort.Slice(stakesAgainstNew, func(i, j int) bool {
		return stakesAgainstNew[i].TimeCreated.Before(stakesAgainstNew[j].TimeCreated)
	})

	log.Printf("stakes against new: %v, new stake: %v", stakesAgainstNew, *newStake)

	// At the end, our modified stakes will be stakesAgainstNew[:filledIdx]
	filledIdx := 0
	for i, oppStake := range stakesAgainstNew {
		// note that oppStake is the stake we are trying to match up with our current stake
		// ex. if our new stake is backing up stake creator, the oppStakes are backing up receiver
		remainder := newStake.SharesStaked - newStake.SharesFilled
		remainingInOppStake := oppStake.SharesStaked - oppStake.SharesFilled
		if remainingInOppStake > remainder {
			// Case 1: We can only fill the opposing stake partially
			stakesAgainstNew[i].SharesFilled += remainder
			newStake.SharesFilled += remainder
			filledIdx++
			break
		} else { // remainingInOppStake <= remainder
			// Case 2: We can fill the opposing stake completely
			stakesAgainstNew[i].SharesFilled += remainingInOppStake
			newStake.SharesFilled += remainingInOppStake
			filledIdx++
		}
	}
	modifiedStakes = stakesAgainstNew[:filledIdx]

	if newStake.BackingCreator {
		underlyingBet.NumStakesFilled += newStake.SharesFilled
		underlyingBet.ReceiverStakedUnfilled -= newStake.SharesFilled
	} else {
		underlyingBet.NumStakesFilled += newStake.SharesFilled
		underlyingBet.CreatorStakedUnfilled -= newStake.SharesFilled
	}

	return
}

// Given a bet and stakes for the bet and a stake that is to be added,
// fill up stakes as much as possible and return the resulting stakes
// Note that this modifies the passed in bet and new stake (as they are passed in by reference)
func handleStakes(
	creatorUnfilledStakes []models.Stake,
	receiverUnfilledStakes []models.Stake,
	newStake *models.Stake,
	underlyingBet *models.Bet) []models.Stake {

	var modifiedStakes []models.Stake

	if newStake.BackingCreator {
		if len(receiverUnfilledStakes) == 0 {
			// Can't fill at all
			underlyingBet.CreatorStakedUnfilled += newStake.SharesStaked
		} else {
			// we actually have to try and fill the receiver stakes
			modifiedStakes = fillStakes(newStake, receiverUnfilledStakes, underlyingBet)
		}
	} else {
		if len(creatorUnfilledStakes) == 0 {
			// Can't fill at all
			underlyingBet.ReceiverStakedUnfilled += newStake.SharesStaked
		} else {
			// we actually have to try and fill the creator stakes
			modifiedStakes = fillStakes(newStake, creatorUnfilledStakes, underlyingBet)
		}
	}

	return modifiedStakes
}

func getTransfersForStakes(creatorStakes []models.Stake, receiverStakes []models.Stake) []models.BalanceTransfer {
	var transfers []models.BalanceTransfer
	// do two-pointers on creator stakes and receiver stakes
	creatorIdx := 0
	receiverIdx := 0
	var currCreatorRemainder int64 = 0
	for creatorIdx < len(creatorStakes) && receiverIdx < len(receiverStakes) {
		if currCreatorRemainder == 0 {
			currCreatorRemainder += creatorStakes[creatorIdx].SharesFilled
			creatorIdx++
			continue
		}
		for receiverIdx < len(receiverStakes) {
			// TODO
		}
	}

	return transfers
}

func PayoutStakes(db dbinterface.DBInterface, bet *models.Bet) error {
	creatorStakes, receiverStakes, err := getAllStakes(db, bet, false)
	if err != nil {
		return err
	}

	transfers := getTransfersForStakes(creatorStakes, receiverStakes)
	// TODO: do the transactions
	friendship.TransferBalancesBulk(db, transfers)

	return nil
}
