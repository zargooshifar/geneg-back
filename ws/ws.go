package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"msgv2-back/database"
	"msgv2-back/models"
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

				for i, item := range clients {
					if item == c {
						clients = append(clients[:i], clients[i+1:]...)
					}
				}

				break
			}

			log.Printf("recv: %s", msg)
			for _, client := range clients {
				if err = client.WriteMessage(mt, msg); err != nil {
					log.Println("write:", err)
					break
				}
			}

		}

	}))

	app.Get("/ws/check", websocket.New(func(c *websocket.Conn) {
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

			tag := models.Tag{}
			count := database.DB.Where("tag_id = ?", msg).First(&tag).RowsAffected

			if count == 0 {
				if err = c.WriteMessage(mt, []byte("#ffffff")); err != nil {
					log.Println("write:", err)
					break
				}

			} else {
				user := models.User{}
				database.DB.Where("id = ?", tag.UserID).First(&user)
				if err = c.WriteMessage(mt, []byte(user.Color)); err != nil {
					log.Println("write:", err)
					break
				}
			}

		}

	}))

}
