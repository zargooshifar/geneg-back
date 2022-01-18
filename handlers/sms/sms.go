package sms

import (
	"fmt"
	"log"
	"math/rand"
	"msgv2-back/database"
	"msgv2-back/models"
	"net/http"
	"strings"
	"time"
)

var api_key = "3460f94ac23b053368992f42513a8d9709a69fde7821474da454363f991c979b"

func SendPin(number string) (string, error) {
	pin := 10000 + rand.Intn(89999)

	//delete previous verifications
	database.DB.Where(&models.VerificationSMS{Number: number}).Delete(&models.VerificationSMS{})

	verification := new(models.VerificationSMS)

	verification.Pin = pin
	verification.Number = number
	verification.Expire = time.Now().Add(2 * time.Minute).Unix()

	if err := database.DB.Create(&verification).Error; err != nil {
		return "", err
	}

	SendOTP(number, pin)

	log.Println("sending pin: %d", pin)
	return verification.ID.String(), nil
}

func SendOTP(number string, pin int) error {

	url := "https://api.ghasedak.me/v2/verification/send/simple"

	payload := strings.NewReader(fmt.Sprintf("receptor=%s&template=humanotp&type=1&param1=%d", number, pin))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("apikey", api_key)
	req.Header.Add("cache-control", "no-cache")

	_, err := http.DefaultClient.Do(req)

	//Handle Error
	if err != nil {
		return err
	}

	//body, _ := ioutil.ReadAll(response.Body)
	//
	//fmt.Println(response)
	//fmt.Println(string(body))

	return nil
}

func SendSMS(number string, message string) error {

	url := "https://api.ghasedak.me/v2/sms/send/simple"
	line := "30005006005013"
	payload := strings.NewReader(fmt.Sprintf("message=%s&receptor=%s&linenumber=%s", message, number, line))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("apikey", api_key)
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	return nil
}
