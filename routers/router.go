package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/walk1ng/gin-photo-gallery-storage/apis/v1"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()

	v1Group := Router.Group("/api/v1")
	{
		// auth
		authGroup := v1Group.Group("/auth")
		{
			authGroup.POST("/add", v1.AddAuth)
			authGroup.POST("/check", v1.CheckAuth)
		}
	}
}
