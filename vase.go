package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

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

	width, height := int(asef.Header.Width), int(asef.Header.Height)
	renderImage(composeFrameImages(width, height, getFrameImages(asef.Frames[0], width, height)))
}

func getFrameImages(frame ase.Frame, width int, height int) []image.Image {
	imgs := make([]image.Image, 0)
	for _, chunk := range frame.Chunks {
		if chunk.GetType() == ase.CelChunkHex {
			celChunk := chunk.(*ase.ChunkCelImage)
			pixelsRGBA := celChunk.ChunkCelRawImageData.Pixels.(ase.PixelsRGBA)
			img := pixelsRGBA.ToImage(int(celChunk.X), int(celChunk.Y), int(celChunk.ChunkCelDimensionData.Width), int(celChunk.ChunkCelDimensionData.Height), width, height)
			imgs = append(imgs, img)
		}
	}

	return imgs
}

func composeFrameImages(baseWidth int, baseHeight int, imgs []image.Image) image.Image {
	final := image.NewRGBA(image.Rect(0, 0, baseWidth, baseHeight))
	for _, img := range imgs {
		draw.Draw(final, img.Bounds(), img, image.Point{}, draw.Over)
	}

	return final
}

func renderImage(img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a < 0x8000 {
				fmt.Print(" ")
				continue
			}
			fmt.Printf("\x1b[48;2;%d;%d;%dm ", r>>8, g>>8, b>>8)
		}
		fmt.Print("\x1b[0m\n")
	}
}
