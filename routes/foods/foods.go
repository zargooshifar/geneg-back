package foods

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/foods"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Get("api/foods/foods", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), handlers.GetItems(models.Food{}))
	app.Get("api/foods/buffet-items", foods.GetBuffetItems)
	app.Post("api/foods/add-card", foods.AddCart)
	app.Get("api/foods/food", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), handlers.GetItem(models.Food{}))
	app.Put("api/foods/food", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR}), handlers.CreateItem(models.Food{}))
	app.Post("api/foods/food", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR}), handlers.UpdateItem(models.Food{}))
	app.Delete("api/foods/food", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR}), handlers.DeleteItem(models.Food{}))

}
