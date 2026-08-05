package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	s := settings{fps: 0, bgColor: "black"}
	var asciiFrames []string
	var asciiSingle string

	reader := bufio.NewReader(os.Stdin)

	printHeader("Video to ASCII")
	for {
		fmt.Println("Main Menu")
		printMenu([]string{
			"1  - Select Media (Video/Photo)",
			"2  - Select Resolution",
			"3  - Convert to ASCII",
			"4  - Preview ASCII",
			"5  - Export (MP4/PNG)",
			"6  - Exit",
		})
		printStatus(&s, asciiFrames, asciiSingle)
		fmt.Println()
		fmt.Print("> ")

		option, err := readInt(reader)
		if err != nil {
			fmt.Println("Invalid option!")
			reader.ReadString('\n')
			continue
		}

		switch option {
		case 1:
			handleSelectMedia(reader, &s)
		case 2:
			handleSelectResolution(reader, &s)
		case 3:
			handleConvert(reader, &s, &asciiFrames, &asciiSingle)
		case 4:
			handlePreview(&s, asciiFrames, asciiSingle)
		case 5:
			handleExport(&s, asciiFrames, asciiSingle)
		case 6:
			fmt.Println("\nGoodbye!")
			return
		default:
			fmt.Println("Invalid option!")
		}
	}
}


func readInt(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return 0, err
	}
	line = strings.TrimSpace(line)
	var n int
	_, err = fmt.Sscanf(line, "%d", &n)
	return n, err
}


func readLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return ""
	}
	return strings.TrimSpace(line)
}

func handleSelectMedia(reader *bufio.Reader, s *settings) {
	fmt.Println("\n--- Select Media Type ---")
	printMenu([]string{
		"1  - Video (multiple frames)",
		"2  - Photo (single frame)",
	})
	fmt.Print("> ")
	mediaType, err := readInt(reader)
	if err != nil {
		fmt.Println("Invalid choice!")
		return
	}

	switch mediaType {
	case 1:
		s.media = "Video"
	case 2:
		s.media = "Photo"
	default:
		fmt.Println("Invalid choice!")
		return
	}
	fmt.Printf(">> Media set to: %s\n", s.media)

	if s.media == "Video" {
		fmt.Println("\n--- Select FPS ---")
		printMenu([]string{
			"1  - 30 FPS (More Smoothness)",
			"2  - 20 FPS (Balanced)",
			"3  - 10 FPS (For Testing)",
		})
		fmt.Print("> ")
		fpsOpt, err := readInt(reader)
		if err != nil || fpsOpt < 1 || fpsOpt > 3 {
			s.fps = 30
		} else {
			s.fps = fpsPresets[fpsOpt]
		}
		fmt.Printf(">> FPS set to: %d\n", s.fps)
	}

	fmt.Printf("\nEnter the full PATH to the %s:\n", s.media)
	fmt.Print("> ")
	s.url = readLine(reader)
	if s.url == "" {
		fmt.Println("Invalid path!")
		return
	}
	fmt.Println(">> Path saved!")

	fmt.Println("\n--- Color Mode ---")
	printMenu([]string{
		"1  - Colorful (preserves colors)",
		"2  - Black & White",
	})
	fmt.Print("> ")
	colorOpt, err := readInt(reader)
	if err != nil {
		s.useColor = false
	} else {
		s.useColor = colorOpt == 1
	}
	if s.useColor {
		fmt.Println(">> Color: Colorful")
	} else {
		fmt.Println(">> Color: Black & White")
	}

	fmt.Println("\n--- Background Color ---")
	printMenu([]string{
		"1  - Black",
		"2  - Gray (dark)",
		"3  - White",
	})
	fmt.Print("> ")
	bgOpt, err := readInt(reader)
	if err != nil {
		s.bgColor = "black"
	} else if bg, ok := bgPresets[bgOpt]; ok {
		s.bgColor = bg
	} else {
		s.bgColor = "black"
	}
	fmt.Printf(">> Background: %s\n", s.bgColor)
}

func handleSelectResolution(reader *bufio.Reader, s *settings) {
	fmt.Println("\n--- Select Resolution ---")
	options := []string{
		"1  - Native (match source file)",
		"2  - 720p  (High Quality)",
		"3  - 480p  (Balanced)",
		"4  - 360p  (Faster Render)",
	}
	printMenu(options)

	fmt.Print("> ")
	choice, err := readInt(reader)
	if err != nil || choice < 1 || choice > 4 {
		fmt.Println("Invalid resolution!")
		return
	}

	if choice == 1 {
		// Native resolution: no ffmpeg scale filter, grid computed at runtime
		s.resolution = ""
		s.resLabel = "Native (auto-detected from source)"
		s.asciiCols = 0
		s.asciiRows = 0
		s.nativeRes = true
	} else {
		preset := resolutionPresets[choice-2] // offset by 1 because option 1 is Native
		s.resolution = preset.filter
		s.resLabel = preset.label
		s.asciiCols = preset.cols
		s.asciiRows = preset.rows
		s.nativeRes = false
	}
	fmt.Printf(">> Resolution set to: %s\n", s.resLabel)
}

func handleConvert(reader *bufio.Reader, s *settings, asciiFrames *[]string, asciiSingle *string) {
	if s.url == "" || (!s.nativeRes && s.resolution == "") {
		fmt.Println("!! Please select media and resolution first.")
		return
	}
	if s.fps == 0 && s.media == "Video" {
		s.fps = 30
	}

	fmt.Println("\nStarting conversion...")
	if s.media == "Photo" {
		img, err := extractSingleFrame(s.url, s.resolution)
		if err != nil {
			fmt.Printf("!! Error: %v\n", err)
			return
		}
		cols, rows := s.asciiCols, s.asciiRows
		if s.nativeRes {
			cols, rows = computeNativeGrid(img)
			fmt.Printf(">> Source dimensions: %dx%d → ASCII grid: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy(), cols, rows)
		}
		result, err := convertSingleToAscii(img, cols, rows, s.useColor)
		if err != nil {
			fmt.Printf("!! Error: %v\n", err)
			return
		}
		*asciiSingle = result
		*asciiFrames = nil
		fmt.Println(">> Photo conversion successful!")
	} else {
		frames, err := extractFrames(s.url, s.resolution, s.fps)
		if err != nil {
			fmt.Printf("!! Error: %v\n", err)
			return
		}
		cols, rows := s.asciiCols, s.asciiRows
		if s.nativeRes && len(frames) > 0 {
			cols, rows = computeNativeGrid(frames[0])
			fmt.Printf(">> Source dimensions: %dx%d → ASCII grid: %dx%d\n", frames[0].Bounds().Dx(), frames[0].Bounds().Dy(), cols, rows)
		}
		*asciiFrames, err = convertToAscii(frames, cols, rows, s.useColor)
		if err != nil {
			fmt.Printf("!! Error: %v\n", err)
			return
		}
		*asciiSingle = ""
		fmt.Printf(">> Conversion successful! %d frames generated.\n", len(*asciiFrames))
	}
}

func handlePreview(s *settings, asciiFrames []string, asciiSingle string) {
	if s.media == "Photo" {
		if asciiSingle == "" {
			fmt.Println("!! Please convert to ASCII first (option 3).")
			return
		}
		fmt.Println("\n--- ASCII Preview ---")
		fmt.Print("\033[H\033[J")
		fmt.Print(asciiSingle)
		fmt.Print("\033[0m")
		fmt.Println()
	} else {
		if len(asciiFrames) == 0 {
			fmt.Println("!! Please convert to ASCII first (option 3).")
			return
		}
		fps := s.fps
		if fps == 0 {
			fps = 30
		}
		fmt.Println("\n--- ASCII Preview ---")
		previewAscii(asciiFrames, fps)
	}
}

func handleExport(s *settings, asciiFrames []string, asciiSingle string) {
	if s.media == "Photo" {
		if asciiSingle == "" {
			fmt.Println("!! Please convert to ASCII first (option 3).")
			return
		}
		fmt.Println("\nExporting to PNG...")
		if err := exportToPNG(asciiSingle, s.url, s.bgColor); err != nil {
			fmt.Printf("!! Error: %v\n", err)
			return
		}
		fmt.Println(">> Export successful!")
	} else {
		if len(asciiFrames) == 0 {
			fmt.Println("!! Please convert to ASCII first (option 3).")
			return
		}
		fps := s.fps
		if fps == 0 {
			fps = 30
		}
		fmt.Println("\nExporting to MP4...")
		if err := exportToMp4(asciiFrames, fps, s.url, s.bgColor); err != nil {
			fmt.Printf("!! Error: %v\n", err)
			return
		}
		fmt.Println(">> Export successful!")
	}
}
