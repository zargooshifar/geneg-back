package tags

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/tags"
	"msgv2-back/models"
)

func Routes(app *fiber.App) {
	app.Post("api/tags/check", tags.Check)

	app.Get("api/tags/tags", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), tags.TagList)
	app.Post("api/tags/tag", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), tags.TagEdit)
	app.Put("api/tags/tag", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), tags.TagCreate)
	app.Delete("api/tags/tag", utils.Secure(models.ROLES{models.ADMIN, models.OPERATOR, models.USER, models.GUEST}), tags.TagDelete)

}
