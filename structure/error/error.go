package error

import "errors"

var (
	RecordNotFound        = errors.New("RecordNotFound")
	RoomIsClose           = errors.New("RoomIsClosed")
	PairNotSuccess        = errors.New("PairNotSuccess")
	ClientNotInHost       = errors.New("ClientNotInHost")
	QueueSmallerThan1     = errors.New("QueueSmallerThan1")
	UserNotInThisRoom     = errors.New("UserNotInThisRoom")
	ChannelClosed         = errors.New("ChannelClosed")
	ChannelNoMsg          = errors.New("ChannelNoMsg")
	InterfaceConvertError = errors.New("InterfaceConvertError")
)
