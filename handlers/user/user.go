package user

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/models"
)

func UpdateUser(c *fiber.Ctx) error {
	user_params := new(models.UserUpdate)

	user := c.Locals("user").(*models.User)
	//model := reflect.New(reflect.TypeOf(item)).Interface()
	c.BodyParser(&user_params)
	base := new(models.Base)
	c.BodyParser(base)

	if err := database.DB.Table("users").Where("id = ?", user.ID).Updates(map[string]interface{}{"first_name": user_params.FirstName, "last_name": user_params.LastName}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
			"err":     err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)

}

func GetProfile(c *fiber.Ctx) error {

	user := c.Locals("user").(*models.User)

	return c.Status(fiber.StatusOK).JSON(user)

}
