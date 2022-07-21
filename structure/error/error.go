package error

import "errors"

var (
	RecordNotFound  = errors.New("RecordNotFound")
	RoomIsClose     = errors.New("RoomIsClosed")
	PairNotSuccess  = errors.New("PairNotSuccess")
	ClientNotInHost = errors.New("ClientNotInHost")
)
