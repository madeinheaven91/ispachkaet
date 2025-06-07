package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/madeinheaven91/ispachkaet/internal"
)

func main() {
	bgColor := flag.String("bg", "ffff00", "RGB background color in hex.")
	title := flag.String("title", "испачкает", "Album title.")
	outputName := flag.String("output", "output", "Output file name.")
	flag.Parse()

	height := 1024
	padding := 32

	img := internal.InitRGBA(*internal.HexStringToRGBA(*bgColor))
	fg := image.Black
	fgLight := image.NewUniform( *internal.HexStringToRGBA("333333"))

	internal.AddPmrc(img, padding)

	c := internal.InitContext(img, fg)
	c.SetSrc(fgLight)
	internal.AddLabel(c, 2*padding, height-2*padding, "испачкает")
	c.SetFontSize(24)
	c.SetSrc(fg)
	internal.AddLabel(c, 2*padding, 4*padding, *title)

	f, _ := os.Create(fmt.Sprintf("output/%s.jpg", *outputName))
	err := jpeg.Encode(f, img, nil)
	if err != nil {
		panic(err)
	}
}
