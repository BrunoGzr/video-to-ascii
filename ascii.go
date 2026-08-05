package main

import (
	"fmt"
	"image"

	ansiart "github.com/Zebbeni/ansiart"
)

// computeNativeGrid calculates ASCII grid dimensions from an image's aspect
// ratio, preserving the source proportions in the rendered output.
//
// The PNG renderer draws characters at width=fontSize*0.6 and height=fontSize*1.5,
// giving each cell an aspect ratio of 0.4 (w/h). To keep the output image
// proportional to the source, we compute:
//
//	rows = cols * (imgH / imgW) * 0.4
//
// We target ~200 columns for landscape. For portrait sources, we increase
// column count so rows stay above ~100 — otherwise the long axis would
// have too few characters to produce recognizable detail.
func computeNativeGrid(img image.Image) (cols, rows int) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w == 0 || h == 0 {
		return 200, 100
	}

	aspect := float64(w) / float64(h) // >1 = landscape, <1 = portrait

	if aspect >= 1 {
		// Landscape: 200 columns is a good baseline
		cols = 200
	} else {
		// Portrait: scale cols up so the tall axis has ~100+ rows
		// rows = cols / aspect * 0.4 → solve for rows = 100:
		// cols = 100 * aspect / 0.4 = 250 * aspect
		cols = int(250.0 * aspect)
		if cols < 100 {
			cols = 100
		}
	}

	rows = int(float64(cols) / aspect * 0.4)
	if rows < 1 {
		rows = 1
	}
	return
}


func convertToAscii(imgs []image.Image, asciiWidth, asciiHeight int, useColor bool) ([]string, error) {
	asciiFrames := make([]string, len(imgs))
	configs := ansiConfig(asciiWidth, asciiHeight)

	for i, frame := range imgs {
		result, err := ansiart.Render(frame, configs)
		if err != nil {
			return nil, fmt.Errorf("failed to render ascii for frame %d: %w", i, err)
		}
		if !useColor {
			asciiFrames[i] = ansiStrip(result)
		} else {
			asciiFrames[i] = result
		}
	}
	return asciiFrames, nil
}

func convertSingleToAscii(img image.Image, asciiWidth, asciiHeight int, useColor bool) (string, error) {
	configs := ansiConfig(asciiWidth, asciiHeight)

	result, err := ansiart.Render(img, configs)
	if err != nil {
		return "", fmt.Errorf("failed to render ascii: %w", err)
	}
	if !useColor {
		return ansiStrip(result), nil
	}
	return result, nil
}

func ansiConfig(width, height int) ansiart.Options {
	configs := ansiart.DefaultOptions()
	configs.Width = width
	configs.Height = height
	configs.CharacterMode = ansiart.Ascii
	configs.ColorBg = false
	configs.TrueColor = true
	return configs
}
