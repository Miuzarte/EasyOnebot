package EasyOnebot

import "errors"

var (
	ErrNoWsConn     = errors.New("no ws connection")
	ErrEmptyData    = errors.New("empty data")
	ErrNoForwardSeg = errors.New("no forward segment")
	ErrNoForwardId  = errors.New("no forward id")
	ErrInvalidId    = errors.New("invalid id")
)
