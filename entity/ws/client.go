package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/user"
)

type Client struct {
	Conn       *websocket.Conn
	Send       chan []byte
	User       user.UserModel
	WantToFind string
	PairClient *Client
	RoomId     uuid.UUID
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

func (c *Client) ReadPump(saveMessageChan chan<- *message.MessageModel, closeChan chan<- uuid.UUID, deleteChan chan<- *Client) {
	defer func() {
		deleteChan <- c
		c.Conn.Close()
		c.PairClient.PairClient = nil
		close(c.Send)
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, messageByte, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var getMessage Message
		err = json.Unmarshal(messageByte, &getMessage)
		if err != nil {
			log.Printf("error decode json: %v", err)
			return
		}

		if c.RoomId == uuid.Nil {
			return
		}
		if getMessage.Type == "message" {
			if c.PairClient != nil {
				c.PairClient.Send <- []byte(messageByte)
			}
			messageModel := &message.MessageModel{
				RoomId:  c.RoomId,
				UserId:  c.User.ID,
				Message: getMessage.Message,
			}
			saveMessageChan <- messageModel
			return
		}
		if getMessage.Type == "leave" {
			if c.PairClient != nil {
				c.PairClient.Send <- []byte(messageByte)
			}
			c.PairClient.RoomId = uuid.Nil
			c.PairClient.PairClient = nil
			c.PairClient.Conn.Close()
			c.RoomId = uuid.Nil
			c.PairClient = nil
			c.Conn.Close()
			closeChan <- c.RoomId
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
