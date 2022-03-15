package main

import (
	"encoding/binary"
	"flag"
	"github.com/goburrow/serial"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"msgv2-back/database"
	"msgv2-back/espserial"
	"msgv2-back/handlers"
	"msgv2-back/routes/auth"
	"msgv2-back/routes/checkin"
	"msgv2-back/routes/face_detection"
	"msgv2-back/routes/foods"
	"msgv2-back/routes/images"
	"msgv2-back/routes/payments"
	"msgv2-back/routes/reserve"
	"msgv2-back/routes/tags"
	"msgv2-back/routes/users"
	"msgv2-back/ws"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	crypto_rand "crypto/rand"
	math_rand "math/rand"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {

	com_port := espserial.FindESP()
	log.Println("serialport: ")
	log.Println(com_port)
	serial_config := &serial.Config{Address: "COM22", BaudRate: 115200, Timeout: time.Second * 1}

	serial_port, err := serial.Open(serial_config)

	if err != nil {
		log.Println("serial failed to connect!")

		log.Fatal(err)
	}
	defer serial_port.Close()
	go espserial.Config(serial_port)

	flag.Parse()

	database.ConnectDB()

	seedRand()

	app := fiber.New(fiber.Config{
		Prefork:   *prod, // go run app.go -prod
		BodyLimit: 50 * 1024 * 1024,
	})

	app.Use(cors.New())

	app.Use(recover.New())
	app.Use(logger.New())

	app.Static("/images", "./images/")

	auth.Routes(app)
	users.Routes(app)
	foods.Routes(app)
	images.Routes(app)
	payments.Routes(app)
	reserve.Routes(app)
	face_detection.Routes(app)
	tags.Routes(app)
	checkin.Routes(app)
	ws.Config(app)

	app.Use(handlers.NotFound)

	log.Fatal(app.Listen(*port))

	app.Listen(":3000")
}

func seedRand() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
