package websocket

type MsgType int

const (
	MsgTypeSpecificUser MsgType = iota
	MsgTypeBroadcast
)

type MsgFrom int

const (
	MsgFromLocal MsgFrom = iota
	MsgFromBroker
)

type Msg struct {
	Type     MsgType
	Sender   string
	Receiver string
	Msg      []byte
	From     MsgFrom
}
