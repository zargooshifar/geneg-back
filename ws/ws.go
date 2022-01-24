package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

func Config(app *fiber.App) {

	clients := []*websocket.Conn{}

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/buffet", websocket.New(func(c *websocket.Conn) {
		clients = append(clients, c)

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)

				break
			}

			for i, client := range clients {

				if client == nil {
					clients = append(clients[:i], clients[i+1:]...)
					return
				} else {
					if err = client.WriteMessage(mt, msg); err != nil {
						log.Println("write:", err)
						break
					}
				}

			}

		}

	}))

}
