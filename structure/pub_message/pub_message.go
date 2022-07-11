package pubmessage

type PublishMessage struct {
	Type     string      `json:"type"`
	SendFrom uint        `json:"sendFrom"`
	SendTo   uint        `json:"sendTo"`
	Payload  interface{} `json:"payload"`
}

type SendToUserMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
