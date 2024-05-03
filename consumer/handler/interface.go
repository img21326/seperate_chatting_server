package handler

import (
	requestmsg "chat_system/consumer/request_msg"
	"encoding/json"
)

type MessageTypeHandlerInterface interface {
	Perform(requestmsg.Msg) error
}

type RequestMsg struct {
	Msg string `json:"msg"`
}

type MsgType int

const (
	WaitingForPairing MsgType = iota
	InRoom
	NotInRoom
	OnMessage
	LeaveRoom
)

type ResponseMsg struct {
	MsgType MsgType `json:"msg_type"`
	Msg     string  `json:"msg"`
}

func NewResponseMsg(msgType MsgType, msg string) ResponseMsg {
	return ResponseMsg{
		MsgType: msgType,
		Msg:     msg,
	}
}

func (msg ResponseMsg) ToJSON() ([]byte, error) {
	return json.Marshal(msg)
}
