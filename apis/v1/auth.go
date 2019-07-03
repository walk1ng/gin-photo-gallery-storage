package v1

import (
	"log"
	"net/http"

	"github.com/walk1ng/gin-photo-gallery-storage/conf"

	"github.com/walk1ng/gin-photo-gallery-storage/utils"

	"github.com/walk1ng/gin-photo-gallery-storage/models"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"
)

// AddAuth func to add a new auth
func AddAuth(context *gin.Context) {
	userName := context.PostForm("user_name")
	password := context.PostForm("password")
	email := context.PostForm("email")

	// validations
	validCheck := validation.Validation{}
	validCheck.Required(userName, "user_name").Message("must have user name")
	validCheck.MaxSize(userName, 16, "user_name").Message("length of user name cannot exceed 16")
	validCheck.MinSize(userName, 6, "user_name").Message("length of user name is at least 6")
	validCheck.Required(password, "password").Message("must have password")
	validCheck.MaxSize(password, 16, "password").Message("length of password cannot exceed 16")
	validCheck.MinSize(password, 6, "password").Message("length of password is at least 6")
	validCheck.Required(email, "email").Message("must have email")
	validCheck.MaxSize(email, 128, "email").Message("email can not exceed 128 chars")

	responseCode := constant.InvalidParams
	if !validCheck.HasErrors() {
		if err := models.AddAuth(userName, password, email); err == nil {
			responseCode = constant.UserAddSuccess
		} else {
			responseCode = constant.UserAlreadyExist
		}
	} else {
		for _, e := range validCheck.Errors {
			log.Println(e.Message)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": userName,
		"msg":  constant.GetMessage(responseCode),
	})
}

// CheckAuth func check if the auth is valid
func CheckAuth(context *gin.Context) {
	userName := context.PostForm("user_name")
	password := context.PostForm("password")

	// validation
	validCheck := validation.Validation{}
	validCheck.Required(userName, "user_name").Message("must have user name")
	validCheck.MaxSize(userName, 16, "user_name").Message("length of user name cannot exceed 16")
	validCheck.MinSize(userName, 6, "user_name").Message("length of user name is at least 6")
	validCheck.Required(password, "password").Message("must have password")
	validCheck.MaxSize(password, 16, "password").Message("length of password cannot exceed 16")
	validCheck.MinSize(password, 6, "password").Message("length of password is at least 6")

	responseCode := constant.InvalidParams

	if !validCheck.HasErrors() {
		if models.CheckAuth(userName, password) {
			if jwtString, err := utils.GenerateJWT(userName); err != nil {
				responseCode = constant.JwtGenerationError
			} else {
				// auth check is pass
				// 1. set jwt to user's cookie
				// 2. add user to redis
				context.SetCookie(constant.Jwt, jwtString,
					constant.CookieMaxAge,
					conf.ServerCfg.Get(constant.ServerPath),
					conf.ServerCfg.Get(constant.ServerDomain),
					true, true)
				if err := utils.AddAuthToRedis(userName); err != nil {
					responseCode = constant.InternalServerError
				} else {
					responseCode = constant.UserAuthSuccess
				}
			}
		} else {
			responseCode = constant.UserAuthError
		}
	} else {
		for _, e := range validCheck.Errors {
			log.Println(e.Message)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": userName,
		"msg":  constant.GetMessage(responseCode),
	})
}
