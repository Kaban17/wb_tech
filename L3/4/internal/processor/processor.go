package processor

import (
	"context"
	"fmt"
	"image"
	"path/filepath"
	"wb_tech/l3_4/internal/model"

	"github.com/disintegration/imaging"
)

type Processor struct {
	basePath string
}

func NewProcessor(basePath string) *Processor {
	return &Processor{basePath: basePath}
}

func (p *Processor) ProcessImage(ctx context.Context, img *model.Image) (string, error) {
	originalPath := filepath.Join(p.basePath, img.OriginalPath)
	
	src, err := imaging.Open(originalPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}

	// 1. Resize/Thumbnail (e.g., fit into 800x800)
	processed := imaging.Fit(src, 800, 800, imaging.Lanczos)

	// 2. Add Watermark
	finalImage := p.addWatermark(processed)

	processedFileName := "processed_" + img.ID + ".jpg"
	processedPath := filepath.Join(p.basePath, processedFileName)

	if err := imaging.Save(finalImage, processedPath); err != nil {
		return "", fmt.Errorf("failed to save processed image: %v", err)
	}

	return processedFileName, nil
}

func (p *Processor) addWatermark(src image.Image) image.Image {
	bounds := src.Bounds()
	
	// Create a small version of itself in the corner as a "watermark"
	watermark := imaging.Thumbnail(src, 100, 100, imaging.Lanczos)
	
	return imaging.Overlay(src, watermark, image.Pt(bounds.Max.X-110, bounds.Max.Y-110), 0.5)
}
