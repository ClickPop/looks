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
	Pathname string `json:"pathname,omitempty"`
	OutputDirectory string `json:"output-directory,omitempty"`
	OutputImageCount float64 `json:"output-image-count,omitempty"`
	Attributes map[string]map[string]map[string]float64 `json:"attributes,omitempty"`
}

type Metadata struct {
	Piece string `json:"piece"`
	Type string `json:"type"`
	Attributes map[string]float64 `json:"attributes"`
}

type FinalData struct {
	Pieces map[string]string `json:"pieces"`
	Rarity int `json:"rarity"`
	Cuteness int `json:"cuteness"`
	Rattitude int `json:"rattitude"`
}

func main() {
	config, err := loadJSON();
	if err != nil {
		handleError(err, exit)
	}
	_, err = os.Stat(config.OutputDirectory)
	if os.IsNotExist(err) {
		os.Mkdir(config.OutputDirectory, 0777)
	}
	for i := 0; i < int(config.OutputImageCount); i++ {
		makeFile(config, i)
	}
	fmt.Printf("Generated %d files in directory %s\n", int(config.OutputImageCount), config.OutputDirectory);
}

func makeFile(config Config, i int)  {
	fmt.Printf("Loading files for image #%d\n", i)
	files, metadata, err := loadFiles(config)
	handleError(err)
	fmt.Printf("Decoding data for image #%d\n", i)
	images, err := getImages(files)
	handleError(err)

	rect := image.Rectangle{images[1].Bounds().Min, images[1].Bounds().Max}
	img := image.NewRGBA(rect)
	origin := image.Point{0, 0}
	fmt.Printf("Layering assets for image #%d\n", i)
	for i := 0; i < len(images); i++ {
		currImg := images[i]
		if i == 0 {
			draw.Draw(img, currImg.Bounds(), currImg, origin, draw.Src)
		} else {
			draw.Draw(img, currImg.Bounds(), currImg, origin, draw.Over)
		}
	}

	out, err := os.Create(fmt.Sprintf("%s/%d.png", config.OutputDirectory, i))
	handleError(err)
	png.Encode(out, img)
	fmt.Printf("Image #%d.png created\n", i)
	var finalMeta FinalData
	finalMeta.Pieces = make(map[string]string, len(metadata))
	for j := 0; j < len(metadata); j++ {
		currMeta := metadata[j]
		finalMeta.Pieces[currMeta.Type] = currMeta.Piece
		finalMeta.Rarity += int(currMeta.Attributes["rarity"])
		finalMeta.Cuteness += int(currMeta.Attributes["cuteness"])
		finalMeta.Rattitude += int(currMeta.Attributes["rattitude"])
	}
	jsonData, err := json.MarshalIndent(finalMeta, "", "  ")
	handleError(err)
	err = os.WriteFile(fmt.Sprintf("%s/%d.json", config.OutputDirectory, i), jsonData, 0666)
	handleError(err)
	fmt.Printf("Metadata #%d.json created\n", i)
}

func loadFiles(config Config) ([]*os.File, []Metadata, error) {
	fileNames := config.PieceOrder
	var files []*os.File
	var metadata []Metadata
	for i := 0; i < len(fileNames); i++ {
		file := fileNames[i]
		piece, meta := handleRarity(config.Attributes[file])
		if piece != "nil" {
			filename := fmt.Sprintf(config.Filename, file, piece)
			reader, err := os.Open(fmt.Sprintf("%s/%s", config.Pathname, filename))
			if err != nil {
				return nil, nil, handleError(err, exit)
			}
			files = append(files, reader)
			metadata = append(metadata, Metadata{Type: file, Piece: piece, Attributes: meta})
		}
	}
	return files, metadata, nil
}

func handleRarity(pieceTypes map[string]map[string]float64) (string, map[string]float64) {
	var chances []string
	for key, val := range pieceTypes {
		rarities := val
		for _, v := range rarities {
			for i := 0; i < int(v); i++ {
				chances = append(chances, key)
			}
		}
	}

	seed := time.Now().Unix() + int64(time.Now().Nanosecond())
	rand.Seed(seed)

	random := rand.Intn(len(chances))
	choice := chances[random]
	return choice, pieceTypes[choice]
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