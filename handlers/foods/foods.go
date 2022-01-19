package foods

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/models"
	"time"
)

const (
	ADD    = "add"
	REMOVE = "remove"
)

type reserve struct {
	//UserID string
	//ReserveID string `json:"reserve_id"`
	FoodID string `json:"food_id"`
	Type   string `json:"type"`
}

type reserveFood struct {
	Food  models.Food `json:"food"`
	Count int         `json:"count"`
}

func GetReserves(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	foods := []models.Food{}

	database.DB.Model(&models.Food{}).Where("expire >= ?", time.Now()).Find(&foods)

	result := []reserveFood{}

	for _, food := range foods {
		if food.Type == models.BUFFET {
			continue
		}
		reserve := models.Reserve{}
		database.DB.Model(&models.Reserve{}).Where(&models.Reserve{
			UserID: user.ID,
			FoodID: food.ID,
		}).Find(&reserve)
		result = append(result, reserveFood{
			Food:  food,
			Count: reserve.Count,
		})
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func ReserveFood(c *fiber.Ctx) error {
	reserve_params := new(reserve)

	user := c.Locals("user").(*models.User)
	c.BodyParser(&reserve_params)

	new_reserve := models.Reserve{
		UserID: user.ID,
		FoodID: uuid.MustParse(reserve_params.FoodID),
	}
	food := models.Food{}
	if food.Expire.After(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.EXPIRE_FOOD,
			"err":     "err",
		})
	}

	reserve_count := database.DB.Table("reserves").Where(&new_reserve).Find(&new_reserve).RowsAffected

	if reserve_count == 0 {

		if err := database.DB.Create(&new_reserve).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": errors.DB_ERROR_SAVING,
				"err":     err,
			})
		}
	}

	if reserve_params.Type == ADD {
		new_reserve.Count += 1
	} else if reserve_params.Type == REMOVE {
		if new_reserve.Count == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": errors.NO_RESERVE,
				"err":     "err",
			})
		}
		new_reserve.Count -= 1
	}

	if err := database.DB.Table("reserves").Select("*").Updates(&new_reserve).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
			"err":     err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(new_reserve)

}
