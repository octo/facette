package backend

import (
	"errors"
)

var (
	ErrInvalidIdentifier = errors.New("invalid item identifier")
	ErrInvalidSlice      = errors.New("invalid slice")
	ErrInvalidStruct     = errors.New("invalid table struct")
	ErrNotAddressable    = errors.New("value is not addressable")
	ErrNotTransformable  = errors.New("value is not transformable")
	ErrUnknownColumn     = errors.New("unknown column")
)
