package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/tarm/serial"
	"log"
)

func Config(app *fiber.App, serial_port *serial.Port) {

	//for {
	//	n, err = serial_port.Read(buf)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	log.Printf( "serial:  %q", buf[:n])
	//}
	//

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

				for i, item := range clients {
					if item == c {
						clients = append(clients[:i], clients[i+1:]...)
					}
				}

				break
			}
			log.Printf("websocket: %s", msg)

			for _, client := range clients {
				if err = client.WriteMessage(mt, msg); err != nil {
					log.Println("write:", err)
					break
				}
			}

		}

	}))

	app.Get("/ws/check", websocket.New(func(c *websocket.Conn) {
		//var (
		//	mt  int
		//	msg []byte
		//	err error
		//)
		for {

			//

		}

	}))

}
