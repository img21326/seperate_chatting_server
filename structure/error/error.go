package error

import "errors"

var (
	RecordNotFound = errors.New("RecordNotFound")
	RoomIsClose    = errors.New("RoomIsClosed")
)
