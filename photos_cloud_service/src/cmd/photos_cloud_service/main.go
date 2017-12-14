package main

import (
	"fmt"
	"models"
	"net/http"
	"strings"
	"time"
	"utils"

	"app_config"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func main() {
	if app_config.InitFromYAML() {
		// setup Gin run mode
		if strings.Compare(app_config.AppConfig.RunMode, "prod") == 0 {
			gin.SetMode(gin.ReleaseMode)
		} else {
			gin.SetMode(gin.DebugMode)
		}


		router := gin.Default()

		// Setup CORS settings
		router.Use(cors.New(cors.Config{
			AllowOrigins:     app_config.AppConfig.AllowOrigins,
			AllowMethods:     app_config.AppConfig.AllowMethods,
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

		// Setup memory limit for multipart forms (default is 32 MiB)
		//router.MaxMultipartMemory = int64(app_config.AppConfig.BodyLimitSize << 20) // MiB
		MaxSize := int64(app_config.AppConfig.BodyLimitSize << 20)
		router.Use(limits.RateLimiter(MaxSize))

		// Handle simple upload client
		if app_config.AppConfig.HandleDemo {
			router.Static("/", "./public")
		}

		// Handle Upload request
		router.POST("/upload", func(c *gin.Context) {
			// 1. Create Uploaded Response Data
			uploadedFiles := make([]string, 0)
			// Multipart form
			form, err := c.MultipartForm()
			if err != nil {
				uploadedFiles = append(uploadedFiles, fmt.Sprintf("get form err: %s", err.Error()))
				c.JSON(http.StatusBadRequest, models.UploadedFiles{
					Files: uploadedFiles,
				})
				return
			}
			files := form.File["files"]

			for _, file := range files {
				realFileType := utils.GetFileType(file)
				if app_config.AppConfig.IsAllowedFileType(realFileType) {
					// allowed to upload
					// get file extension
					ext := utils.GetFileExtension(file.Filename)
					// 1. generate new name
					fileName := fmt.Sprintf("%s.%s", uuid.NewV4(), ext)
					folderName := utils.GetCurrentFormatDate()
					filePath := fmt.Sprintf("%s%s", app_config.AppConfig.SaveLocation, folderName)
					// 2. save file
					// 2.1. create file path
					utils.CreateDirIfNotExist(filePath)
					// 2.2. save to disk
					err := c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", filePath, fileName))
					if err != nil {
						//c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
						uploadedFiles = append(uploadedFiles, fmt.Sprintf("get form err: %s", err.Error()))
					} else {
						uploadedFiles = append(uploadedFiles, fmt.Sprintf("%s%s/%s", app_config.AppConfig.DownloadDomain, folderName, fileName))
					}
				} else {
					// not allowed to upload
					uploadedFiles = append(uploadedFiles, "file_type is not allowed")
				}
			}
			c.JSON(http.StatusOK, models.UploadedFiles{
				Files: uploadedFiles,
			})
		})
	}
}