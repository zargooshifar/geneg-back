package tags

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/models"
	"strconv"
)

type tagCheck struct {
	TagID string `json:"tag_id"`
}
type tagUser struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}

func Check(c *fiber.Ctx) error {
	tagCheck := tagCheck{}
	c.BodyParser(&tagCheck)

	tag := models.Tag{}
	count := database.DB.Model(&models.Tag{}).Where("tag_id = ?", tagCheck.TagID).First(&tag).RowsAffected
	if count == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"assigned": false,
		})
	}

	user := models.User{}
	database.DB.Model(&models.User{}).Where("id = ?", tag.UserID).First(&user)

	result := tagUser{
		Name: user.FirstName + " " + user.LastName,
		Id:   user.ID,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"assigned": true,
		"user":     result,
	})
}

func TagList(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	items := []models.Tag{}

	//ordering and paination parameters
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	order := c.Query("order")
	offset := (page - 1) * limit
	count := int64(0)

	database.DB.Model(&models.Tag{}).Offset(offset).Limit(limit).
		Order(order).
		Where("user_id = ?", user.ID).
		Find(items).Offset(-1).Count(&count)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count":   count,
		"results": items,
	})
}

func TagCreate(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	tag := models.Tag{}
	c.BodyParser(&tag)
	tag.UserID = user.ID
	if err := database.DB.Model(&models.Tag{}).Create(tag).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
		})
	}
	return c.Status(200).JSON(tag)
}

func TagEdit(c *fiber.Ctx) error {
	temp := models.Tag{}
	c.BodyParser(&temp)
	base := new(models.Base)
	c.BodyParser(base)
	if err := database.DB.Model(&models.Tag{}).Where("id = ?", base.ID.String()).Select("*").Updates(temp).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
			"err":     err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(temp)
}

func TagDelete(c *fiber.Ctx) error {
	id := c.Query("id")
	tag := models.Tag{}
	err := database.DB.Delete(tag, uuid.MustParse(id)).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	return c.Status(200).JSON(tag)
}
