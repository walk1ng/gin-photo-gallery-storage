package utils

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"go.uber.org/zap"

	"github.com/walk1ng/gin-photo-gallery-storage/conf"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"
)

var containerURL azblob.ContainerURL
var azStorageAccountName string
var azStorageAccountKey string
var azStorageContainerName string

func init() {
	azStorageAccountName = conf.ServerCfg.Get(constant.AzStorageAccountName)
	azStorageAccountKey = conf.ServerCfg.Get(constant.AzStorageAccountKey)
	azStorageContainerName = conf.ServerCfg.Get(constant.AzStorageContainerName)
	URL, _ := url.Parse(fmt.Sprintf(constant.AzStorageBlobURLEndpointFormat, azStorageAccountName, azStorageContainerName))

	// create a default request pipeline with storage account name and key
	cred, err := azblob.NewSharedKeyCredential(azStorageAccountName, azStorageAccountKey)
	if err != nil {
		AppLogger.Fatal(fmt.Sprintf("az: invalid credential with error: %s.", err.Error()), zap.String("service", "init()"))
	}

	p := azblob.NewPipeline(cred, azblob.PipelineOptions{})
	containerURL = azblob.NewContainerURL(*URL, p)
}

// Upload func upload a photo to azure blob storage
func Upload(photoID uint, fileName string, file *os.File) string {
	uploadID := fmt.Sprintf(constant.PhotoUpdateIDFormat, photoID)
	go AsyncUpload(uploadID, photoID, fileName, file)
	return uploadID
}

// AsyncUpload func upload a photo to the azure blob storage async
func AsyncUpload(uploadID string, photoID uint, fileName string, file *os.File) {
	// set upload status in redis
	if !SetUploadStatus(uploadID, 1) {
		AppLogger.Info("failed to set upload status before upload.", zap.String("service", "AsyncUpload()"))
		return
	}

	// upload the photo to az blob storage
	blobURL := containerURL.NewBlockBlobURL(fileName)
	_, err := azblob.UploadFileToBlockBlob(context.Background(), file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16,
	})

	// if failed to upload, send callback to redis to delete photo
	if err != nil {
		AppLogger.Info(err.Error(), zap.String("service", "AsyncUpload()"))
		if !SendToChannel(constant.PhotoDeleteChannel, fmt.Sprintf("%d", photoID)) {
			AppLogger.Info("failed to send delete-photo message to channel.", zap.String("service", "AsyncUpload()"))
		}
		return
	}

	// if success to upload, send callback to redis to update url for the photo
	photoURL := fmt.Sprintf(constant.AzStorageBlobURLEndpointFormat, azStorageAccountName, azStorageContainerName) + "/" + fileName
	updateURLMessage := fmt.Sprintf("%d-%s", photoID, photoURL)
	if !SendToChannel(constant.PhotoURLUpdateChannel, updateURLMessage) {
		AppLogger.Info("failed to send update-photo-url message to channel.", zap.String("service", "AsyncUpload()"))
	}
	return
}
