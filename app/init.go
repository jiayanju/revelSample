package app

import (
	"chat/app/chatservice"
	"github.com/revel/revel"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"gopkg.in/stomp.v1"
	"net/http"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
	revel.OnAppStart(initSystem)
	revel.OnAppStart(installHandlers)
}

var (
	service chatservice.ChatRedisService
)

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func initSystem() {
	service = chatservice.NewChatRedisService()

}

func installHandlers() {

	websocketHandler := sockjs.NewHandler("/websocket/sockjs/room", sockjs.DefaultOptions, func(session sockjs.Session) {
		revel.TRACE.Println("begin handle websocket")
		room := "/topic/room1"
		conn, _ := stomp.Dial("tcp", "localhost:61613", stomp.Options{})
		sub, _ := conn.Subscribe(room, stomp.AckAuto)

		newMessages := make(chan string)
		go func() {
			for {
				msg, err := session.Recv()
				if err != nil {
					close(newMessages)
					return
				}
				newMessages <- msg
			}
		}()

		msgs := service.GetValues("room2")

		count := len(msgs)
		revel.TRACE.Printf("msg length : %d", count)
		if count > 0 {
			for i := range msgs {
				session.Send(msgs[i])
			}
		}

		for {
			select {
			case msg := <-sub.C:
				if session.Send(string(msg.Body)) != nil {
					break
				}

			case receivedMsg, ok := <-newMessages:
				if !ok {
					break
				}
				conn.Send(room, "", []byte(receivedMsg), nil)
				service.SaveValueToList("room2", receivedMsg)
			}
		}

	})

	var (
		serverMux    = http.NewServeMux()
		revelHandler = revel.Server.Handler
	)
	serverMux.Handle("/websocket/sockjs/room/", websocketHandler)
	serverMux.Handle("/", revelHandler)
	revel.Server.Handler = serverMux
	revel.TRACE.Println("Register websocket handler")
}
