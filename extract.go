package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)


func extractFrames(fileDirectory string, resolution string, fps int) ([]image.Image, error) {
	tempDir, err := os.MkdirTemp("", "frames_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	outPattern := filepath.Join(tempDir, "frame_%09d.jpg")

	kwargs := ffmpeg.KwArgs{"r": fps}
	if resolution != "" {
		kwargs["vf"] = resolution
	}

	err = ffmpeg.Input(fileDirectory).
		Output(outPattern, kwargs).Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract frames with ffmpeg: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(tempDir, "frame_*.jpg"))
	if err != nil {
		return nil, fmt.Errorf("failed to list extracted frames: %w", err)
	}

	var frames []image.Image
	for _, framePath := range files {
		file, err := os.Open(framePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open frame %s: %w", framePath, err)
		}
		img, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to decode frame %s: %w", framePath, err)
		}
		frames = append(frames, img)
	}
	return frames, nil
}


func extractSingleFrame(fileDirectory string, resolution string) (image.Image, error) {
	absPath, err := filepath.Abs(fileDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	tempDir, err := os.MkdirTemp("", "frame_single_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	outPath := filepath.Join(tempDir, "frame.jpg")

	kwargs := ffmpeg.KwArgs{"frames": 1, "update": 1}
	if resolution != "" {
		kwargs["vf"] = resolution
	}

	err = ffmpeg.Input(absPath).
		Output(outPath, kwargs).Run()
	if err != nil {
		f, openErr := os.Open(absPath)
		if openErr != nil {
			return nil, fmt.Errorf("failed to open image: %w", openErr)
		}
		defer f.Close()
		img, _, decErr := image.Decode(f)
		if decErr != nil {
			return nil, fmt.Errorf("failed to extract frame with ffmpeg: %v; failed to decode image: %w", err, decErr)
		}
		return img, nil
	}

	f, err := os.Open(outPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open extracted frame: %w", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode extracted frame: %w", err)
	}
	return img, nil
}
