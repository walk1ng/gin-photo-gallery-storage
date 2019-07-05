package utils

import (
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
	"github.com/walk1ng/gin-photo-gallery-storage/conf"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"
)

var RedisClient *redis.Client
var InitComplete = make(chan struct{}, 1)

// Initialize the redis client
func init() {
	host := conf.ServerCfg.Get(constant.RedisHost)
	port := conf.ServerCfg.Get(constant.RedisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})
	InitComplete <- struct{}{}
}

// AddAuthToRedis func add an auth to redis mean the user has logged in
func AddAuthToRedis(username string) error {
	key := fmt.Sprintf("%s%s", constant.LoginUser, username)
	err := RedisClient.Set(key, username, constant.LoginMaxAge*time.Second).Err()
	if err != nil {
		AppLogger.Info(err.Error(), zap.String("service", "AddAuthToRedis()"))
		return err
	}
	return nil
}

// IsAuthInRedis func check if an auth exists in redis
func IsAuthInRedis(username string) bool {
	key := fmt.Sprintf("%s%s", constant.LoginUser, username)
	err := RedisClient.Get(key).Err()
	if err != nil {
		AppLogger.Info(err.Error(), zap.String("service", "IsAuthInRedis()"))
		return false
	}
	return true
}

// RemoveAuthFromRedis func remove the auth from redis
func RemoveAuthFromRedis(username string) bool {
	key := fmt.Sprintf("%s%s", constant.LoginUser, username)
	err := RedisClient.Del(key).Err()
	if err != nil {
		AppLogger.Info(err.Error(), zap.String("service", "RemoveAuthFromRedis()"))
		return false
	}
	return true
}

// SetUploadStatus func set the upload status for a photo
func SetUploadStatus(key string, value int) bool {
	err := RedisClient.Set(key, value, 0).Err()
	if err != nil {
		AppLogger.Info(err.Error(), zap.String("service", "SetUploadStatus()"))
		return false
	}
	return true
}

// GetPhotoUploadStatus func get the upload status for a photo
func GetUploadStatus(key string) int {
	val := RedisClient.Get(key).Val()
	if val == "" {
		return -2 // no such key
	}
	status, _ := strconv.Atoi(val)
	return status
}

// SendToChannel func send a message to a channel
func SendToChannel(channel, message string) bool {
	err := RedisClient.Publish(channel, message).Err()
	if err != nil {
		AppLogger.Info(err.Error(), zap.String("service", "SendToChannel()"))
		return false
	}
	return true
}
