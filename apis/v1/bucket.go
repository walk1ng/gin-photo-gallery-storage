package v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"
	"github.com/walk1ng/gin-photo-gallery-storage/models"
)

// AddBucket func add a new bucket.
func AddBucket(context *gin.Context) {
	responseCode := constant.InvalidParams
	bucketToAdd := models.Bucket{}
	if err := context.ShouldBindWith(&bucketToAdd, binding.Form); err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketToAdd.AuthID, "auth_id").Message("must have auth id")
	validCheck.Required(bucketToAdd.Name, "bucket_name").Message("must have bucket name")
	validCheck.MaxSize(bucketToAdd.Name, 64, "bucket_name").Message("length of bucket name cannot exceed 64")

	if !validCheck.HasErrors() {
		if err := models.AddBucket(&bucketToAdd); err != nil {
			if err == models.ErrBucketExists {
				responseCode = constant.BucketAlreadyExist
			} else {
				responseCode = constant.InternalServerError
			}
		} else {
			responseCode = constant.BucketAddSuccess
		}
	} else {
		for _, e := range validCheck.Errors {
			log.Println(e.Message)
		}
	}

	data := make(map[string]string)
	data["bucket_name"] = bucketToAdd.Name

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// DeleteBucket func delete an exist bucket.
func DeleteBucket(context *gin.Context) {
	responseCode := constant.InvalidParams
	bucketID, err := strconv.Atoi(context.Query("bucket_id"))
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketID, "bucket_id").Message("must have bucket id")
	validCheck.Min(bucketID, 1, "bucket_id").Message("bucket id should be positive")

	if !validCheck.HasErrors() {
		if err := models.DeleteBucket(uint(bucketID)); err != nil {
			if err == models.ErrNoSuchBucket {
				responseCode = constant.BucketNotExist
			} else {
				responseCode = constant.InternalServerError
			}
		} else {
			responseCode = constant.BucketDeleteSuccess
		}
	} else {
		for _, e := range validCheck.Errors {
			log.Println(e.Message)
		}
	}

	data := make(map[string]interface{})
	data["bucket_id"] = bucketID
	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// UpdateBucket func to update an existed bucket.
func UpdateBucket(context *gin.Context) {
	responseCode := constant.InvalidParams
	bucketToUpdate := models.Bucket{}
	if err := context.ShouldBindWith(&bucketToUpdate, binding.Form); err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketToUpdate.ID, "bucket_id").Message("must have bucket id")
	validCheck.MaxSize(bucketToUpdate.Name, 64, "bucket_name").Message("name of bucket cannot exceed 64")

	if !validCheck.HasErrors() {
		if err := models.UpdateBucket(&bucketToUpdate); err != nil {
			if err == models.ErrNoSuchBucket {
				responseCode = constant.BucketNotExist
			} else {
				responseCode = constant.InternalServerError
			}
		} else {
			responseCode = constant.BucketUpdateSuccess
		}
	} else {
		for _, e := range validCheck.Errors {
			log.Println(e.Message)
		}
	}

	data := make(map[string]interface{})
	data["bucket_id"] = bucketToUpdate.ID
	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// GetBucketByID func get bucket by its ID.
func GetBucketByID(context *gin.Context) {
	responseCode := constant.InvalidParams
	bucketID, err := strconv.Atoi(context.Query("bucket_id"))
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketID, "bucket_id").Message("must have bucket id")
	validCheck.Min(bucketID, 1, "bucket_id").Message("bucket id should be positive")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if bucket, err := models.GetBucketByID(uint(bucketID)); err != nil {
			if err == models.ErrNoSuchBucket {
				responseCode = constant.BucketNotExist
			} else {
				responseCode = constant.InternalServerError
			}
		} else {
			responseCode = constant.BucketGetSuccess
			data["bucket"] = bucket
		}
	} else {
		for _, e := range validCheck.Errors {
			log.Println(e.Message)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// GetBucketByAuthID func get buckets by auth id.
func GetBucketByAuthID(context *gin.Context) {
	responseCode := constant.InvalidParams
	authID, err := strconv.Atoi(context.Query("auth_id"))
	offset := context.GetInt("offset")
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(authID, "auth_id").Message("must have auth id")
	validCheck.Min(authID, 1, "auth_id").Message("auth id should be positive")
	validCheck.Min(offset, 0, "page_offset").Message("page offset must be LE than 0")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if buckets, err := models.GetBucketByAuthID(uint(authID), offset); err != nil {
			responseCode = constant.InternalServerError
		} else {
			responseCode = constant.BucketGetSuccess
			data["buckets"] = buckets
		}
	} else {
		for _, e := range validCheck.Errors {
			log.Println(e.Message)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})

}
