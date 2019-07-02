package middlewares

import (
	"log"
	"net/http"

	"github.com/walk1ng/gin-photo-gallery-storage/utils"

	"github.com/gin-gonic/gin"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"
)

// GetAuthMiddleware func is a wrapper func to return a auth middleware
func GetAuthMiddleware() func(*gin.Context) {
	return func(context *gin.Context) {
		jwtString, err := context.Cookie(constant.Jwt)

		// cannot find the jwt in cookie, it mean that the user has not loggined yet or cookie missing.
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{
				"code": constant.JwtMissingError,
				"data": make(map[string]string),
				"msg":  constant.GetMessage(constant.JwtMissingError),
			})
			context.Abort()
			return
		}

		// parse jwt
		claim, err := utils.ParseJWT(jwtString)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{
				"code": constant.JwtParseError,
				"data": make(map[string]string),
				"msg":  constant.GetMessage(constant.JwtParseError),
			})
			context.Abort()
			return
		}

		if utils.IsAuthInRedis(claim.UserName) {
			context.Set("user_name", claim.UserName)
			context.Next()
		} else {
			// auth is expired
			context.JSON(http.StatusBadRequest, gin.H{
				"code": constant.UserAuthTimeout,
				"data": make(map[string]string),
				"msg":  constant.GetMessage(constant.UserAuthTimeout),
			})
			context.Abort()
		}
	}
}
