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

	database.DB.Model(&models.Food{}).Where("expire >= ? AND type <> ?", time.Now(), models.BUFFET).Find(&foods)

	result := []reserveFood{}

	for _, food := range foods {
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

func GetBuffetItems(c *fiber.Ctx) error {
	foods := []models.Food{}
	database.DB.Model(&models.Food{}).Where("expire >= ? AND type = ?", time.Now(), models.BUFFET).Find(&foods)
	return c.Status(fiber.StatusOK).JSON(foods)
}

type reserve_result struct {
	UserName string `json:"user_name"`
	Count    int    `json:"count"`
	Food     string `json:"food"`
}

func GetTodayReserves(c *fiber.Ctx) error {
	foods := []models.Food{}
	start := time.Now().Format("2006-01-02 00:00:00")
	end := time.Now().Add(time.Hour * 24).Format("2006-01-02 00:00:00")
	database.DB.Model(&models.Food{}).Where("expire >= ? AND expire < ? AND type = ?", start, end, models.LAUNCH).Find(&foods)

	result := []reserve_result{}

	for _, food := range foods {

		reserves := []models.Reserve{}
		database.DB.Model(&models.Reserve{}).Where("food_id = ?", food.ID).Find(&reserves)

		for _, r := range reserves {
			user := models.User{}
			database.DB.Where("id = ?", r.UserID).Find(&user)
			new_result := reserve_result{
				Food:     food.Name,
				UserName: user.FirstName + " " + user.LastName,
				Count:    r.Count,
			}

			result = append(result, new_result)
		}

	}

	return c.Status(fiber.StatusOK).JSON(result)
}

type card_item struct {
	Count int `json:"count"`
}

type card struct {
	UserID string               `json:"user_id"`
	Card   map[string]card_item `json:"card"`
}

func AddCart(c *fiber.Ctx) error {
	card := card{}
	c.BodyParser(&card)

	for key := range card.Card {

		new_reserve := models.Reserve{
			UserID: uuid.MustParse(card.UserID),
			FoodID: uuid.MustParse(key),
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

		new_reserve.Count += card.Card[key].Count

		if err := database.DB.Table("reserves").Select("*").Updates(&new_reserve).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": errors.DB_ERROR_SAVING,
				"err":     err,
			})
		}

	}
	return c.Status(fiber.StatusOK).JSON(card)
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
