package internal

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/golang/freetype"
)

func HexStringToRGBA(input string) *color.RGBA {
	if len(input) != 6 {
		log.Fatalf("Couldn't parse hex color %s.\n", input)
		os.Exit(1)
	}

	r, err := strconv.ParseUint(input[0:2], 16, 16)
	if err != nil {
		log.Fatalf("Couldn't parse hex color %s.\n%q", input, err)
		os.Exit(1)
	}
	g, err := strconv.ParseUint(input[2:4], 16, 16)
	if err != nil {
		log.Fatalf("Couldn't parse hex color %s.\n%q", input, err)
		os.Exit(1)
	}
	b, err := strconv.ParseUint(input[4:6], 16, 16)
	if err != nil {
		log.Fatalf("Couldn't parse hex color %s.\n%q", input, err)
		os.Exit(1)
	}
	res := color.RGBA{
		uint8(r), uint8(g), uint8(b), 0xff,
	}
	return &res
}

func InitRGBA(color color.RGBA) *image.RGBA {
	width := 1024
	height := 1024
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := range width {
		for y := range height {
			img.Set(x, y, color)
		}
	}
	return img
}

func InitContext(img *image.RGBA, fg *image.Uniform) *freetype.Context {
	bytes, err := os.ReadFile("assets/times.ttf")
	if err != nil {
		log.Fatalf("Couldn't read times.ttf.\n%q", err)
		os.Exit(1)
	}
	font, err := freetype.ParseFont(bytes)
	if err != nil {
		log.Fatalf("Couldn't parse times.ttf.\n%q", err)
		os.Exit(1)
	}
	c := freetype.NewContext()
	c.SetDPI(240)
	c.SetFont(font)
	c.SetFontSize(16)
	c.SetDst(img)
	c.SetSrc(fg)
	c.SetClip(img.Bounds())
	return c
}

func AddLabel(c *freetype.Context, x, y int, label string) {
	substrings := strings.SplitSeq(label, `\n`)
	for s := range substrings {
		_, err := c.DrawString(s, freetype.Pt(x, y))
		if err != nil {
			log.Printf("Couldn't draw label %s at (%d, %d). %q", label, x, y, err)
		}
		y += int(math.Floor((float64(c.PointToFixed(16)>>6) * 1.5)))
	}
}

func AddLabelFromEnd(c *freetype.Context, x, y int, label string) {
	substrings := strings.SplitSeq(label, `\n`)
	for s := range substrings {
		// FIXME: костыль
		advanced, err := c.DrawString(s, freetype.Pt(0, -100))
		if err != nil {
			log.Printf("Couldn't draw label %s at (%d, %d). %q", s, x, y, err)
		}
		_, err = c.DrawString(s, freetype.Pt(x-advanced.X.Round(), y))
		if err != nil {
			log.Printf("Couldn't draw label %s at (%d, %d). %q", label, x, y, err)
		}
		y += int(math.Floor((float64(c.PointToFixed(16)>>6) * 1.5)))
	}
}

func AddLabelCentered(c *freetype.Context, x, y int, label string) {
	substrings := strings.SplitSeq(label, `\n`)
	for s := range substrings {
		// FIXME: костыль
		advanced, err := c.DrawString(s, freetype.Pt(0, -100))
		if err != nil {
			log.Printf("Couldn't draw label %s at (%d, %d). %q", s, x, y, err)
		}
		_, err = c.DrawString(s, freetype.Pt(x-advanced.X.Round()/2, y))
		if err != nil {
			log.Printf("Couldn't draw label %s at (%d, %d). %q", label, x, y, err)
		}
		y += int(math.Floor((float64(c.PointToFixed(16)>>6) * 1.5)))
	}
}

func GetImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, err := jpeg.Decode(f)
	return image, err
}

func AddPmrc(img *image.RGBA, padding int) {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	pmrc, err := GetImageFromFilePath("assets/pmrc.jpg")
	if err != nil {
		log.Printf("Couldn't open image %s. %q", "assets/pmrc.jpg", err)
	}
	pmrcW := pmrc.Bounds().Max.X
	pmrcH := pmrc.Bounds().Max.Y

	for x := range pmrcW {
		for y := range pmrcH {
			r, g, b, a := pmrc.At(x, y).RGBA()
			col := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			img.Set(height-pmrcW-padding+x, width-pmrcH-padding+y, col)
		}
	}
}
