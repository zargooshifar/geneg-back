package auth

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/handlers/auth/utils"
	"msgv2-back/handlers/sms"
	"msgv2-back/models"
	"net"
	"strings"
	"time"
)

func Login(c *fiber.Ctx) error {
	l := new(models.Login)

	if err := c.BodyParser(l); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}

	user := new(models.User)

	if count := database.DB.Where(&models.User{Username: l.Username}).First(&user).RowsAffected; count == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.USER_NOT_EXIST,
		})
	}

	//check password
	if !utils.VerifyPassword(user.Password, l.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.ERROR_WRONG_PASSWORD,
		})
	}

	accessToken, refreshToken := utils.GenerateTokens(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access":  accessToken,
		"refresh": refreshToken,
	})

}

func CheckUserName(c *fiber.Ctx) error {

	username := new(models.LoginUserName)

	if err := c.BodyParser(username); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}
	exists := (database.DB.Where(&models.User{Username: username.Username}).First(&models.User{}).RowsAffected > 0)
	if !exists {

		remoteIP := getIPAdress(c)
		log.Println(remoteIP)
		if remoteIP != "81.16.121.206" && remoteIP != "192.168.31.1" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": errors.CANT_REGISTER_OUTSIDE_CORP,
			})
		}

		verification_id, err := sms.SendPin(username.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"verification_id": verification_id,
			"exists":          exists,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exists": exists,
	})
}

func ForgetPin(c *fiber.Ctx) error {

	username := new(models.LoginUserName)

	if err := c.BodyParser(username); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}

	exists := (database.DB.Where(&models.User{Username: username.Username}).First(&models.User{}).RowsAffected > 0)

	verification_id, err := sms.SendPin(username.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"verification_id": verification_id,
		"exists":          exists,
	})

}

func ResetPassword(c *fiber.Ctx) error {

	reg := new(models.Registration)
	if err := c.BodyParser(reg); err != nil {
		return c.Status(fiber.StatusBadRequest).Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}

	verification := new(models.VerificationSMS)
	if count := database.DB.Where(&models.VerificationSMS{ID: uuid.MustParse(reg.Verification)}).First(&verification).RowsAffected; count == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.VERIFICATION_NOT_EXIST,
		})
	}

	if !verification.Confirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.VERIFICATION_NOT_CONFIRMED,
		})
	}

	user := models.User{}

	if count := database.DB.Where(&models.User{Username: verification.Number}).First(&user).RowsAffected; count == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.USER_NOT_EXIST,
		})
	}

	// Hashing the password with a random salt
	password := []byte(reg.Password)
	hashedPassword, err := utils.GenerateHashPassword(password)

	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)

	if err := database.DB.Table("users").Where("id = ?", user.ID).Updates(map[string]interface{}{"password": string(hashedPassword)}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
		})
	}

	//delete verification...
	database.DB.Delete(verification)

	// setting up the authorization cookies
	accessToken, refreshToken := utils.GenerateTokens(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access":  accessToken,
		"refresh": refreshToken,
	})

}

func VerifyPin(c *fiber.Ctx) error {
	max_attempts := 5
	pin := new(models.Pin)
	if err := c.BodyParser(pin); err != nil {
		return c.Status(fiber.StatusBadRequest).Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}

	verification := new(models.VerificationSMS)

	exists := (database.DB.Where(&models.VerificationSMS{ID: pin.ID}).First(&verification).RowsAffected > 0)
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_VERIFICATION_ID,
		})
	}

	if verification.Expire < time.Now().Unix() || verification.Attempts > max_attempts {
		database.DB.Delete(verification)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.EXPIRE_PIN,
		})
	}

	if verification.Pin == pin.Pin {
		verification.Confirm = true
		database.DB.Save(verification)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":         true,
			"message":         "",
			"verification_id": verification.ID.String(),
		})
	} else {
		verification.Attempts += 1
		database.DB.Save(verification)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_PIN,
		})
	}

}

func CompleteRegister(c *fiber.Ctx) error {

	reg := new(models.Registration)
	if err := c.BodyParser(reg); err != nil {
		return c.Status(fiber.StatusBadRequest).Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}

	// validate if the email, username and password are in correct format
	//regErrors := utils.ValidateRegister(reg)
	//if regErrors.Err {
	//	return c.Status(fiber.StatusBadRequest).JSON(regErrors)
	//}

	verification := new(models.VerificationSMS)
	if count := database.DB.Where(&models.VerificationSMS{ID: uuid.MustParse(reg.Verification)}).First(&verification).RowsAffected; count == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.VERIFICATION_NOT_EXIST,
		})
	}

	if !verification.Confirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.VERIFICATION_NOT_CONFIRMED,
		})
	}

	if count := database.DB.Where(&models.User{Username: verification.Number}).First(new(models.User)).RowsAffected; count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.USER_EXIST,
		})
	}

	user := new(models.User)

	// Hashing the password with a random salt
	password := []byte(reg.Password)
	hashedPassword, err := utils.GenerateHashPassword(password)

	if err != nil {
		panic(err)
	}
	user.Username = verification.Number
	user.FirstName = reg.FirstName
	user.LastName = reg.LastName
	user.Password = string(hashedPassword)
	user.Role = "user"

	////first user is admin!
	//users_count := database.DB.Find(&[]models.User{}).RowsAffected
	//if(users_count == 0){
	//	user.Role = "admin"
	//}
	//fmt.Println("users: ", users_count)

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.DB_ERROR_SAVING,
		})
	}

	//delete verification...
	database.DB.Delete(verification)

	// setting up the authorization cookies
	accessToken, refreshToken := utils.GenerateTokens(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access":  accessToken,
		"refresh": refreshToken,
	})
}

func Refresh(c *fiber.Ctx) error {
	r := new(models.RefreshToken)
	if err := c.BodyParser(r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.WRONG_INPUT,
		})
	}
	access, refresh, error := utils.RefreshTokens(r)
	if len(error) > 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access":  access,
		"refresh": refresh,
	})
}

type ipRange struct {
	start net.IP
	end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var privateRanges = []ipRange{
	ipRange{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	ipRange{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	ipRange{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	ipRange{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	//ipRange{
	//	start: net.ParseIP("192.168.0.0"),
	//	end:   net.ParseIP("192.168.255.255"),
	//},
	ipRange{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func getIPAdress(c *fiber.Ctx) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(c.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	return ""
}

func Logout(c *fiber.Ctx) error {

	log.Println(getIPAdress(c))
	//TODO: should remove current refresh key user is using
	//TODO: additionally we can add a functionality to delete all user refresh tokens, to remove all active logins
	return c.SendStatus(fiber.StatusOK)
}
