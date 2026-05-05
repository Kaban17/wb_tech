package processor

import (
	"context"
	"fmt"
	"image"
	"math/rand"
	"path/filepath"
	"time"
	"wb_tech/l3_4/internal/model"

	"github.com/disintegration/imaging"
)

type Processor struct {
	basePath string
	rng      *rand.Rand
}

func NewProcessor(basePath string) *Processor {
	return &Processor{
		basePath: basePath,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (p *Processor) ProcessImage(ctx context.Context, img *model.Image) (string, error) {
	// Simulate 30% failure chance
	if p.rng.Float32() < 0.3 {
		return "", fmt.Errorf("random processing error (30%% chance triggered)")
	}

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
