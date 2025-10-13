package main

import (
	"fmt"
	"image"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/starillume/ase"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("error: ase filepath not provided")
		os.Exit(1)
	}
	
	fd, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("error: could not open file with filepath %s\n: %s", os.Args[1], err)
		os.Exit(1)
	}
	defer fd.Close()

	asef, err := ase.DeserializeFile(fd)
	if err != nil {
		fmt.Println("error: could not deserialize ase file: ", err)
		os.Exit(1)
	}
	
	handleInterrupt()

	width, height := int(asef.Header.Width), int(asef.Header.Height)
	if len(asef.Frames) > 1 {
		renderAnimation(asef.Frames, width, height)
	} else {
		renderFrame(asef.Frames[0].Image)
	}
}

func handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func () {
		<-c
		fmt.Print("\033[?25h")
		fmt.Print("\033[0m")
		os.Exit(0)
	}()
}

func renderFrame(img image.Image) {
	bounds := img.Bounds()
	var lastColored bool
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		lastColored = false
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a < 0x8000 {
				if lastColored {
					fmt.Print("\x1b[0m")
					lastColored = false
				}
				fmt.Print(" ")
				continue
			}

			if !lastColored {
				lastColored = true
			}

			fmt.Printf("\x1b[48;2;%d;%d;%dm ", r>>8, g>>8, b>>8)
		}

		if lastColored {
			fmt.Print("\x1b[0m")
		}
		fmt.Print("\n")
	}
}

func renderAnimation(frames []*ase.Frame, width int, height int) {
	fmt.Print("\033[?25l")
	for {
		for _, frame := range frames {
			fmt.Print("\033[2J\033[H")
			renderFrame(frame.Image)
			time.Sleep(time.Millisecond * time.Duration(frame.Duration))
		}
	}
}
