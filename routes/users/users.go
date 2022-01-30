package users

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/user"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Get("api/admin/users", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER}), handlers.GetItems(models.User{}))
	app.Get("api/admin/user", utils.Secure(models.ROLES{models.ADMIN}), handlers.GetItem(models.User{}))
	app.Put("api/admin/user", utils.Secure(models.ROLES{models.ADMIN}), handlers.CreateItem(models.User{}))
	app.Post("api/admin/user", utils.Secure(models.ROLES{models.ADMIN}), handlers.UpdateItem(models.User{}))
	app.Delete("api/admin/user", utils.Secure(models.ROLES{models.ADMIN}), handlers.DeleteItem(models.User{}))
	app.Post("api/user/update", utils.Secure(models.ROLES{models.USER, models.ADMIN, models.GUEST, models.OPERATOR}), user.UpdateUser)
	app.Get("api/user/get_profile", utils.Secure(models.ROLES{models.USER, models.ADMIN, models.GUEST, models.OPERATOR}), user.GetProfile)

}
