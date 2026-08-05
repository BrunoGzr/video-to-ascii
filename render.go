package main

import (
	"fmt"
	"image/color"
	"os"
	"regexp"
	"strconv"
	"strings"

	gg "github.com/fogleman/gg"
)

var ansiFgRe = regexp.MustCompile(`\x1b\[38;2;(\d+);(\d+);(\d+)m`)

func bgColorFor(name string) color.Color {
	switch name {
	case "black":
		return color.RGBA{R: 15, G: 15, B: 15, A: 255}
	case "gray":
		return color.RGBA{R: 45, G: 45, B: 45, A: 255}
	case "white":
		return color.RGBA{R: 240, G: 240, B: 240, A: 255}
	default:
		return color.Black
	}
}

func fgColorFor(name string) color.Color {
	switch name {
	case "black":
		return color.RGBA{R: 240, G: 240, B: 240, A: 255}
	case "gray":
		return color.RGBA{R: 235, G: 235, B: 235, A: 255}
	case "white":
		return color.RGBA{R: 15, G: 15, B: 15, A: 255}
	default:
		return color.White
	}
}

func loadMonoFontFace(dc *gg.Context, fontSize float64) error {
	fontPaths := []string{
		"C:\\Windows\\Fonts\\cour.ttf",
		"C:\\Windows\\Fonts\\consola.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf",
	}
	for _, fontPath := range fontPaths {
		if _, err := os.Stat(fontPath); err == nil {
			if err := dc.LoadFontFace(fontPath, fontSize); err == nil {
				return nil
			}
		}
	}
	return fmt.Errorf("no monospace font found in known paths: %v", fontPaths)
}


func boostColor(r, g, b float64) (float64, float64, float64) {
	lum := 0.299*r + 0.587*g + 0.114*b
	if lum >= 0.35 {
		return r, g, b
	}
	if lum < 0.01 {
		return 0.35, 0.35, 0.35
	}
	scale := 0.35 / lum
	r = clamp01(r * scale)
	g = clamp01(g * scale)
	b = clamp01(b * scale)
	return r, g, b
}

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}


func renderAsciiToPNG(ascii string, width, height int, fontSize float64, outputPath string, bgColor string) error {
	dc := gg.NewContext(width, height)
	dc.SetColor(bgColorFor(bgColor))
	dc.Clear()

	if err := loadMonoFontFace(dc, fontSize); err != nil {
		return fmt.Errorf("failed to load font face: %w", err)
	}

	lines := strings.Split(ascii, "\n")
	lineHeight := fontSize * 1.5

	marginX := 10.0
	marginY := lineHeight
	for _, line := range lines {
		renderAsciiLine(dc, line, marginX, marginY, bgColor)
		marginY += lineHeight
	}
	return dc.SavePNG(outputPath)
}

func renderAsciiLine(dc *gg.Context, line string, x, y float64, bgColor string) {
	rem := line
	cx := x
	curR, curG, curB := 0.0, 0.0, 0.0

	for len(rem) > 0 {
		for strings.HasPrefix(rem, "\x1b[") {
			endM := strings.IndexByte(rem, 'm')
			if endM == -1 {
				rem = ""
				break
			}
			seq := rem[:endM+1]
			rem = rem[endM+1:]

			if m := ansiFgRe.FindStringSubmatch(seq); m != nil {
				r, _ := strconv.Atoi(m[1])
				g, _ := strconv.Atoi(m[2])
				b, _ := strconv.Atoi(m[3])
				curR = float64(r) / 255
				curG = float64(g) / 255
				curB = float64(b) / 255
			}
		}

		if len(rem) == 0 {
			break
		}

		nextEsc := strings.IndexByte(rem, '\x1b')
		var chunk string
		if nextEsc == -1 {
			chunk = rem
			rem = ""
		} else {
			chunk = rem[:nextEsc]
			rem = rem[nextEsc:]
		}

		if len(chunk) == 0 {
			continue
		}

		if strings.Contains(line, "\x1b[") {
			boostR, boostG, boostB := boostColor(curR, curG, curB)
			dc.SetRGB(boostR, boostG, boostB)
		} else {
			dc.SetColor(fgColorFor(bgColor))
		}

		dc.DrawStringAnchored(chunk, cx, y, 0, 1)
		w, _ := dc.MeasureString(chunk)
		cx += w
	}
}


func ansiStrip(s string) string {
	var result []rune
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

func measureAscii(ascii string, fontSize float64) (width, height int) {
	lines := strings.Split(ascii, "\n")
	lineHeight := fontSize * 1.5

	maxChars := 0
	for _, line := range lines {
		clean := ansiStrip(line)
		if len(clean) > maxChars {
			maxChars = len(clean)
		}
	}

	charWidth := fontSize * 0.6
	width = int(float64(maxChars)*charWidth + 20)
	height = int(float64(len(lines))*lineHeight + lineHeight)
	return
}
