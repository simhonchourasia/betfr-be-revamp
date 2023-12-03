package models

import (
	"time"

	"github.com/google/uuid"
)

type Stake struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UnderlyingID   uuid.UUID // id of underlying bet that the stake is for
	OwnerName      string    // name of user who owns this stake
	SharesStaked   int64     // how many shares the owner wishes to stake
	SharesFilled   int64     // between 0 and SharesStaked
	BackingCreator bool      // whether this is backing up the creator or the receiver of the underlying bet
	Comment        string
	TimeCreated    time.Time // we fill earlier stakes for a given bet first
}
