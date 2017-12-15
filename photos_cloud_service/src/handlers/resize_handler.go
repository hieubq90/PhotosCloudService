package handlers

import (
	"fmt"

	"gopkg.in/h2non/bimg.v1"
)

func DoResize(img string) {
	buffer, err := bimg.Read(img)
	if err != nil {
		fmt.Printf("[DoResize] Load file error: %s\n", err.Error())
	}

	newImage, err := bimg.NewImage(buffer).Resize(320, 320)
	if err != nil {
		fmt.Printf("[DoResize] Resize error: %s\n", err.Error())
	}


	bimg.Write("new.jpg", newImage)
}