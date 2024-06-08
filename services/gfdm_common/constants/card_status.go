package constants

type CardStatus int

const (
	CardStatusError CardStatus = iota
	CardStatusNew
	CardStatusSuccess
	CardStatusExpired
	CardStatusMaintenance
	CardStatusEndOfLife
)
