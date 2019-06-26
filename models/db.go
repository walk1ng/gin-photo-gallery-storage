package models

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/walk1ng/gin-photo-gallery-storage/conf"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"
)

var db *gorm.DB

// BaseModel struct
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT" form:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"default: CURRENT_TIMESTAMP" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default: CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" form:"updated_at"`
}

// Init the database connection
func init() {
	dbType := conf.ServerCfg.Get(constant.DBType)
	dbHost := conf.ServerCfg.Get(constant.DBHost)
	dbPort := conf.ServerCfg.Get(constant.DBPort)
	dbUser := conf.ServerCfg.Get(constant.DBUser)
	dbPwd := conf.ServerCfg.Get(constant.DBPwd)
	dbName := conf.ServerCfg.Get(constant.DBName)

	db, err := gorm.Open(dbType, fmt.Sprintf(constant.DBConnect, dbUser, dbPwd, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatalln("failed to connect database!")
	}

	db.SingularTable(true)

}
