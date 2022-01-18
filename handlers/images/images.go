package images

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/models"
)

func Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("document")
	if(err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
			"err":     err,
		})
	}
	file_params := models.Image{}
	c.BodyParser(&file_params)
	log.Println(file_params.Name)
	image := models.Image{
		Name: file_params.Name,
	}
	err = database.DB.Create(&image).Error
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
			"err":     err,
		})
	}

	path := fmt.Sprintf("./images/foods/%s", image.ID)


	image.Path = path;


	log.Println(file_params)

	if err == nil {
		log.Println("saving file!")
		c.SaveFile(file, path)
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = database.DB.Updates(&image).Error
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
			"err":     err,
		})
	}

	return  c.Status(fiber.StatusOK).JSON(image)
}
