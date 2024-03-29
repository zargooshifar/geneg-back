package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/models"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).Send([]byte("the page you are looking for is not exist."))
}

func FilterQuery(c *fiber.Ctx) string {
	rawParams := string(c.Request().URI().QueryString())
	params := strings.Split(rawParams, "&")
	query := ""
	for index, part := range params {
		for k, r := range part {
			if r == '=' {
				if (len(part[k+1:]) == 0) || (part[:k] == "limit") || (part[:k] == "page") || (part[:k] == "order") {
					continue
				}

				value, _ := url.QueryUnescape(part[k+1:])
				if part[:k] == "expire" {
					query += part[:k] + " >= '" + value + "'"
				} else {
					query += part[:k] + " LIKE '%" + value + "%'"
				}
				if index < len(params)-1 {
					query += " AND "
				}
			}
		}
	}
	return strings.TrimSuffix(query, " AND ")
}

func GetItems(item interface{}) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)

		filter := FilterQuery(c)

		accessed_item := make([]reflect.StructField, 0)

		fields := reflect.VisibleFields(reflect.TypeOf(item))
		for _, field := range fields {
			_roles, roled := field.Tag.Lookup("access")
			if roled {
				roles := models.ROLES{}
				roles = strings.Split(_roles, ",")
				if roles.Has(user.Role) {
					accessed_item = append(accessed_item, field)
				}
			} else {
				accessed_item = append(accessed_item, field)
			}
		}

		items := reflect.New(reflect.SliceOf(reflect.StructOf(accessed_item))).Interface()

		//ordering and paination parameters
		limit, _ := strconv.Atoi(c.Query("limit"))
		page, _ := strconv.Atoi(c.Query("page"))
		order := c.Query("order")
		offset := (page - 1) * limit
		count := int64(0)

		database.DB.Model(item).Offset(offset).Limit(limit).
			Order(order).
			Where(filter).
			Find(items).Offset(-1).Count(&count)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count":   count,
			"results": items,
		})
	}
}

func GetItem(item interface{}) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		temp := item
		//model := reflect.New(reflect.TypeOf(item)).Interface()
		id := c.Query("id")
		if err := database.DB.Find(temp, "id = ?", id).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": errors.DB_ERROR_SAVING,
			})
		}
		return c.Status(200).JSON(temp)
	}
}

func CreateItem(item interface{}) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		temp := reflect.New(reflect.TypeOf(item)).Interface()
		c.BodyParser(&temp)
		log.Println(temp)
		if err := database.DB.Model(item).Create(temp).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": errors.DB_ERROR_SAVING,
			})
		}
		return c.Status(200).JSON(temp)
	}
}

func UpdateItem(item interface{}) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		temp := item
		//model := reflect.New(reflect.TypeOf(item)).Interface()
		c.BodyParser(&temp)
		base := new(models.Base)
		c.BodyParser(base)
		if err := database.DB.Model(item).Where("id = ?", base.ID.String()).Select("*").Updates(temp).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": errors.DB_ERROR_SAVING,
				"err":     err,
			})
		}
		return c.Status(fiber.StatusOK).JSON(temp)
	}
}

func DeleteItem(item interface{}) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Query("id")
		temp := reflect.New(reflect.TypeOf(item)).Interface() // for some reason item instance itself not work with gorm, so i create a instance with same type!
		err := database.DB.Delete(temp, uuid.MustParse(id)).Error
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err,
			})
		}
		return c.Status(200).JSON(temp)
	}
}
