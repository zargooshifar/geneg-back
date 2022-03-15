package espserial

import (
	"bufio"
	"github.com/goburrow/serial"
	"io/ioutil"
	"log"
	"msgv2-back/database"
	"msgv2-back/models"
	"strconv"
	"strings"
)

func FindESP() string {
	contents, _ := ioutil.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "ttyUSB") {
			return "/dev/" + f.Name()
		}
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return ""
}

func GetUserString(user *models.User, is_buffet bool, is_checkin bool) string {
	buffet := "0"
	checkin := "0"
	if is_checkin {
		checkin = "1"
	}
	if is_buffet {
		buffet = "1"
	}

	return buffet + ":" + checkin + ":" + user.ID.String() + ":" + user.FirstName + " " + user.LastName + ":" + strconv.Itoa(user.Balance) + ":" + user.Color + "\r\n"
}

func Config(serial_port serial.Port) {

	scanner := bufio.NewScanner(serial_port)
	for scanner.Scan() {
		data := scanner.Text()

		log.Println("serial: ", data)

		parts := strings.Split(data, ":")
		if len(parts) == 3 {
			is_buffet := parts[0] == "0"
			is_checkin := parts[1] == "1"
			tag_id := strings.TrimSpace(parts[2])
			if is_checkin {
				DoCheckin(tag_id, serial_port)
			} else if is_buffet {
				DoBuffet(tag_id)
			}
		}

	}
}

func DoCheckin(tag_id string, serial_port serial.Port) {

	tag := models.Tag{}
	count := database.DB.Where("tag_id = ?", tag_id).First(&tag).RowsAffected

	if count == 0 {

		if _, err := serial_port.Write([]byte("0:1:0:-:0:#ffffff\r\n")); err != nil {
			log.Println("write:", err)
			return
		}

	} else {
		user := models.User{}
		database.DB.Where("id = ?", tag.UserID).First(&user)
		checkin := models.CheckIn{}
		//checkin.User = user
		checkin.UserID = user.ID
		checkin.Tagged = true
		err := database.DB.Create(&checkin).Error

		if err == nil {
			serial_port.Write([]byte(GetUserString(&user, false, true)))
		} else {
			log.Println("write:", err)
			return
		}
	}
}

func DoBuffet(tag_id string) {

}
