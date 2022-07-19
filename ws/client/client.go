package client

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/img21326/fb_chat/structure/message"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/img21326/fb_chat/structure/ws"
)

type Client struct {
	Conn       *websocket.Conn
	Ctx        context.Context
	CtxCancel  context.CancelFunc
	Send       chan []byte
	User       user.User
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

func (c *Client) ReadPump(PublishChan chan *pubmessage.PublishMessage) {

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		select {
		case <-c.Ctx.Done():
			log.Printf("[websocket client] stop listen user: %v", c.User.ID)
			return
		default:
			_, messageByte, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("websocket unexcept error: %v", err)
				}
				c.CtxCancel()
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
				messageModel := message.Message{
					RoomId:  c.RoomId,
					UserId:  c.User.ID,
					Message: getMessage.Message,
					Time:    getMessage.Time.Time,
				}
				publishMessage := pubmessage.PublishMessage{
					Type:     "message",
					SendFrom: c.User.ID,
					SendTo:   c.PairId,
					Payload:  messageModel,
				}
				PublishChan <- &publishMessage
				continue
			}
			if getMessage.Type == "leave" {
				log.Printf("[websocket client] get leave message from user: %v", c.User.ID)
				publishMessage := pubmessage.PublishMessage{
					Type:     "leave",
					SendFrom: c.User.ID,
					SendTo:   c.PairId,
					Payload:  c.RoomId,
				}
				PublishChan <- &publishMessage
				return
			}
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
		case <-c.Ctx.Done():
			return
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
