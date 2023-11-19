package models

import (
	"github.com/google/uuid"
)

type FriendshipStatus int8

const (
	NotConnected FriendshipStatus = iota
	Pending
	Friends
)

type Friendship struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	InitiatorName string
	ReceiverName  string
	FriendStatus  FriendshipStatus
	Balance       int64 // how much initiator owes receiver
}

type FriendReq struct {
	InitiatorName string
	ReceiverName  string
}

type FriendReqResolution struct {
	InitiatorName string
	ReceiverName  string
	ReqStatus     RequestStatus
}

type BalanceTransfer struct {
	LoserName  string
	WinnerName string
	Amount     int64
}
