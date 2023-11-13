package models

import (
	"time"

	"github.com/google/uuid"
)

type BetStatus int8
type BetOutcome int8

const (
	InvalidBet BetStatus = iota
	PendingBet
	DeclinedBet
	OngoingBet
	Resolved
	Conflicted
)

const (
	Undecided BetOutcome = iota
	CreatorWon
	ReceiverWon
)

type Bet struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatorID       uuid.UUID
	ReceiverID      uuid.UUID
	OverallStatus   BetStatus
	Outcome         BetOutcome
	CreatorAmount   int64 // amount that creator receives upon winning per share
	ReceiverAmount  int64 // amount that receiver receives upon winning per share
	NumShares       int64
	CreatorOutcome  BetOutcome // what the creator says the outcome was
	ReceiverOutcome BetOutcome // what the receiver says the outcome was
	Title           string
	Description     string
	ExpiryTime      time.Time
	CreatorStaked   int64 // shares staked in favour of creator
	ReceiverStaked  int64 // shares staked in favour of receiver
}