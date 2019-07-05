package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/walk1ng/gin-photo-gallery-storage/conf"

	"github.com/gin-gonic/gin"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"
	"github.com/walk1ng/gin-photo-gallery-storage/utils"
)

// GetRefreshMiddleware func is a wrapper func to return a refresh middleware
func GetRefreshMiddleware() func(*gin.Context) {
	return func(context *gin.Context) {
		// firstly, get the user_name which set by the auth middleware
		if userName, exist := context.Get("user_name"); exist {
			// generate a new jwt for the user
			jwtString, err := utils.GenerateJWT(userName.(string))
			if err != nil {
				utils.AppLogger.Info(err.Error(), zap.String("service", "GetRefreshMiddleware()"))
				data := make(map[string]string)
				data["user_name"] = userName.(string)
				context.JSON(http.StatusBadRequest, gin.H{
					"code": constant.JwtGenerationError,
					"data": data,
					"msg":  constant.GetMessage(constant.JwtGenerationError),
				})
				context.Abort()
				return
			}

			// save the new jwt in user's cookie
			context.SetCookie(constant.Jwt, jwtString,
				constant.CookieMaxAge,
				conf.ServerCfg.Get(constant.ServerPath),
				conf.ServerCfg.Get(constant.ServerDomain),
				true, true)

			// refresh user in the redis, it mean to refresh the key's expiration
			err = utils.AddAuthToRedis(userName.(string))
			if err != nil {
				utils.AppLogger.Info(err.Error(), zap.String("service", "GetRefreshMiddleware()"))
				context.JSON(http.StatusBadRequest, gin.H{
					"code": constant.InternalServerError,
					"data": make(map[string]string),
					"msg":  constant.GetMessage(constant.InternalServerError),
				})
				context.Abort()
				return
			}

			context.Next()
		}

		context.Abort()
	}
}
