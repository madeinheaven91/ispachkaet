package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/madeinheaven91/ispachkaet/internal"
)

var (
	bgColor    = flag.String("bg", "ffff00", "RGB background color in hex.")
	title      = flag.String("title", "испачкает", "Album title. To add newlines add '\\n'.")
	author     = flag.String("author", "испачкает", "Author name.")
	outputName = flag.String("output", "output", "Output file name.")
	alignment  = flag.String("alignment", "left", "Title text alignment. Possible values: left, center, right")
)

func main() {
	flag.Parse()

	size := 1024
	padding := 32
	fg := image.Black
	fgLight := image.NewUniform(*internal.HexStringToRGBA("333333"))

	img := internal.InitRGBA(*internal.HexStringToRGBA(*bgColor))
	// Add parental advisory
	internal.AddPmrc(img, padding)

	// Adding labels
	c := internal.InitContext(img, fg)
	c.SetSrc(fgLight)
	internal.AddLabel(c, 2*padding, size-2*padding, *author)
	c.SetFontSize(24)
	c.SetSrc(fg)
	switch *alignment {
	case "left":
		internal.AddLabel(c, 2*padding, 4*padding, *title)
	case "center":
		internal.AddLabelCentered(c, size/2, size/2-2*padding, *title)
	case "right":
		internal.AddLabelFromEnd(c, size-2*padding, 4*padding, *title)
	default:
		log.Fatalf("Unknown alignment argument: %s", *alignment)
	}

	// Save the file
	f, _ := os.Create(fmt.Sprintf("output/%s.png", *outputName))
	err := png.Encode(f, img)
	if err != nil {
		log.Fatalf("Couldn't save image to output/%s.png.\n%q", *outputName, err)
		os.Exit(1)
	}
}
