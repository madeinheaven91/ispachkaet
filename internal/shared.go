package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"strconv"

	"github.com/golang/freetype"
)

func HexStringToRGBA(input string) *color.RGBA {
	if len(input) != 6 {
		panic(fmt.Sprintf("expected input length 6, got %d", len(input)))
	}

	r, err := strconv.ParseUint(input[0:2], 16, 16)
	if err != nil {
		panic("couldn't parse hex color")
	}
	g, err := strconv.ParseUint(input[2:4], 16, 16)
	if err != nil {
		panic("couldn't parse hex color")
	}
	b, err := strconv.ParseUint(input[4:6], 16, 16)
	if err != nil {
		panic("couldn't parse hex color")
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
		panic(err)
	}
	font, err := freetype.ParseFont(bytes)
	if err != nil {
		panic(err)
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
	_, err := c.DrawString(label, freetype.Pt(x, y))
	if err != nil {
		fmt.Println(err)
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
		panic(err)
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
