package models

type RequestStatus int8

const (
	NoAction RequestStatus = iota
	Accepted
	Declined
)
