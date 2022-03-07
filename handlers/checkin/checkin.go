package checkin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/models"
	"strconv"
)

func CheckInAll(c *fiber.Ctx) error {

	items := []models.CheckIn{}
	//ordering and paination parameters
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	order := c.Query("order")
	offset := (page - 1) * limit
	count := int64(0)

	database.DB.Model(&models.CheckIn{}).Offset(offset).Limit(limit).
		Order(order).
		Find(&items).Offset(-1).Count(&count)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count":   count,
		"results": items,
	})
}

type CheckInCreateItem struct {
	UserID string `json:"user_id"`
}

func CheckInCreate(c *fiber.Ctx) error {

	user_item := CheckInCreateItem{}
	c.BodyParser(&user_item)
	user := models.User{}

	if err := database.DB.Model(&models.User{}).Where("id = ?", user_item.UserID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.USER_NOT_EXIST,
		})
	}

	item := models.CheckIn{}
	item.Tagged = false
	item.User = user
	item.UserID = &user.ID

	if err := database.DB.Model(&models.CheckIn{}).Create(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
		})
	}
	return c.Status(200).JSON(item)
}

func CheckInDelete(c *fiber.Ctx) error {
	id := c.Query("id")
	item := models.CheckIn{}
	err := database.DB.Delete(&item, uuid.MustParse(id)).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	return c.Status(200).JSON(item)
}
