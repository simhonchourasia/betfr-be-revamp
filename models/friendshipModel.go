package models

import (
	"github.com/google/uuid"
)

type FriendshipStatus int8

const (
	NotConnected FriendshipStatus = iota
	Pending
	Declined
	Friends
)

type Friendship struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	InitiatorID   uuid.UUID
	ReceiverID    uuid.UUID
	InitiatorName string
	ReceiverName  string
	FriendStatus  FriendshipStatus
	Balance       int64 // how much initiator owes receiver
}
