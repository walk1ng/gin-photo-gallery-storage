package models

type Auth struct {
	UserName string `json:"userName" gorm:"type:varchar(16)"`
	Password string `json:"password" gorm:"type:varchar(255)"`
	Email    string `json:"email" gorm:"type:varchar(128)"`
}
