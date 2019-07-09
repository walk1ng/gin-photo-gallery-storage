package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/walk1ng/gin-photo-gallery-storage/utils"
	"go.uber.org/zap"

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

	var err error
	db, err = gorm.Open(dbType, fmt.Sprintf(constant.DBConnect, dbUser, dbPwd, dbHost, dbPort, dbName))
	if err != nil {
		utils.AppLogger.Fatal(err.Error(), zap.String("service", "init()"))
	}

	db.SingularTable(true)
	if !db.HasTable(&Auth{}) {
		db.CreateTable(&Auth{})
	}

	if !db.HasTable(&Bucket{}) {
		db.CreateTable(&Bucket{})
	}

	if !db.HasTable(&Photo{}) {
		db.CreateTable(&Photo{})
	}

	// run a goroutine never exit to listen to redis callbacks
	go listenRedisCallback()

}

func listenRedisCallback() {
	// wait until the utils package is initialized
	<-utils.InitComplete

	// subscribe the redis channels
	updateChan := utils.RedisClient.Subscribe(constant.PhotoURLUpdateChannel).Channel()
	deleteChan := utils.RedisClient.Subscribe(constant.PhotoDeleteChannel).Channel()

	for {
		select {
		case msg := <-updateChan:
			photoID, _ := strconv.Atoi(msg.Payload[:strings.Index(msg.Payload, "-")])
			photoURL := msg.Payload[strings.Index(msg.Payload, "-")+1:]
			dberr := UpdatePhotoURL(uint(photoID), photoURL)
			if dberr != nil {
				utils.AppLogger.Info("callback error: update photo url.", zap.String("service", "listenRedisCallback()"))
			} else {
				utils.SetUploadStatus(fmt.Sprintf(constant.PhotoUpdateIDFormat, photoID), 0)
			}

		case msg := <-deleteChan:
			photoID, _ := strconv.Atoi(msg.Payload)
			if err := DeletePhotoByID(uint(photoID)); err != nil {
				utils.AppLogger.Info("callback error: delete photo.", zap.String("service", "listenRedisCallback()"))
			} else {
				utils.SetUploadStatus(fmt.Sprintf(constant.PhotoUpdateIDFormat, photoID), -1)
			}
		}
	}
}
