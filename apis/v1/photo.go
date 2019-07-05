package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/walk1ng/gin-photo-gallery-storage/utils"
	"go.uber.org/zap"

	"github.com/astaxie/beego/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"
	"github.com/walk1ng/gin-photo-gallery-storage/models"
)

// AddPhoto func add a new photo.
func AddPhoto(context *gin.Context) {
	responseCode := constant.InvalidParams
	photoToAdd := models.Photo{}

	photoFile, fileErr := context.FormFile("photo")
	if fileErr != nil {
		utils.AppLogger.Info(fileErr.Error(), zap.String("service", "AddPhoto()"))
	}

	paramErr := context.ShouldBindWith(&photoToAdd, binding.Form)
	if paramErr != nil {
		utils.AppLogger.Info(paramErr.Error(), zap.String("service", "AddPhoto()"))
	}

	if fileErr != nil || paramErr != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(photoToAdd.AuthID, "auth_id").Message("must have auth id")
	validCheck.Required(photoToAdd.BucketID, "bucket_id").Message("must have bucket id")
	validCheck.Required(photoToAdd.Name, "photo_name").Message("must have photo name")
	validCheck.MaxSize(photoToAdd.Name, 255, "photo_name").Message("length of photo's name cannot exceed 255")

	data := make(map[string]interface{})
	photoToAdd.Tag = strings.Join(photoToAdd.Tags, ";")

	if !validCheck.HasErrors() {
		if photoToAdd, uploadID, err := models.AddPhoto(&photoToAdd, photoFile); err != nil {
			if err == models.ErrPhotoExists {
				responseCode = constant.PhotoAlreadyExist
			} else {
				responseCode = constant.InternalServerError
			}
		} else {
			responseCode = constant.PhotoAddInProcess
			data["photo"] = *photoToAdd
			data["photo_upload_id"] = uploadID
		}
	} else {
		for _, e := range validCheck.Errors {
			utils.AppLogger.Info(e.Message, zap.String("service", "AddPhoto()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// DeletePhoto func delete an existed photo.
func DeletePhoto(context *gin.Context) {
	responseCode := constant.InvalidParams
	bucketID, err := strconv.Atoi(context.PostForm("bucket_id"))
	photoName := context.PostForm("photo_name")
	if err != nil {
		utils.AppLogger.Info(err.Error(), zap.String("service", "DeletePhoto()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(bucketID, "bucket_id").Message("must have bucket id")
	validCheck.Required(photoName, "photo_name").Message("must have photo name")
	validCheck.MaxSize(photoName, 255, "photo_name").Message("length of photo's name cannot exceed 255")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if err := models.DeletePhotoByBucketIDAndPhotoName(uint(bucketID), photoName); err != nil {
			if err == models.ErrNoSuchPhoto {
				responseCode = constant.PhotoNotExist
			} else {
				responseCode = constant.InternalServerError
			}
		} else {
			responseCode = constant.PhotoDeleteSuccess
		}
	} else {
		for _, e := range validCheck.Errors {
			utils.AppLogger.Info(e.Message, zap.String("service", "DeletePhoto()"))
		}
	}

	data["photo_name"] = photoName
	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// UpdatePhoto func update an existed photo.
func UpdatePhoto(context *gin.Context) {
	responseCode := constant.InvalidParams
	photoToUpdate := models.Photo{}

	err := context.ShouldBindWith(&photoToUpdate, binding.Form)
	if err != nil {
		utils.AppLogger.Info(err.Error(), zap.String("service", "UpdatePhoto()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(photoToUpdate.ID, "photo_id").Message("must have photo id")
	validCheck.Min(int(photoToUpdate.ID), 1, "photo_id").Message("photo id must be positive")
	validCheck.Required(photoToUpdate.Name, "photo_name").Message("must have photo name")
	validCheck.MaxSize(photoToUpdate.Name, 255, "photo_name").Message("length of photo's name cannot exceed 255")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		photoToUpdate.Tag = strings.Join(photoToUpdate.Tags, ";")
		if photo, err := models.UpdatePhoto(&photoToUpdate); err != nil {
			if err == models.ErrNoSuchPhoto {
				responseCode = constant.PhotoNotExist
			} else {
				responseCode = constant.InternalServerError
			}
		} else {
			responseCode = constant.PhotoUpdateSuccess
			data["photo"] = *photo
		}
	} else {
		for _, e := range validCheck.Errors {
			utils.AppLogger.Info(e.Message, zap.String("service", "UpdatePhoto()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// GetPhotoByID func get photo by its ID.
func GetPhotoByID(context *gin.Context) {
	responseCode := constant.InvalidParams
	photoID, err := strconv.Atoi(context.Query("photo_id"))

	if err != nil {
		utils.AppLogger.Info(err.Error(), zap.String("service", "GetPhotoByID()"))
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}
	validCheck.Required(photoID, "photo_id").Message("must have photo id")
	validCheck.Min(photoID, 1, "photo_id").Message("photo id should be positive")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if photo, err := models.GetPhotoByID(uint(photoID)); err != nil {
			if err == models.ErrNoSuchPhoto {
				responseCode = constant.PhotoNotExist
			} else {
				responseCode = constant.InternalServerError
			}
		} else {
			responseCode = constant.PhotoGetSuccess
			photo.Tags = strings.Split(photo.Tag, ";")
			data["photo"] = *photo
		}
	} else {
		for _, e := range validCheck.Errors {
			utils.AppLogger.Info(e.Message, zap.String("service", "GetPhotoByID()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// GetPhotoByBucketID func get photos by bucket ID.
func GetPhotoByBucketID(context *gin.Context) {
	responseCode := constant.InvalidParams
	bucketID, err := strconv.Atoi(context.Query("bucket_id"))
	offset := context.GetInt("offset")

	if err != nil {
		utils.AppLogger.Info(err.Error(), zap.String("service", "GetPhotoByBucketID()"))
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
	validCheck.Min(offset, 0, "page_offset").Message("page offset must be >= 0")

	data := make(map[string]interface{})
	if !validCheck.HasErrors() {
		if photos, err := models.GetPhotosByBucketID(uint(bucketID), offset); err != nil {
			responseCode = constant.InternalServerError
		} else {
			responseCode = constant.PhotoGetSuccess
			for i, photo := range photos {
				photos[i].Tags = strings.Split(photo.Tag, ";")
			}
			data["photo"] = photos
		}
	} else {
		for _, e := range validCheck.Errors {
			utils.AppLogger.Info(e.Message, zap.String("service", "GetPhotoByBucketID()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// GetPhotoUploadStatus func get the upload status of photo by photo ID.
func GetPhotoUploadStatus(context *gin.Context) {
	responseCode := constant.InvalidParams
	uploadID := context.Query("upload_id")

	validCheck := validation.Validation{}
	validCheck.Required(uploadID, "upload_id").Message("must have upload id")

	data := make(map[string]interface{})
	data["upload_id"] = uploadID

	if !validCheck.HasErrors() {
		responseCode = models.GetPhotoUploadStatus(uploadID)
	} else {
		for _, e := range validCheck.Errors {
			utils.AppLogger.Info(e.Message, zap.String("service", "GetPhotoUploadStatus()"))
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}
