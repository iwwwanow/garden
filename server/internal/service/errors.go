package service

import "errors"

var (
	ErrUsernameTaken       = errors.New("username already taken")
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrNotFound            = errors.New("not found")
	ErrAlreadyWatered      = errors.New("already watered today")
	ErrFlowerDried         = errors.New("flower is dried")
	ErrMaxFlowers          = errors.New("max 64 active flowers reached")
	ErrInsufficientSeeds   = errors.New("insufficient seeds")
	ErrPlantLimitReached   = errors.New("only 1 flower can be planted per day")
)
