package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)


func downloadsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, "Downloads"), nil
}



func buildOutputPath(inputPath, extension string) (string, error) {
	dl, err := downloadsDir()
	if err != nil {
		return "", err
	}

	base := filepath.Base(inputPath)
	// Strip the original extension
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	outputName := name + "_Ascii" + extension
	return filepath.Join(dl, outputName), nil
}



func exportToMp4(asciiFrames []string, fps int, inputPath string, bgColor string) error {
	if len(asciiFrames) == 0 {
		return fmt.Errorf("no ascii frames to export")
	}

	tempDir, err := os.MkdirTemp("", "mp4_frames_*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	fontSize := 28.0



	imgWidth, imgHeight := measureAscii(asciiFrames[0], fontSize)

	for i, ascii := range asciiFrames {
		framePath := filepath.Join(tempDir, fmt.Sprintf("frame_%09d.png", i+1))
		if err := renderAsciiToPNG(ascii, imgWidth, imgHeight, fontSize, framePath, bgColor); err != nil {
			return fmt.Errorf("failed to render frame %d: %w", i, err)
		}
	}

	pattern := filepath.Join(tempDir, "frame_%09d.png")

	outputFile, err := buildOutputPath(inputPath, ".mp4")
	if err != nil {
		return fmt.Errorf("failed to determine output path: %w", err)
	}

	err = ffmpeg.Input(pattern, ffmpeg.KwArgs{"r": fps}).
		Output(outputFile, ffmpeg.KwArgs{
			"c:v":     "libx264",
			"pix_fmt": "yuv420p",
		}).Run()
	if err != nil {
		return fmt.Errorf("failed to render MP4: %w", err)
	}
	fmt.Printf("Video saved as %s\n", outputFile)
	return nil
}

func exportToPNG(ascii string, inputPath string, bgColor string) error {
	fontSize := 28.0
	imgWidth, imgHeight := measureAscii(ascii, fontSize)

	outputFile, err := buildOutputPath(inputPath, ".png")
	if err != nil {
		return fmt.Errorf("failed to determine output path: %w", err)
	}

	if err := renderAsciiToPNG(ascii, imgWidth, imgHeight, fontSize, outputFile, bgColor); err != nil {
		return err
	}
	fmt.Printf("Image saved as %s\n", outputFile)
	return nil
}
