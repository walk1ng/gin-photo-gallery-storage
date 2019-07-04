package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/walk1ng/gin-photo-gallery-storage/apis/v1"
	"github.com/walk1ng/gin-photo-gallery-storage/middlewares"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()

	authMiddleware := middlewares.GetAuthMiddleware()
	refreshMiddleware := middlewares.GetRefreshMiddleware()
	paginationMiddleware := middlewares.GetPaginationMiddleware()

	v1Group := Router.Group("/api/v1")
	{
		// auth
		authGroup := v1Group.Group("/auth")
		{
			authGroup.POST("/add", v1.AddAuth)
			authGroup.POST("/check", v1.CheckAuth)
		}

		// bucket
		bucketGroup := v1Group.Group("/bucket")
		{
			bucketGroup.POST("/add", authMiddleware, refreshMiddleware, v1.AddBucket)
			bucketGroup.DELETE("/delete", authMiddleware, refreshMiddleware, v1.DeleteBucket)
			bucketGroup.PUT("/update", authMiddleware, refreshMiddleware, v1.UpdateBucket)
			bucketGroup.GET("/get_by_id", authMiddleware, refreshMiddleware, v1.GetBucketByID)
			bucketGroup.GET("/get_by_auth_id", authMiddleware, refreshMiddleware, paginationMiddleware, v1.GetBucketByAuthID)
		}

		// photo
		photoGroup := v1Group.Group("/photo")
		{
			photoGroup.POST("/add", authMiddleware, refreshMiddleware, v1.AddPhoto)
			photoGroup.DELETE("/delete", authMiddleware, refreshMiddleware, v1.DeletePhoto)
			photoGroup.PUT("/update", authMiddleware, refreshMiddleware, v1.UpdatePhoto)
			photoGroup.GET("/get_by_id", authMiddleware, refreshMiddleware, v1.GetPhotoByID)
			photoGroup.GET("/get_by_bucket_id", authMiddleware, refreshMiddleware, paginationMiddleware, v1.GetPhotoByBucketID)
			photoGroup.GET("/upload_status", authMiddleware, refreshMiddleware, v1.GetPhotoUploadStatus)
		}
	}
}
