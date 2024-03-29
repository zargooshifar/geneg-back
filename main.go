package main

import (
	"encoding/binary"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"msgv2-back/database"
	"msgv2-back/handlers"
	"msgv2-back/radius"
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

	crypto_rand "crypto/rand"
	math_rand "math/rand"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {

	//go espserial.Config()

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

	go radius.Setup()
	log.Printf("Starting server on :3000")

	log.Fatal(app.Listen(*port))

	//app.Listen(":3000")
}

func seedRand() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
