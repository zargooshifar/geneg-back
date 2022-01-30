package images

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/images"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Get("api/images/images", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), handlers.GetItems(models.Image{}))
	app.Get("api/images/image", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), handlers.GetItem(models.Image{}))
	app.Put("api/images/upload", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR}), images.Upload)

	app.Post("api/images/image", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR}), handlers.UpdateItem(models.Image{}))
	app.Delete("api/images/image", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR}), handlers.DeleteItem(models.Image{}))

}
