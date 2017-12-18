package handlers

import (
	"app_config"
	"fmt"
	"utils"

	"gopkg.in/h2non/bimg.v1"
)

var (
	options320 = bimg.Options{
		Width:        320,
		Crop:         false,
		Compression:  9,
		Embed: true,
	}

	options720 = bimg.Options{
		Width:        720,
		Crop:         false,
		Compression:  9,
		Embed: true,
	}
)

func DoResize(img string) {
	buffer, err := bimg.Read(img)

	if err != nil {
		fmt.Printf("[DoResize] Load file error: %s\n", err.Error())
	}
	size, _ := bimg.NewImage(buffer).Size()
	filePath, fileName := utils.GetFilePathAndName(img)

	for _, opt := range app_config.AppConfig.Resize_Options {
		if size.Width > opt {
			if opt == 320 {
				newImage, err := bimg.NewImage(buffer).Process(options320)
				if err != nil {
					fmt.Printf("[DoResize] Resize 320 error: %s\n", err.Error())
				}

				utils.CreateDirIfNotExist(filePath+"320")
				bimg.Write(fmt.Sprintf("%s320%s", filePath, fileName), newImage)
			} else if opt == 720 {
				newImage, err := bimg.NewImage(buffer).Process(options720)
				if err != nil {
					fmt.Printf("[DoResize] Resize 720 error: %s\n", err.Error())
				}

				utils.CreateDirIfNotExist(filePath+"720")
				bimg.Write(fmt.Sprintf("%s720%s", filePath, fileName), newImage)
			}
		}
	}
}

//func DoResize(img string) {
//	img := image.(image.Rect(0, 0, tt.origWidth, tt.origHeight))
//}