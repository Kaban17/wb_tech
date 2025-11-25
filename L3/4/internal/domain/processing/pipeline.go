package processing

import (
	"image"
	"os"

	"wb_tech/l3_4/internal/domain/entity"

	"github.com/disintegration/imaging"
)

// ProcessFunc - функция обработки изображения
type ProcessFunc func(image.Image) (image.Image, error)

// Pipeline - цепочка обработки изображений
type Pipeline struct {
	steps []ProcessFunc
}

// NewPipeline создает новый pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{
		steps: make([]ProcessFunc, 0),
	}
}

// Add добавляет шаг обработки в pipeline
func (p *Pipeline) Add(step ProcessFunc) *Pipeline {
	p.steps = append(p.steps, step)
	return p
}

// Execute выполняет все шаги pipeline последовательно
func (p *Pipeline) Execute(img image.Image) (image.Image, error) {
	var err error
	result := img

	for _, step := range p.steps {
		result, err = step(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// LoadImage загружает изображение из файла
func LoadImage(i *entity.Image) (image.Image, error) {
	file, err := os.Open(i.OriginalName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// ProcessWithPipeline загружает изображение и применяет pipeline
func ProcessWithPipeline(i *entity.Image, pipeline *Pipeline) (image.Image, error) {
	img, err := LoadImage(i)
	if err != nil {
		return nil, err
	}

	return pipeline.Execute(img)
}

// ResizeStep создает шаг изменения размера
func ResizeStep(width, height int) ProcessFunc {
	return func(img image.Image) (image.Image, error) {
		resized := imaging.Resize(img, width, height, imaging.Lanczos)
		return resized, nil
	}
}

// BlurStep создает шаг размытия
func BlurStep(sigma float64) ProcessFunc {
	return func(img image.Image) (image.Image, error) {
		blurred := imaging.Blur(img, sigma)
		return blurred, nil
	}
}

// GrayscaleStep создает шаг преобразования в градации серого
func GrayscaleStep() ProcessFunc {
	return func(img image.Image) (image.Image, error) {
		grayscale := imaging.Grayscale(img)
		return grayscale, nil
	}
}

// CropStep создает шаг обрезки изображения
func CropStep(rect image.Rectangle) ProcessFunc {
	return func(img image.Image) (image.Image, error) {
		cropped := imaging.Crop(img, rect)
		return cropped, nil
	}
}
