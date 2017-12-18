package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	IMG_TYPE_JPEG = "image/jpeg"
	IMG_TYPE_PNG = "image/png"
)

func GetFileType(fileHeader *multipart.FileHeader) string {
	file, _ := fileHeader.Open() // or get your file from a file system
	defer file.Close()
	buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	if _, err := file.Read(buff); err != nil {
		fmt.Println(err) // do something with that error
		return ""
	}
	return http.DetectContentType(buff)
}

func GetFileExtension(fileType string) string {
	if strings.Compare(fileType, IMG_TYPE_JPEG) == 0 {
		return "jpg"
	}
	if strings.Compare(fileType, IMG_TYPE_PNG) == 0 {
		return "png"
	}
	return "unknown"
}

func GetFilePathAndName(filename string) (string, string) {
	index := strings.LastIndex(filename, "/")
	return filename[:index+1], filename[index:]
}

func GetCurrentFormatDate() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
