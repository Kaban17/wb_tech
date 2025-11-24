package processing

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	img_ "wb_tech/l3_4/internal/domain/entity"

	"github.com/disintegration/imaging"
)

func Resize(i *img_.Image, width, height int) (*image.Image, error) {
	file, err := os.Open(i.OriginalName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	resized := imaging.Resize(img, width, height, imaging.Lanczos)

	var result image.Image = resized
	return &result, nil
}
