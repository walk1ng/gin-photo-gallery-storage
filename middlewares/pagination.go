package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/walk1ng/gin-photo-gallery-storage/constant"

	"github.com/gin-gonic/gin"
)

var errInvalidPageNo = errors.New("page no can not be negative")

// GetPaginationMiddleware func is a wrapper function to return a pagination middleware
func GetPaginationMiddleware() func(*gin.Context) {
	return func(context *gin.Context) {
		responseCode := constant.PaginationSuccess
		pageNo := context.Query("page")
		if pageNo == "" {
			responseCode = constant.InvalidParams
		} else {
			pageOffset, err := getPagination(pageNo)
			if err != nil {
				responseCode = constant.InvalidParams
			} else {
				context.Set("offset", pageOffset)
			}
		}

		if responseCode == constant.InvalidParams {
			data := make(map[string]string)
			data["page"] = pageNo
			context.JSON(http.StatusBadRequest, gin.H{
				"code": responseCode,
				"data": data,
				"msg":  constant.GetMessage(responseCode),
			})
			context.Abort()
		}

		// forward to the next middleware
		context.Next()
	}
}

// getPagination func which calculates the offset given the page number
func getPagination(pageNo string) (int, error) {
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		return 0, err
	}
	if pageNoInt < 0 {
		return 0, errInvalidPageNo
	}
	return pageNoInt * constant.PageSize, nil
}
