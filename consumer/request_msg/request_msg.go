package requestmsg

import (
	"chat_system/websocket"
	"encoding/json"
)

type RequestMsgEvent int

const (
	JoinRoom RequestMsgEvent = iota
	LeaveRoom
	Typing
	OnMessage
)

type Msg struct {
	MsgEvent RequestMsgEvent `json:"msg_event"`
	Msg      string          `json:"msg"`
	Sender   string
}

func NewRequestMsg(msg websocket.Msg) (Msg, error) {
	var requestMsg Msg
	err := json.Unmarshal(msg.Msg, &requestMsg)
	if err != nil {
		return Msg{}, err
	}
	requestMsg.Sender = msg.Sender
	return requestMsg, nil
}
