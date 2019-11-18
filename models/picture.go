package models

import (
	"fmt"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/otaviokr/pics-at-home/utils"
)

// Picture is the core entity of this program. Any image is defined as a 'picture'.
type Picture struct {
	gorm.Model
	Path        string
	Description string
}

// Validate checks if the Picture entity has valid data, before saving. If field ID is defined (higher than 0),
// it will check if it exists in the database.
func (pic *Picture) Validate(db *gorm.DB) (map[string]interface{}, bool) {
	picDB := &Picture{}

	if strings.TrimSpace(pic.Path) == "" {
		return utils.Message(false, "Pic entity does not have a valid Path attribute!"), false

	} else if pic.ID > 0 {
		err := db.Table("pictures").First(picDB, pic.ID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.Message(false, fmt.Sprintf("Picture not found in Database: %d", pic.ID)), false
			}

			return utils.Message(false, "Connection error. Failed to validated picture ID."), false
		}

	} else {
		// Avoid to save a picture that's already registered in database, i.e., we expect to get ErrRecordNotFound error.
		//err := db.Table("pictures").Where("path = ?", pic.Path).First(picDB).Error
		err := db.Table("pictures").Where("path = ?", pic.Path).First(picDB).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return utils.Message(false, fmt.Sprintf("Unexpected error while validating picture path in database: %s - Error: %v", pic.Path, err)), false
			}
		} else if picDB.ID > 0 {
			return utils.Message(false, fmt.Sprintf("Picture path already stored in database: %s", pic.Path)), false
		}
	}

	return utils.Message(false, "Validation passed"), true
}

// Create saves the new picture in database.
func (pic *Picture) Create(db *gorm.DB) map[string]interface{} {
	if response, ok := pic.Validate(db); !ok {
		return response
	}

	if pic.CreatedAt.IsZero() {
		pic.CreatedAt = time.Now()
	}

	if err := db.Create(pic).Error; err != nil {
		return utils.Message(false, fmt.Sprintf("Error inserting into database: %v", err))
	}

	if pic.ID < 1 {
		return utils.Message(false, "Failed to insert picture data into database.")
	}

	response := utils.Message(true, "New picture inserted")
	response["picture"] = pic
	return response
}

// BUG(otaviokr) : The first query should get only the IDs; we do not need the entire rows.

// GetRandomPicture will return a random picture from database.
func GetRandomPicture(db *gorm.DB) image.Image {
	pic := &Picture{}
	//var picIDs []uint
	//err := GetDB().Table("pictures").Select("id").Order("id ASC").Find(&picIDs).Error
	var picIDs []Picture
	err := db.Table("pictures").Order("id ASC").Find(&picIDs).Error
	if err != nil {
		fmt.Printf("Failed to get pictures IDs in database. %v\n", err)
		return nil
	}

	if len(picIDs) > 0 {
		picDB := picIDs[rand.Intn(len(picIDs))]
		err = db.Table("pictures").First(pic, picDB.ID).Error
		if err != nil {
			fmt.Printf("Failed to get picture by specific ID: %d\n", picDB.ID)
			return nil
		}

		picFile, err := os.Open(pic.Path)
		if err != nil {
			fmt.Printf("Could not open picture: %s\n", pic.Path)
			return nil
		}

		img, err := jpeg.Decode(picFile)
		if err != nil {
			fmt.Printf("Could not decode picture: %s\n", pic.Path)
			return nil
		}
		picFile.Close()

		return img
	}

	return nil
}

// GetRandomPictureInfo will return a random picture from database.
func GetRandomPictureInfo(db *gorm.DB) Picture {
	pic := Picture{}
	//var picIDs []uint
	//err := GetDB().Table("pictures").Select("id").Order("id ASC").Find(&picIDs).Error
	var picIDs []Picture
	err := db.Table("pictures").Order("id ASC").Find(&picIDs).Error
	if err != nil {
		fmt.Printf("Failed to get pictures IDs in database. %v\n", err)
		return Picture{}
	}

	if len(picIDs) > 0 {
		picDB := picIDs[rand.Intn(len(picIDs))]
		err = db.Table("pictures").First(&pic, picDB.ID).Error
		if err != nil {
			fmt.Printf("Failed to get picture by specific ID: %d\n", picDB.ID)
			return Picture{}
		}
	}
	return pic
}

// GetRecentPics will return a list of n recently added pictures.
func GetRecentPics(n uint, db *gorm.DB) []Picture {
	var pics []Picture

	err := db.Table("pictures").Order("created_at DESC").Limit(n).Find(&pics).Error
	if err != nil {
		fmt.Printf("Failed to get pictures IDs in database. %v\n", err)
		return nil
	}

	return pics
}

// GetPictureByID will return the entry determined by ID given.
func GetPictureByID(id uint, db *gorm.DB) *Picture {
	pic := &Picture{}

	err := db.Table("pictures").Where("id = ?", id).First(pic).Error
	if err != nil {
		fmt.Println("Failed to get picture from database.")
		return nil
	}

	return pic
}
