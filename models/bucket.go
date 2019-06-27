package models

import (
	"errors"
	"log"

	"github.com/walk1ng/gin-photo-gallery-storage/constant"
)

// Bucket struct model represent bucket table
type Bucket struct {
	BaseModel
	AuthID      uint   `json:"auth_id" gorm:"type:int" form:"auth_id"`
	Name        string `json:"name" gorm:"type:varchar(64)" form:"bucket_name"`
	State       int    `json:"state" gorm:"type:tinyint(1)" form:"state"`
	Size        int    `json:"size" gorm:"type:int" form:"bucket_size"`
	Description string `json:"description" gorm:"type:text" form:"description"`
}

var ErrBucketExists = errors.New("bucket already exists")
var ErrNoSuchBucket = errors.New("no such bucket")

// AddBucket func add a new bucket
func AddBucket(bucketToAdd *Bucket) error {
	trx := db.Begin()
	defer trx.Commit()

	// check if the bucket exists
	bucket := Bucket{}
	trx.Set("gorm:query_option", "FOR UPDAE").
		Where("auth_id = ? AND name = ? AND state = ?", bucketToAdd.AuthID, bucketToAdd.Name, 1).
		First(&bucket)

	if bucket.ID > 0 {
		return ErrBucketExists
	}

	// insert the bucket
	bucket.AuthID = bucketToAdd.AuthID
	bucket.Name = bucketToAdd.Name
	bucket.State = 1
	bucket.Size = 0
	bucket.Description = bucketToAdd.Description

	if err := trx.Create(&bucket).Error; err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DeleteBucket func delete an existed bucket
func DeleteBucket(bucketID uint) error {
	trx := db.Begin()
	defer trx.Commit()

	result := trx.Where("id = ? AND state = ?", bucketID, 1).Delete(Bucket{})
	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return ErrNoSuchBucket
	}

	return nil
}

// UpdateBucket func update an existed bucket
func UpdateBucket(bucketToUpdate *Bucket) error {
	trx := db.Begin()
	defer trx.Commit()

	bucket := Bucket{}
	bucket.ID = bucketToUpdate.ID
	result := trx.Model(&bucket).Update(*bucketToUpdate)
	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return ErrNoSuchBucket
	}

	return nil
}

// GetBucketByID func get bucket by bucket id
func GetBucketByID(bucketID uint) (Bucket, error) {
	trx := db.Begin()
	defer trx.Commit()

	bucket := Bucket{}
	trx.Where("id = ?", bucketID).First(&bucket)

	if bucket.ID > 0 {
		return bucket, nil
	}

	return bucket, ErrNoSuchBucket

}

// GetBucketByAuthID func get all buckets by the given user
func GetBucketByAuthID(authID uint, offset int) ([]Bucket, error) {
	trx := db.Begin()
	defer trx.Commit()

	buckets := make([]Bucket, 0, constant.PageSize)
	err := trx.Where("auth_id = ?", authID).
		Offset(offset).
		Limit(constant.PageSize).
		Find(&buckets).
		Error

	if err != nil {
		return buckets, err
	}
	return buckets, nil
}
