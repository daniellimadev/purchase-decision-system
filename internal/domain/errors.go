package domain

import "errors"

var (
	ErrPurchaseNotFound      = errors.New("purchase not found")
	ErrInvalidAmount         = errors.New("invalid purchase amount")
	ErrPurchaseAlreadyDecided = errors.New("purchase already has a decision")
)
