package main

import (
	"fmt"
	"wb_tech/l3_4/internal/domain/entity"
	"wb_tech/l3_4/internal/domain/processing"
)

func main() {
	fmt.Println("Hello, world from L3/4!")

	img := entity.NewImage("image.png", entity.PNG, entity.StatusCompleted)

	pipeline1 := processing.NewPipeline().
		Add(processing.ResizeStep(800, 600)).
		Add(processing.BlurStep(5.0))

	result1, err := processing.ProcessWithPipeline(img, pipeline1)
	if err != nil {
		fmt.Println("Error processing image:", err)
		return
	}

	err = processing.SaveImage(img, result1)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Println("Pipeline 1: Resized and blurred image saved!")

}
