package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
)

// Auth struct model represent auth table
type Auth struct {
	BaseModel
	UserName string `json:"user_name" gorm:"type:varchar(16)"`
	Password string `json:"password" gorm:"type:varchar(255)"`
	Email    string `json:"email" gorm:"type:varchar(128)"`
}

var ErrAuthExist = errors.New("auth already exists")

// AddAuth func to add a new auth
func AddAuth(username, password, email string) error {
	// transaction
	trx := db.Begin()
	defer trx.Commit()

	auth := Auth{}

	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_name = ?", username).
		First(&auth)

	if auth.ID > 0 {
		return ErrAuthExist
	}

	hash := md5.New()
	io.WriteString(hash, password)
	auth.UserName = username
	auth.Password = fmt.Sprintf("%x", hash.Sum(nil))
	auth.Email = email
	err := trx.Create(&auth).Error
	if err != nil {
		return err
	}
	return nil
}

// CheckAuth func check if the auth is valid
func CheckAuth(username, password string) bool {
	trx := db.Begin()
	defer trx.Commit()

	auth := Auth{}

	hash := md5.New()
	io.WriteString(hash, password)
	password = fmt.Sprintf("%x", hash.Sum(nil))

	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_name = ? AND password = ?", username, password).
		First(&auth)

	if auth.ID > 0 {
		return true
	}
	return false
}
