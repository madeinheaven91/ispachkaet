package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/golang/freetype"
)

func initRGBA(color color.RGBA) *image.RGBA {
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

func initContext(img *image.RGBA, fg *image.Uniform) *freetype.Context {
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

func main() {
	width := 1024
	height := 1024

	img := initRGBA(color.RGBA{0xff, 0xff, 0x00, 0xff})
	fg := image.Black
	fgLight := image.NewUniform(color.RGBA{0x33, 0x33, 0x33, 0xff})

	pmrc, err := getImageFromFilePath("assets/pmrc.jpg")
	if err != nil {
		panic(err)
	}
	pmrcW := pmrc.Bounds().Max.X
	pmrcH := pmrc.Bounds().Max.Y
	padding := 32

	for x := range pmrcW {
		for y := range pmrcH {
			r, g, b, a := pmrc.At(x, y).RGBA()
			col := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			img.Set(height-pmrcW-padding+x, width-pmrcH-padding+y, col)
		}
	}

	c := initContext(img, fg)
	c.SetSrc(fgLight)
	addLabel(c, 2 * padding, height - 2 * padding, "испачкает")
	c.SetFontSize(24)
	c.SetSrc(fg)
	addLabel(c, 2 * padding, 4 * padding, "испачкает")

	f, _ := os.Create("output/image.jpg")
	err = jpeg.Encode(f, img, nil)
	if err != nil {
		panic(err)
	}
}

func addLabel(c *freetype.Context, x, y int, label string) {
	_, err := c.DrawString(label, freetype.Pt(x, y))
	if err != nil {
		fmt.Println(err)
	}
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, err := jpeg.Decode(f)
	return image, err
}
