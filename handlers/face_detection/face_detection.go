package face_detection

import (
	"encoding/base64"
	nerrors "errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
	"log"
	"msgv2-back/database"
	"msgv2-back/errors"
	"msgv2-back/models"
	"os"
	"strings"
)

const (
	cascade    = "haarcascade_frontalface_default.xml"
	model_path = "faces.model"
)

func ReTrain()  {
	
}


func Train(userID uuid.UUID, imageURLs []models.FaceImage) error {

	faces := []gocv.Mat{}
	labels := []int{}

	faceID := models.FaceID{
		UserID: userID,
	}
	count := database.DB.Where("user_id = ?", imageURLs[0].UserID).Find(&faceID).RowsAffected
	if count == 0 {
		database.DB.Create(&faceID)
	}

	log.Println("faceID", faceID)

	for _, imageURL := range imageURLs {
		image, err := ReadBase64Image(imageURL.DataURL)
		if err != nil {
			//return nerrors.New(errors.FAILED_CONVERT_BASE64_TO_IMAGE)
			log.Println(err)
			continue
		}

		face, err := ExtractFace(image)
		if err != nil {
			log.Println(err)
			continue
		}
		faces = append(faces, face)
		labels = append(labels, faceID.ID)
		image.Close()
	}
	log.Println("images: ", len(imageURLs))
	log.Println("face founded: ", len(faces))

	recognizer := contrib.NewLBPHFaceRecognizer()

	_, err := os.Stat(model_path)
	if err == nil {
		log.Println("loading models file")
		recognizer.LoadFile(model_path)
	}

	log.Println("start training...")
	recognizer.Update(faces, labels)
	log.Println("finish training...")
	recognizer.SaveFile("faces.model")

	return nil
}


func ExtractFace(image gocv.Mat) (gocv.Mat, error) {

	classifier := gocv.NewCascadeClassifier()
	classifier.Load(cascade)
	faces := classifier.DetectMultiScale(image)
	if len(faces) != 1 {
		log.Println("face not founded")
		return gocv.Mat{}, nerrors.New(errors.NOT_SINGLE_FACE)
	}
	face := image.Region(faces[0])
	return face, nil
}

func ReadBase64Image(url string) (gocv.Mat, error) {

	encoded_data := strings.Split(url, ",")[1]
	decoded_data, err := base64.StdEncoding.DecodeString(encoded_data)
	if err != nil {
		return gocv.Mat{}, err
	}

	image, err := gocv.IMDecode(decoded_data, gocv.IMReadGrayScale)
	if err != nil {
		return gocv.Mat{}, err
	}

	return image, nil
}

type faces struct {
	Faces []string `json:"faces"`
}

func UploadFace(c *fiber.Ctx) error {
	faces := new(faces)
	user := c.Locals("user").(*models.User)
	c.BodyParser(&faces)

	faceImages := []models.FaceImage{}
	for _, face := range faces.Faces {
		item := models.FaceImage{
			UserID:  user.ID,
			DataURL: face,
		}
		err := database.DB.Model(&models.FaceImage{}).Create(&item).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": errors.DB_ERROR_SAVING,
			})
		}
		faceImages = append(faceImages, item)
	}

	err := Train(user.ID, faceImages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": "success",
	})

}

func Detect(c *fiber.Ctx) error {
	faces := new(faces)

	log.Println("load recognizer")
	recognizer := contrib.NewLBPHFaceRecognizer()
	_, err := os.Stat(model_path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errors.NOT_TRAINED,
		})
	}
	recognizer.LoadFile(model_path)
	log.Println("recognizer loaded!")

	c.BodyParser(&faces)

	users := _users{}

	for index, url := range faces.Faces {
		log.Println("detect image #", index)
		image, err := ReadBase64Image(url)
		if err != nil {
			log.Println("failed to convert image!")
			continue
		}
		face, err := ExtractFace(image)
		if err != nil {
			log.Println("failed to convert extract face!")
			continue
		}

		log.Println("predicting...")
		id := recognizer.Predict(face)
		log.Println("predicted: " , id)
		faceID := models.FaceID{}
		database.DB.Where("id = ?", id).Find(&faceID)
		user := models.User{}
		database.DB.Where("id = ?", faceID.UserID).Find(&user)
		if(!users.hasUser(user)){
			users = append(users, user)
		}

	}

	log.Println(users)

	return c.Status(fiber.StatusOK).JSON(users)

}

type _users []models.User

func (list _users) hasUser(a models.User) bool {
	for _, b := range list {
		if b.ID == a.ID {
			return true
		}
	}
	return false
}