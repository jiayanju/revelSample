package controllers

import (
	"chat/app"
	"github.com/revel/revel"
)

const (
	ROOM_LIST_KEY = "RoomList"
)

type ChatController struct {
	*revel.Controller
}

func (c *ChatController) CreateRoom() revel.Result {
	roomName := c.Params.Get("roomName")
	revel.TRACE.Println("Room Name" + roomName)
	app.ChatService.SaveValueToList(ROOM_LIST_KEY, roomName)
	return c.Render()
}
