package main

import (
	"fmt"
	"html/template"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/kevin-cantwell/gifx"
)

func main() {
	app := cli.NewApp()
	app.Name = "gifx"
	app.Usage = "A command-line tool for exposing the color palette of an image."
	app.Action = func(c *cli.Context) {
		var images []ImageData

		gifImg, err := gif.DecodeAll(os.Stdin)
		if err != nil {
			exit(err)
		}

		f, err := ioutil.TempFile("", "palette")
		if err != nil {
			exit(err)
		}
		masterImageName := f.Name()
		if err := gif.EncodeAll(f, gifImg); err != nil {
			exit(err)
		}
		f.Close()

		var origPalette []Color
		if palette, ok := gifImg.Config.ColorModel.(color.Palette); ok {
			for _, clr := range palette {
				r, g, b, a := clr.RGBA()
				r, g, b, a = r/0x101, g/0x101, b/0x101, a/0x101 // Convert 16bit to 8bit
				origPalette = append(origPalette, Color{
					RGBA: []uint32{r, g, b, a},
				})
			}
		}

		for _, frame := range gifImg.Image {
			f, err := ioutil.TempFile("", "palette")
			if err != nil {
				exit(err)
			}
			imageData := ImageData{
				Orig:        gifImg,
				OrigPalette: origPalette,
				Image:       masterImageName,
				Frame:       f.Name(),
				X:           frame.Bounds().Min.X,
				Y:           frame.Bounds().Min.Y,
			}
			if err := gif.Encode(f, frame, nil); err != nil {
				exit(err)
			}
			f.Close()

			for _, clr := range frame.Palette {
				r, g, b, a := clr.RGBA()
				r, g, b, a = r/0x101, g/0x101, b/0x101, a/0x101 // Convert 16bit to 8bit
				imageData.Palette = append(imageData.Palette, Color{
					RGBA: []uint32{r, g, b, a},
				})
			}
			images = append(images, imageData)
		}

		f, err = ioutil.TempFile("", "palette")
		if err != nil {
			exit(err)
		}

		t, err := template.New("palette").Parse(gifx.PaletteHTML)
		if err != nil {
			exit(err)
		}
		t.Execute(f, images)
		f.Close()
		if err := exec.Command("open", f.Name()).Run(); err != nil {
			exit(err)
		}
	}
	app.Run(os.Args)
}

type ImageData struct {
	Image       string
	Frame       string
	X           int
	Y           int
	Palette     []Color
	Orig        *gif.GIF
	OrigPalette []Color
}

type Color struct {
	Hex  string
	RGBA []uint32 // Should be 8-bit rgba values despite 32 bit integer
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
