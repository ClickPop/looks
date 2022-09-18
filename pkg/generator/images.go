package generator

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"log"
)

func getImages(files []*bytes.Reader) ([]*image.Image, error) {
	var images []*image.Image
	for i := 0; i < len(files); i++ {
		img, err := png.Decode(files[i])
		if err != nil {
			return nil, err
		}
		images = append(images, &img)
	}
	return images, nil
}

func buildImage(images []*image.Image, i int) *image.RGBA {
	baseImage := *images[0]
	rect := image.Rectangle{baseImage.Bounds().Min, baseImage.Bounds().Max}
	img := image.NewRGBA(rect)
	origin := image.Point{0, 0}
	log.Printf("Layering assets for image #%d\n", i)
	for i := 0; i < len(images); i++ {
		currImg := *images[i]
		if i == 0 {
			draw.Draw(img, currImg.Bounds(), currImg, origin, draw.Src)
		} else {
			draw.Draw(img, currImg.Bounds(), currImg, origin, draw.Over)
		}
	}
	return img
}
