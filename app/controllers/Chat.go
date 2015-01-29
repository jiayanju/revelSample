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
	revel.TRACE.Println("Room Name : " + roomName)
	c.Validation.Required(roomName)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Home)
	}

	app.ChatService.SaveValueToList(ROOM_LIST_KEY, roomName)

	return c.Redirect("/chat/room/%s", roomName)
}

func (c *ChatController) Room(roomName string) revel.Result {
	revel.TRACE.Println("Room Name In Room Method : " + roomName)
	return c.Render(roomName)
}
