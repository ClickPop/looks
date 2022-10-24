package generator

import (
	"image"
	"image/draw"
	"image/png"
)

type ImageWithPosition struct {
	Image    *image.Image
	Position int
}

func getImages(files <-chan *FileWithPosition, images chan<- *ImageWithPosition, errChan chan<- error) {
	for {
		select {
		case file, open := <-files:
			if !open {
				close(images)
				return
			}
			img, err := png.Decode(file.File)
			if err != nil {
				errChan <- err
			}
      image := ImageWithPosition{Image: &img, Position: file.Position}
      println(image.Image)
      images <- &image
		}
	}
}

func buildImage(images []*image.Image) *image.RGBA {
	baseImage := *images[0]
	rect := image.Rectangle{baseImage.Bounds().Min, baseImage.Bounds().Max}
	img := image.NewRGBA(rect)
	origin := image.Point{0, 0}
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
