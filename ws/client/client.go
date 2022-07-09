package client

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/user"
	"github.com/img21326/fb_chat/ws"
	"github.com/img21326/fb_chat/ws/messageType"
)

type Client struct {
	Conn       *websocket.Conn
	Send       chan []byte
	User       user.UserModel
	WantToFind string
	RoomId     uuid.UUID
	PairId     uint
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (c *Client) ReadPump(PublishChan chan<- messageType.PublishMessage, unRegisterChan chan<- *Client, deletePairChan chan<- *Client) {
	defer func() {
		unRegisterChan <- c
		deletePairChan <- c
		c.Conn.Close()
		close(c.Send)
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, messageByte, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket unexcept error: %v", err)
			}
			break
		}
		var getMessage ws.WebsocketMessage
		err = json.Unmarshal(messageByte, &getMessage)
		if err != nil {
			log.Printf("error decode json: %v", err)
			continue
		}
		log.Printf("[websocket client] get message from user: %v, message: %+v", c.User.ID, getMessage)

		// 沒有在房間裡 不做任何動作
		if c.RoomId == uuid.Nil {
			continue
		}
		if getMessage.Type == "message" {
			messageModel := message.MessageModel{
				RoomId:  c.RoomId,
				UserId:  c.User.ID,
				Message: getMessage.Message,
				Time:    time.Time(getMessage.Time),
			}
			publishMessage := messageType.PublishMessage{
				Type:     "message",
				SendFrom: c.User.ID,
				SendTo:   c.PairId,
				Payload:  messageModel,
			}
			PublishChan <- publishMessage
			continue
		}
		if getMessage.Type == "leave" {
			publishMessage := messageType.PublishMessage{
				Type:     "leave",
				SendFrom: c.User.ID,
				SendTo:   c.PairId,
				Payload:  c.RoomId,
			}
			PublishChan <- publishMessage
			c.RoomId = uuid.Nil
			c.PairId = 0
			c.Conn.Close()
			break
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
