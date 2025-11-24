package processing

import (
	"image"
	"os"
	"wb_tech/l3_4/internal/domain/entity"
	img_ "wb_tech/l3_4/internal/domain/entity"

	"github.com/disintegration/imaging"
)

func chooseFormat(i *img_.Image) imaging.Format {
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

func Save(i *img_.Image, resized *image.Image) error {
	filename := i.ID + "_" + i.OriginalName
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = imaging.Encode(file, *resized, imaging.Format(chooseFormat(i)))
	if err != nil {
		return err
	}

	return nil
}
