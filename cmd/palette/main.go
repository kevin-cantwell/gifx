package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/AdRoll/goamz/s3"
	"github.com/codegangsta/cli"

	"github.com/kevin-cantwell/palette/internal/conf"
	"github.com/kevin-cantwell/palette/internal/stitcher"
)

func main() {
	app := cli.NewApp()
	app.Name = "palette"
	app.Usage = "A command-line tool for exposing the color palette of an image."
	app.Action = func(c *cli.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(os.Stdin, &buf)
		img, format, err := image.Decode(tee)
		if err != nil {
			exit(err)
		}
		fmt.Printf("Detected %v\n", format)

		var images []ImageData

		switch img.(type) {
		case *image.Paletted:
			gifImg, err := gif.DecodeAll(&buf)
			if err != nil {
				exit(err)
			}
			for _, frame := range gifImg.Image {
				f, err := ioutil.TempFile("", "palette")
				if err != nil {
					exit(err)
				}
				imageData := ImageData{Image: f.Name()}
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
		}

		f, err := ioutil.TempFile("", "palette")
		if err != nil {
			exit(err)
		}
		t, err := template.ParseFiles("palette.html")
		if err != nil {
			exit(err)
		}
		t.Execute(f, images)
		f.Close()
		if err := exec.Command("open", f.Name()).Run(); err != nil {
			exit(err)
		}
	}
	// app.Flags = []cli.Flag{
	// 	cli.BoolFlag{
	// 		Name:  "d, decompress",
	// 		Usage: "Decompresses the input instead of compressing the output.",
	// 	},
	// }
	app.Run(os.Args)
}

type ImageData struct {
	Image   string
	Palette []Color
}

type Color struct {
	Hex  string
	RGBA []uint32 // Should be 8-bit rgba values despite 32 bit integer
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func StitchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("image")
	if err != nil {
		StitchResponder{Err: err.Error()}.Respond(w, http.StatusBadRequest)
		return
	}
	defer file.Close()

	staticImg, _, err := image.Decode(file)
	if err != nil {
		log.Println("ERROR decoding static image:", err)
		StitchResponder{Err: err.Error()}.Respond(w, http.StatusInternalServerError)
		return
	}

	gifName := r.FormValue("gif")
	fmt.Println("DEBUG fetching gif", gifName)
	gifReader, err := conf.TimehopUploadsS3Bucket.GetReader("/gifs/" + gifName + ".gif")
	if err != nil {
		log.Println("ERROR fetching gif from s3:", err)
		StitchResponder{Err: err.Error()}.Respond(w, http.StatusInternalServerError)
		return
	}

	gifImg, err := gif.DecodeAll(gifReader)
	if err != nil {
		log.Println("ERROR decoding gif image:", err)
		StitchResponder{Err: err.Error()}.Respond(w, http.StatusInternalServerError)
		return
	}

	stitchr := stitcher.GIF{StaticImage: staticImg, GIFImage: gifImg}
	outGif, err := stitchr.Stitch()

	var buf bytes.Buffer
	gif.EncodeAll(&buf, outGif)

	userID := r.FormValue("user_id")
	if userID == "" {
		userID = "000"
	}
	uploadPath := fmt.Sprintf("/gifreactions/%s/%d.gif", userID, time.Now().Unix())
	if err := conf.TimehopUploadsS3Bucket.PutReader(uploadPath, &buf, int64(buf.Len()), "image/gif", s3.PublicRead, s3.Options{}); err != nil {
		StitchResponder{Err: err.Error()}.Respond(w, http.StatusInternalServerError)
		return
	}

	StitchResponder{
		GIF: "http://timehop.uploads.s3.amazonaws.com" + uploadPath,
	}.Respond(w, http.StatusOK)
}

type StitchResponder struct {
	GIF string `json:"gif,omitempty"`
	MP4 string `json:"mp4,omitempty"`
	Err string `json:"error,omitempty"`
}

func (resp StitchResponder) Respond(w http.ResponseWriter, status int) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	body, _ := json.Marshal(resp)
	return w.Write(body)
}

// var (
// 	timehopPalette = []color.Color{
// 		color.RGBA{0, 0, 0, 100},       // White
// 		color.RGBA{255, 255, 255, 100}, // Black
// 		color.RGBA{252, 95, 52, 100},   // Melon Red
// 		color.RGBA{131, 215, 218, 100}, // Teal
// 		color.RGBA{20, 182, 239, 100},  // Sky Blue
// 		color.RGBA{219, 111, 163, 100}, // Eggplant Purple
// 	}
// )
