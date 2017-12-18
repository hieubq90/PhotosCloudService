package main

import (
	"context"
	"fmt"
	"handlers"
	"log"
	"models"
	"net/http"
	"os"
	"os/signal"
	"runtime"
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
		runtime.GOMAXPROCS(app_config.AppConfig.RuntimeMaxProcs)
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

		// make chanel to receive handle new uploaded file need resize
		newFileChanel := make(chan string, 100000)

		// start a RESIZE_HANDLER function
		go func(c chan string) {
			defer close(c)
			for {
				newImage := <-c
				go handlers.DoResize(newImage)
			}
		}(newFileChanel)

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
					ext := utils.GetFileExtension(realFileType)
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
						newImage := fmt.Sprintf("%s/%s", filePath, fileName)
						newFileChanel <- newImage
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


		// Start listen
		endPoint := fmt.Sprintf("%s:%d", app_config.AppConfig.ListenHost, app_config.AppConfig.ListenPort)
		fmt.Printf("[PhotosCloudService] Started listen on: %s\n", endPoint)
		srv := &http.Server{
			Addr:    endPoint,
			Handler: router,
		}

		go func() {
			// service connections
			fmt.Printf("[PhotosCloudService] Started listen on: %s\n", endPoint)
			if err := srv.ListenAndServe(); err != nil {
				fmt.Printf("[PhotosCloudService] Starting error: %s\n", err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Println("Shutdown PhotosCloudService ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("PhotosCloudService Shutdown:", err)
		}
		fmt.Println("PhotosCloudService exiting")
	}
}