package main

import (
	"fmt"
	"strings"
	"time"
)


type settings struct {
	media       string
	url         string
	resolution  string
	resLabel    string // display name for the current resolution
	fps         int
	useColor    bool
	bgColor     string
	asciiCols   int
	asciiRows   int
	nativeRes   bool // when true, grid is computed from the decoded frame's aspect ratio
}


type resolutionPreset struct {
	label  string
	filter string
	cols   int
	rows   int
}

var resolutionPresets = []resolutionPreset{
	{"720p  (High Quality)", "scale=1280:720", 200, 100},
	{"480p  (Balanced)", "scale=854:480", 160, 80},
	{"360p  (Faster Render)", "scale=640:360", 120, 60},
}


var fpsPresets = map[int]int{1: 30, 2: 20, 3: 10}


var bgPresets = map[int]string{1: "black", 2: "gray", 3: "white"}


func printHeader(title string) {
	width := max(len(title)+8, 44)
	border := strings.Repeat("=", width)
	pad := (width - len(title) - 4) / 2
	if pad < 0 {
		pad = 0
	}
	fmt.Printf("\n%s\n", border)
	fmt.Printf("==%s%s  ==\n", strings.Repeat(" ", pad), strings.ToUpper(title))
	fmt.Printf("%s\n\n", border)
}


func printMenu(options []string) {
	width := 42
	border := "+" + strings.Repeat("-", width-2) + "+"
	fmt.Println(border)
	for _, opt := range options {
		padded := opt
		if len(padded) < width-4 {
			padded = padded + strings.Repeat(" ", width-4-len(padded))
		}
		fmt.Printf("| %s |\n", padded)
	}
	fmt.Println(border)
}


func printStatus(s *settings, asciiFrames []string, asciiSingle string) {
	if s.media != "" {
		fmt.Printf("  Media: %s\n", s.media)
	}
	if s.resLabel != "" {
		fmt.Printf("  Resolution: %s\n", s.resLabel)
	}
	if s.fps > 0 {
		fmt.Printf("  FPS: %d\n", s.fps)
	}
	if s.useColor {
		fmt.Println("  Color: Colorful")
	} else {
		fmt.Println("  Color: Black & White")
	}
	fmt.Printf("  Background: %s\n", s.bgColor)
	if s.url != "" {
		fmt.Printf("  File: %s\n", s.url)
	}
	if len(asciiFrames) > 0 {
		fmt.Printf("  Frames: %d converted\n", len(asciiFrames))
	}
	if asciiSingle != "" {
		fmt.Println("  Photo: converted")
	}
}


func previewAscii(frames []string, fps int) {
	if len(frames) == 0 {
		fmt.Println("No Ascii to preview")
		return
	}

	delay := time.Second / time.Duration(fps)
	for _, frame := range frames {
		fmt.Print("\033[H\033[J")
		fmt.Print(frame)
		fmt.Print("\033[0m")
		time.Sleep(delay)
	}
	fmt.Print("\033[0m")
}


func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
