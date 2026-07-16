package service

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImageProcessorCreatesJPEGThumbnailForPNG(t *testing.T) {
	srcDir := t.TempDir()
	srcPath := filepath.Join(srcDir, "source.png")
	dstPath := filepath.Join(srcDir, "thumb.jpg")

	if err := writeTestPNG(srcPath, 640, 480); err != nil {
		t.Fatalf("write png: %v", err)
	}

	processed, err := NewImageProcessor().Process(srcPath, dstPath)
	if err != nil {
		t.Fatalf("process image: %v", err)
	}
	if !processed {
		t.Fatal("expected png input to be processed")
	}

	file, err := os.Open(dstPath)
	if err != nil {
		t.Fatalf("open thumbnail: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("decode thumbnail as jpeg: %v", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != 320 || bounds.Dy() != 240 {
		t.Fatalf("expected thumbnail to preserve aspect ratio at 320px max edge, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestImageProcessorReturnsUnprocessedForUnsupportedInput(t *testing.T) {
	srcDir := t.TempDir()
	srcPath := filepath.Join(srcDir, "source.svg")
	dstPath := filepath.Join(srcDir, "thumb.jpg")

	if err := os.WriteFile(srcPath, []byte("<svg xmlns='http://www.w3.org/2000/svg'></svg>"), 0o644); err != nil {
		t.Fatalf("write svg: %v", err)
	}

	processed, err := NewImageProcessor().Process(srcPath, dstPath)
	if err != nil {
		t.Fatalf("process unsupported image: %v", err)
	}
	if processed {
		t.Fatal("expected svg input to be left unprocessed")
	}
	if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
		t.Fatalf("expected thumbnail file to be absent, got err=%v", err)
	}
}

func TestImageProcessorCreatesJPEGThumbnailForGIF(t *testing.T) {
	srcDir := t.TempDir()
	srcPath := filepath.Join(srcDir, "source.gif")
	dstPath := filepath.Join(srcDir, "thumb.jpg")

	if err := writeTestGIF(srcPath, 200, 400); err != nil {
		t.Fatalf("write gif: %v", err)
	}

	processed, err := NewImageProcessor().Process(srcPath, dstPath)
	if err != nil {
		t.Fatalf("process image: %v", err)
	}
	if !processed {
		t.Fatal("expected gif input to be processed")
	}

	file, err := os.Open(dstPath)
	if err != nil {
		t.Fatalf("open thumbnail: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("decode thumbnail as jpeg: %v", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != 160 || bounds.Dy() != 320 {
		t.Fatalf("expected thumbnail to preserve aspect ratio at 320px max edge, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func writeTestPNG(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 200, G: 50, B: 25, A: 255})
		}
	}
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func writeTestGIF(path string, width, height int) error {
	img := image.NewPaletted(image.Rect(0, 0, width, height), []color.Color{
		color.Black,
		color.White,
	})
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetColorIndex(x, y, uint8((x+y)%2))
		}
	}
	buf := new(bytes.Buffer)
	if err := gif.Encode(buf, img, nil); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func TestImageProcessorLeavesSmallJPEGWithinBounds(t *testing.T) {
	srcDir := t.TempDir()
	srcPath := filepath.Join(srcDir, "source.jpg")
	dstPath := filepath.Join(srcDir, "thumb.jpg")

	if err := writeTestJPEG(srcPath, 120, 90); err != nil {
		t.Fatalf("write jpeg: %v", err)
	}

	processed, err := NewImageProcessor().Process(srcPath, dstPath)
	if err != nil {
		t.Fatalf("process jpeg: %v", err)
	}
	if !processed {
		t.Fatal("expected jpeg input to be processed")
	}

	file, err := os.Open(dstPath)
	if err != nil {
		t.Fatalf("open thumbnail: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("decode thumbnail as jpeg: %v", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != 120 || bounds.Dy() != 90 {
		t.Fatalf("expected small jpeg size to stay unchanged, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func writeTestJPEG(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 15, G: 120, B: 200, A: 255})
		}
	}
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 95}); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func TestImageProcessorRejectsEmptySourcePath(t *testing.T) {
	processed, err := NewImageProcessor().Process("", filepath.Join(t.TempDir(), "thumb.jpg"))
	if err == nil {
		t.Fatal("expected empty source path to fail")
	}
	if processed {
		t.Fatal("expected empty source path to remain unprocessed")
	}
	if !strings.Contains(err.Error(), "source") {
		t.Fatalf("expected source validation error, got %v", err)
	}
}
