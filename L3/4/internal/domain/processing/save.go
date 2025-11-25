package processing

import (
	"image"
	"os"
	"wb_tech/l3_4/internal/domain/entity"

	"github.com/disintegration/imaging"
)

func ChooseFormat(i *entity.Image) imaging.Format {
	switch i.Format {
	case entity.JPEG:
		return imaging.JPEG
	case entity.PNG:
		return imaging.PNG
	case entity.GIF:
		return imaging.GIF
	default:
		return imaging.BMP
	}
}

// SaveImage сохраняет обработанное изображение
func SaveImage(i *entity.Image, img image.Image) error {
	filename := i.ID + "_" + i.OriginalName
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = imaging.Encode(file, img, ChooseFormat(i))
	if err != nil {
		return err
	}

	return nil
}

// SaveImageToPath сохраняет изображение по указанному пути
func SaveImageToPath(path string, img image.Image, format imaging.Format) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return imaging.Encode(file, img, format)
}
