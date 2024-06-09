package types

import "errors"

var (
	NotEnoughBalance = errors.New("Insufficient balance in funding account")
)
