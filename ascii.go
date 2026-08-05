package main

import (
	"fmt"
	"image"

	ansiart "github.com/Zebbeni/ansiart"
)

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
