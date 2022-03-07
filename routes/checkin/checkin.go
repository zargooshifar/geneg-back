package checkin

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/checkin"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Get("api/checkin/checkins", utils.Secure(models.ROLES{models.ADMIN}), checkin.CheckInAll)
	app.Put("api/checkin/checkin", utils.Secure(models.ROLES{models.ADMIN}), checkin.CheckInCreate)
	app.Delete("api/checkin/checkin", utils.Secure(models.ROLES{models.ADMIN}), checkin.CheckInDelete)

}
