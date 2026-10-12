package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrNotCancelable = errors.New("order cannot be cancelled in its current status")
	ErrUnknownStatus = errors.New("unknown order status")
)
