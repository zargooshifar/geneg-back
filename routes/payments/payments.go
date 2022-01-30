package payments

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/payments"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Get("api/payments/payments", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), payments.GetPayments)
}
