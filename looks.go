package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	PieceOrder []string `json:"piece-order,omitempty"`
	Filename string `json:"filename,omitempty"`
	Rarity map[string]map[string]float64 `json:"rarity,omitempty"`
}

func main() {
	config, err := loadJSON();
	if err != nil {
		handleError(err, exit)
	}
	files, err := loadFiles(config)
	handleError(err)
	images, err := getImages(files)
	handleError(err)

	rect := image.Rectangle{images[1].Bounds().Min, images[1].Bounds().Max}
	img := image.NewRGBA(rect)
	origin := image.Point{0, 0}
	for i := 0; i < len(images); i++ {
		currImg := images[i]
		if i == 0 {
			draw.Draw(img, currImg.Bounds(), currImg, origin, draw.Src)
		} else {
			draw.Draw(img, currImg.Bounds(), currImg, origin, draw.Over)
		}
	}

	out, err := os.Create("test.png")
	if err != nil {
		fmt.Println(err)
	}
	
	png.Encode(out, img)
}

func loadFiles(config Config) ([]*os.File, error) {
	fileNames := config.PieceOrder
	var files []*os.File
	for i := 0; i < len(fileNames); i++ {
		file := fileNames[i]
		piece := handleRarity(config.Rarity[file])
		// We need to add the randomness logic here.
		if piece != "nil" {
			reader, err := os.Open(fmt.Sprintf("rat-parts_%s-%s.png", file, piece))
			if err != nil {
				return nil, handleError(err, exit)
			}
			files = append(files, reader)
		}
	}
	return files, nil
}

func handleRarity(rarities map[string]float64) string {
	var chances []string
	for k, v := range rarities {
		for i := 0; i < int(v); i++ {
			chances = append(chances, k)
		}
	}

	seed := time.Now().Unix() + int64(time.Now().Nanosecond())
	rand.Seed(seed)

	random := rand.Intn(len(chances))

	return chances[random]
}

func getImages(files []*os.File) ([]image.Image, error) {
	var images []image.Image
	for i := 0; i < len(files); i++ {
		img, err := png.Decode(files[i])
		if err != nil {
			return nil, handleError(err, exit)
		}
		images = append(images, img)
	}
	return images, nil
}

func handleError(err error, callbacks ...func(err error) error) error {
	if err != nil {
		fmt.Println(err)
		for i := 0; i < len(callbacks); i++ {
			callback := callbacks[i]
			val := callback(err)
			if (val != nil) {
				return val
			}
		}
	}
	return err
}

func exit(_ error) error {
	os.Exit(1)
	return nil
}

func loadJSON() (Config, error) {
	var config Config
	data, err := ioutil.ReadFile("./config.json")
	handleError(err)
	err = json.Unmarshal(data, &config)
	handleError(err)
	return config, handleError(err)
}

// func parseRarity(rarity map[string]interface{}, callbacks ...func(rarity float64)) {
// 	for _, val := range rarity {
// 		if reflect.TypeOf(val) == reflect.TypeOf(make(map[string]interface{})) {
// 			parseRarity(val.(map[string]interface{}), callbacks...)
// 		} else if  {
// 			for i := 0; i < len(callbacks); i++ {
// 				callback := callbacks[i]
// 				callback(val.(float64))
// 			}
// 		}
// 	}
// }