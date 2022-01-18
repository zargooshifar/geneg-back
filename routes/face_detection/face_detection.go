package face_detection

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/face_detection"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Get("api/face/faces", handlers.GetItems(models.FaceImage{}))
	//app.Get("api/face_detection", face_detection.Train)
	app.Post("api/face/face", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER}), face_detection.UploadFace)
	app.Post("api/face/detect", face_detection.Detect)
	//app.Delete("api/face/face", utils.Secure(models.ROLES{models.ADMIN}), handlers.DeleteItem(models.FaceImage{}))
}
