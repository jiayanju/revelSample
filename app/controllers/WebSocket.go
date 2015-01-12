package controllers

import (
	"code.google.com/p/go.net/websocket"
	"github.com/revel/revel"
	"gopkg.in/stomp.v1"
)

type WebSocket struct {
	*revel.Controller
}

func (c WebSocket) Room(user string) revel.Result {
	return c.Render(user)
}

func (c WebSocket) RoomSocket(user string, ws *websocket.Conn) revel.Result {
	room := "/topic/room1"
	conn, _ := stomp.Dial("tcp", "localhost:61613", stomp.Options{})
	sub, _ := conn.Subscribe(room, stomp.AckAuto)

	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(newMessages)
				return
			}
			newMessages <- msg
		}
	}()

	for {
		select {
		case msg := <-sub.C:
			if websocket.JSON.Send(ws, string(msg.Body)) != nil {
				return nil
			}

		case receivedMsg, ok := <-newMessages:
			if !ok {
				return nil
			}
			conn.Send(room, "", []byte(receivedMsg), nil)

		}
	}

}
