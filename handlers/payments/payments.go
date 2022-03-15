package payments

import (
	"github.com/gofiber/fiber/v2"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/handlers"
	"msgv2-back/models"
	"strconv"
)

func GetPayments(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	filter := handlers.FilterQuery(c)
	payments := []models.Payment{}
	if len(filter) > 0 {
		filter += " AND "
	}
	filter += "user_id = '" + user.ID.String() + "'"
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	order := c.Query("order")
	offset := (page - 1) * limit
	count := int64(0)
	database.DB.Model(&models.Payment{}).Offset(offset).Limit(limit).
		Order(order).
		Where(filter).
		Find(&payments).Offset(-1).Count(&count)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count":   count,
		"results": payments,
	})
}

func CreatePayments(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	payment := models.Payment{}
	c.BodyParser(&payment)

	payment.Description = "پرداختی - اضافه شده توسط " + user.FirstName + " " + user.LastName

	if err := database.DB.Model(&models.Payment{}).Create(&payment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
		})
	}
	return c.Status(fiber.StatusOK).JSON(payment)
}
