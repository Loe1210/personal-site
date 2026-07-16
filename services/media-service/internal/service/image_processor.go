package service

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/png"
)

const imageProcessorMaxEdge = 320

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) Process(srcPath string, dstPath string) (bool, error) {
	if p == nil {
		return false, fmt.Errorf("image processor is required")
	}
	if strings.TrimSpace(srcPath) == "" {
		return false, fmt.Errorf("source path is required")
	}
	if strings.TrimSpace(dstPath) == "" {
		return false, fmt.Errorf("destination path is required")
	}

	file, err := os.Open(srcPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return false, nil
	}

	thumb := resizeWithinBounds(src, imageProcessorMaxEdge)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return false, err
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = out.Close()
	}()

	if err := jpeg.Encode(out, thumb, &jpeg.Options{Quality: 85}); err != nil {
		_ = os.Remove(dstPath)
		return false, err
	}

	return true, nil
}

func resizeWithinBounds(src image.Image, maxEdge int) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return src
	}

	scale := 1.0
	if width > maxEdge || height > maxEdge {
		sw := float64(maxEdge) / float64(width)
		sh := float64(maxEdge) / float64(height)
		if sw < sh {
			scale = sw
		} else {
			scale = sh
		}
	}
	if scale >= 1 {
		return cloneImage(src)
	}

	newWidth := int(float64(width) * scale)
	newHeight := int(float64(height) * scale)
	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	drawNearest(dst, src)
	return dst
}

func cloneImage(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), src, bounds.Min, draw.Src)
	return dst
}

func drawNearest(dst *image.RGBA, src image.Image) {
	dstBounds := dst.Bounds()
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()
	if dstBounds.Dx() == 0 || dstBounds.Dy() == 0 {
		return
	}

	for y := 0; y < dstBounds.Dy(); y++ {
		srcY := srcBounds.Min.Y + (y*srcHeight)/dstBounds.Dy()
		if srcY >= srcBounds.Max.Y {
			srcY = srcBounds.Max.Y - 1
		}
		for x := 0; x < dstBounds.Dx(); x++ {
			srcX := srcBounds.Min.X + (x*srcWidth)/dstBounds.Dx()
			if srcX >= srcBounds.Max.X {
				srcX = srcBounds.Max.X - 1
			}
			dst.Set(x, y, color.NRGBAModel.Convert(src.At(srcX, srcY)))
		}
	}
}
