package reserve

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/foods"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Put("api/reserves/reserve", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), foods.ReserveFood)
	app.Get("api/reserves/reserves", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), foods.GetReserves)
	app.Get("api/reserves/today", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR}), foods.GetTodayReserves)
}
