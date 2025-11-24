package main

import (
	"fmt"

	"wb_tech/l3_4/internal/domain/entity"
	"wb_tech/l3_4/internal/domain/processing"
)

func main() {
	fmt.Println("Hello, world from L3/4!")
	img := entity.NewImage("image.png", entity.PNG, entity.StatusCompleted)
	resized, err := processing.Resize(img, 100, 100)
	if err != nil {
		fmt.Println("Error resizing image:", err)
		return
	}
	err = processing.Save(img, resized)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}
	fmt.Println("Resized image saved successfully!")
}
