package models

import (
	"errors"
	"log"
	"mime/multipart"
	"os"

	"github.com/walk1ng/gin-photo-gallery-storage/utils"

	"github.com/walk1ng/gin-photo-gallery-storage/constant"

	"github.com/jinzhu/gorm"
)

// Photo struct model repesent the photo table
type Photo struct {
	BaseModel
	AuthID      uint     `json:"auth_id" gorm:"type:int" form:"auth_id"`
	BucketID    uint     `json:"bucket_id" gorm:"type:int" form:"bucket_id"`
	Name        string   `json:"name" gorm:"type:varchar(255)" form:"name"`
	Tag         string   `json:"tag" gorm:"type:varchar(255)" form:"tag"`
	Tags        []string `json:"tags" gorm:"-" form:"tags"`
	URL         string   `json:"url" gorm:"type:varchar(255)" form:"url"`
	Description string   `json:"description" gorm:"type:text" form:"description"`
	State       int      `json:"state" gorm:"type:tinyint(1)" form:"state"`
}

var ErrPhotoExists = errors.New("photo already exists")
var ErrNoSuchPhoto = errors.New("no such photo")
var ErrPhotoFileBroken = errors.New("photo file is broken")

// AddPhoto func add a new photo
func AddPhoto(photoToAdd *Photo, photoFileHeader *multipart.FileHeader) (*Photo, string, error) {
	trx := db.Begin()
	defer trx.Commit()

	// check if the photo exist
	photo := Photo{}
	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("bucket_id = ? AND name = ?", photoToAdd.BucketID, photoToAdd.Name).
		First(&photo)

	if photo.ID > 0 {
		return nil, "", ErrPhotoExists
	}

	photo.AuthID = photoToAdd.AuthID
	photo.BucketID = photoToAdd.BucketID
	photo.Name = photoToAdd.Name
	photo.Tag = photoToAdd.Tag
	photo.Description = photoToAdd.Description
	photo.State = 1

	// insert the new photo to photo table
	err := trx.Create(&photo).Error
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	// update the related bucket
	err = trx.Model(&Bucket{}).
		Where("id = ?", photoToAdd.BucketID).
		Update("size", gorm.Expr("size + ?", 1)).
		Error

	if err != nil {
		trx.Rollback()
		log.Println(err)
		return nil, "", err
	}

	// todo: upload the photo to the cloud
	if photoFile, err := photoFileHeader.Open(); err == nil {
		uploadID := utils.Upload(photo.ID, photo.Name, photoFile.(*os.File))
		return &photo, uploadID, nil
	} else {
		log.Println(err)
		return nil, "", ErrPhotoFileBroken
	}
}

// DeletePhotoByID func delete a photo by ID
func DeletePhotoByID(photoID uint) error {
	trx := db.Begin()
	defer trx.Commit()

	result := trx.Where("id = ?", photoID).Delete(Photo{})
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return ErrNoSuchPhoto
	}

	return nil
}

// DeletePhotoByBucketIDAndPhotoName func delete a photo by its bucket id and its name
func DeletePhotoByBucketIDAndPhotoName(bucketID uint, name string) error {
	trx := db.Begin()
	defer trx.Commit()

	result := trx.Where("bucket_id = ? AND name = ?", bucketID, name).Delete(Photo{})
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return ErrNoSuchPhoto
	}
	return nil
}

// UpdatePhoto func update a photo
func UpdatePhoto(photoToUpdate *Photo) (*Photo, error) {
	trx := db.Begin()
	defer trx.Commit()

	photo := Photo{}
	photo.ID = photoToUpdate.ID

	result := trx.Model(&photo).Updates(*photoToUpdate)
	if err := result.Error; err != nil {
		log.Println(err)
		return nil, err
	}

	if result.RowsAffected == 0 {
		return &photo, ErrNoSuchPhoto
	}

	return &photo, nil
}

// UpdatePhotoURL func update the url of a photo
func UpdatePhotoURL(photoID uint, url string) error {
	trx := db.Begin()
	defer trx.Commit()

	photo := Photo{}
	photo.ID = photoID

	err := trx.Model(&photo).Update("url", url).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetPhotoByID func get the photo by its photo ID
func GetPhotoByID(photoID uint) (*Photo, error) {
	trx := db.Begin()
	defer trx.Commit()

	photo := Photo{}
	err := trx.Where("id = ?", photoID).First(&photo).Error
	if err != nil || photoID == 0 {
		log.Println(err)
		return &photo, err
	}

	return &photo, nil
}

// GetPhotosByBucketID func get photos by its bucket ID
func GetPhotosByBucketID(bucketID uint, offset int) ([]Photo, error) {
	trx := db.Begin()
	defer trx.Commit()

	photos := make([]Photo, 0, constant.PageSize)
	err := trx.Where("bucket_id = ?", bucketID).
		Offset(offset).
		Limit(constant.PageSize).
		Find(&photos).
		Error

	if err != nil {
		return photos, err
	}

	return photos, nil
}

// GetPhotoUploadStatus func check photo upload status
func GetPhotoUploadStatus(uploadID string) int {
	status := utils.GetUploadStatus(uploadID)
	switch status {
	case -2:
		return constant.PhotoNotExist
	case -1:
		return constant.PhotoUploadError
	case 0:
		return constant.PhotoUploadSuccess
	case 1:
		return constant.PhotoAddInProcess
	default:
		return constant.InvalidParams
	}
}
