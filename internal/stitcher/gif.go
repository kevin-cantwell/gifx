package stitcher

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
)

var (
	timehopPalette = []color.Color{
		color.RGBA{0, 0, 0, 100},       // White
		color.RGBA{255, 255, 255, 100}, // Black
		color.RGBA{252, 95, 52, 100},   // Melon
		color.RGBA{131, 215, 218, 100}, // Teal
		color.RGBA{20, 182, 239, 100},  // Sky
		color.RGBA{219, 111, 163, 100}, // Eggplant
		color.RGBA{213, 219, 68, 100},  // Yoda
		color.RGBA{251, 207, 39, 100},  // RubberDuck
	}
)

type GIF struct {
	StaticImage image.Image
	GIFImage    *gif.GIF
}

func (g *GIF) Stitch() (*gif.GIF, error) {
	if g.GIFImage.Config.Width != g.StaticImage.Bounds().Dx() {
		return nil, errors.New("stitch: provided images must have the same width")
	}

	outGif := &gif.GIF{
		BackgroundIndex: g.GIFImage.BackgroundIndex,
		Delay:           g.GIFImage.Delay,
		Disposal:        g.GIFImage.Disposal,
		LoopCount:       g.GIFImage.LoopCount,
		// Must specify the new width/height in config
		Config: image.Config{
			// ColorModel: g.GIFImage.Config.ColorModel,
			Width:  g.GIFImage.Config.Width,
			Height: g.GIFImage.Config.Height + g.StaticImage.Bounds().Dy(),
		},
	}

	// outGif.Config.

	workers := make([]chan *image.Paletted, len(g.GIFImage.Image))
	for i := 0; i < len(workers); i++ {
		workers[i] = make(chan *image.Paletted)
	}

	for i, frame := range g.GIFImage.Image {
		go func(i int, frame *image.Paletted, worker chan<- *image.Paletted) {

			// // Add as many of the timehop palette colors as possible
			// remaining := 256 - len(frame.Palette)
			// for j := 0; j < remaining && j < len(timehopPalette); j++ {
			//  clr := g.StaticImage.ColorModel().Convert(timehopPalette[j])
			//  frame.Palette = append(frame.Palette, clr)
			// }

			if i == 1 {
				fmt.Println(i, "--------------------------")
				for _, color := range frame.Palette {
					fmt.Println(color)
				}
			}

			newPaletted := image.NewPaletted(image.Rect(0, 0, outGif.Config.Width, outGif.Config.Height), frame.Palette) // palette.Plan9)

			// Draw the static image in each frame
			draw.Draw(newPaletted, g.StaticImage.Bounds(), g.StaticImage, g.StaticImage.Bounds().Min, draw.Src)
			gifDrawRect := frame.Bounds().Add(image.Point{0, g.StaticImage.Bounds().Dy()})
			// Draw the gif frame
			draw.Draw(newPaletted, gifDrawRect, frame, frame.Bounds().Min, draw.Src)

			worker <- newPaletted
		}(i, frame, workers[i])
	}

	outGif.Image = make([]*image.Paletted, len(workers))
	for i := 0; i < len(workers); i++ {
		outGif.Image[i] = <-workers[i]
		fmt.Printf("DEBUG gif stitched %d of %d\n", i+1, len(workers))
	}

	return outGif, nil
}
