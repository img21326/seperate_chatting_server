package websocket

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientInterface interface {
	ID() string
	SendMsg(msg []byte)
	Close()
}

type Client struct {
	id   string
	ctx  context.Context
	conn *websocket.Conn
	hub  HubInterface

	sendChan  chan []byte
	closeOnce sync.Once
}

func NewClient(conn *websocket.Conn, opts ...func(*Client)) ClientInterface {
	client := &Client{
		id:       uuid.New().String(),
		conn:     conn,
		sendChan: make(chan []byte, 256),
	}

	for _, opt := range opts {
		opt(client)
	}

	client.run()
	client.hub.AddClient(client)
	return client
}

func (c *Client) WithID(id string) {
	c.id = id
}

func (c *Client) WithHub(hub HubInterface) {
	c.hub = hub
}

func (c *Client) WithContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) SendMsg(msg []byte) {
	c.sendChan <- msg
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.sendChan)
		c.ctx.Done()
		c.conn.Close()
		c.hub.RemoveClient(c.id)
	})
}

func (c *Client) readLoop() {
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				return
			}
			c.hub.ReceiveMsgFromLocalClient(Msg{
				Sender: c.id,
				Msg:    msg,
				From:   MsgFromLocal,
			})
		}
	}
}

func (c *Client) run() {
	go c.writeLoop()
	go c.readLoop()
}

func (c *Client) writeLoop() {
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg, ok := <-c.sendChan:
			if !ok {
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		}
	}
}
