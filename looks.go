package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Config struct {
	PieceOrder       []string                                 `json:"piece-order,omitempty"`
	Filename         string                                   `json:"filename,omitempty"`
	Pathname         string                                   `json:"pathname,omitempty"`
	OutputDirectory  string                                   `json:"output-directory,omitempty"`
	OutputImageCount float64                                  `json:"output-image-count,omitempty"`
	Attributes       map[string]map[string]map[string]float64 `json:"attributes,omitempty"`
	MaxWorkers       float64                                  `json:"max-workers,omitempty"`
	DescriptionData  map[string]DescriptionData               `json:"description-data,omitempty"`
}

type Metadata struct {
	Piece      string             `json:"piece"`
	Type       string             `json:"type"`
	Attributes map[string]float64 `json:"attributes"`
}

type FinalData struct {
	Pieces      map[string]string `json:"pieces"`
	Rarity      int               `json:"rarity"`
	Cunning     int               `json:"cunning"`
	Cuteness    int               `json:"cuteness"`
	Rattitude   int               `json:"rattitude"`
	Description string            `json:"description"`
}

type DescriptionData struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Descriptors []string `json:"descriptors,omitempty"`
	Hobbies     []string `json:"hobbies,omitempty"`
}

type Job struct {
	id     int
	config Config
}

func main() {
	startTime := time.Now()
	config, err := loadJSON()
	if err != nil {
		handleError(err, exit)
	}
	_, err = os.Stat(config.OutputDirectory)
	if os.IsNotExist(err) {
		os.Mkdir(config.OutputDirectory, 0777)
	}

	jobs := make(chan Job, int(config.OutputImageCount))

	num_workers := config.MaxWorkers

	if num_workers == 0 {
		num_workers = 3
	}

	log.Printf("Spinning up %d workers", int(num_workers))
	for w := 0; w < int(num_workers); w++ {
		wg.Add(1)
		go makeFile(jobs)
	}

	for i := 0; i < int(config.OutputImageCount); i++ {
		jobs <- Job{i, config}
	}
	close(jobs)

	wg.Wait()
	log.Printf("Generated %d files in directory %s in %d seconds.\n", int(config.OutputImageCount), config.OutputDirectory, int(time.Since(startTime).Seconds()))
	checkHashes()
}

func makeFile(jobs <-chan Job) {
	for job := range jobs {
		config := job.config
		i := job.id
		log.Printf("Loading files for image #%d\n", i)
		files, metadata, err := loadFiles(config)
		handleError(err)
		log.Printf("Decoding data for image #%d\n", i)
		images, err := getImages(files)
		handleError(err)

		rect := image.Rectangle{images[1].Bounds().Min, images[1].Bounds().Max}
		img := image.NewRGBA(rect)
		origin := image.Point{0, 0}
		log.Printf("Layering assets for image #%d\n", i)
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
		out.Close()
		log.Printf("Image #%d.png created\n", i)
		var finalMeta FinalData
		finalMeta.Pieces = make(map[string]string, len(metadata))
		for j := 0; j < len(metadata); j++ {
			currMeta := metadata[j]
			finalMeta.Pieces[currMeta.Type] = currMeta.Piece
			finalMeta.Rarity += int(currMeta.Attributes["rarity"])
			finalMeta.Cunning += int(currMeta.Attributes["cunning"])
			finalMeta.Cuteness += int(currMeta.Attributes["cuteness"])
			finalMeta.Rattitude += int(currMeta.Attributes["rattitude"])
		}
		finalMeta.Description = config.buildDescription(finalMeta)
		jsonData, err := json.MarshalIndent(finalMeta, "", "  ")
		handleError(err)
		err = os.WriteFile(fmt.Sprintf("%s/%d.json", config.OutputDirectory, i), jsonData, 0666)
		handleError(err)
		log.Printf("Metadata #%d.json created\n", i)
	}
	wg.Done()
}

func loadFiles(config Config) ([]bytes.Reader, []Metadata, error) {
	fileNames := config.PieceOrder
	var files []bytes.Reader
	var metadata []Metadata
	for i := 0; i < len(fileNames); i++ {
		file := fileNames[i]
		piece, meta := handleRarity(config.Attributes[file])
		if piece != "nil" {
			filename := fmt.Sprintf(config.Filename, file, piece)
			data, err := os.ReadFile(fmt.Sprintf("%s/%s", config.Pathname, filename))
			if err != nil {
				return nil, nil, handleError(err, exit)
			}
			files = append(files, *bytes.NewReader(data))
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

func getImages(files []bytes.Reader) ([]image.Image, error) {
	var images []image.Image
	for i := 0; i < len(files); i++ {
		img, err := png.Decode(&files[i])
		if err != nil {
			return nil, handleError(err, exit)
		}
		images = append(images, img)
	}
	return images, nil
}

func handleError(err error, callbacks ...func(err error) error) error {
	if err != nil {
		log.Fatalln(err)
		for i := 0; i < len(callbacks); i++ {
			callback := callbacks[i]
			val := callback(err)
			if val != nil {
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

func checkHashes() {
	hashes := make(map[string]string)
	log.Println("Checking hashes for collisions...")
	filepath.WalkDir("rats", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if strings.Contains(d.Name(), ".png") {
			data, _ := os.ReadFile(path)
			reader := bytes.NewReader(data)
			hasher := md5.New()
			_, err := io.Copy(hasher, reader)
			if err != nil {
				log.Fatal(err)
			}
			hash := hex.EncodeToString(hasher.Sum(nil))
			if hashes[hash] != "" {
				log.Fatal("COLLISION", hashes[hash], path)
				os.Exit(1)
			}
			hashes[hash] = path
		}
		return err
	})
	log.Println("All hashes unique")
}

func (c *Config) buildDescription(meta FinalData) string {
	primaryStat := "default"

	if meta.Cunning > meta.Cuteness && meta.Cunning > meta.Rattitude {
		primaryStat = "cunning"
	} else if meta.Cuteness > meta.Cunning && meta.Cuteness > meta.Rattitude {
		primaryStat = "cuteness"
	} else if meta.Rattitude > meta.Cunning && meta.Rattitude > meta.Cuteness {
		primaryStat = "rattitude"
	}

	randomDescriptor := getRandomDescriptor(c.DescriptionData[primaryStat].Descriptors)
	randomHobbies := getRandomHobbies(c.DescriptionData[primaryStat].Hobbies, 3)
	ratType := c.DescriptionData[primaryStat].Name

	return fmt.Sprintf("This little rat is a %s, that means %s. Their favorite hobbies include %s.", ratType, randomDescriptor, randomHobbies)
}

func getRandomDescriptor(descriptors []string) string {
	rand.Seed(time.Now().Unix() + int64(time.Now().Nanosecond()))
	return descriptors[rand.Intn(len(descriptors))]
}

func getRandomHobbies(hobbies []string, n int) string {
	rand.Seed(time.Now().Unix() + int64(time.Now().Nanosecond()))
	var randomHobbies []string

	if len(hobbies) < n {
		n = len(hobbies)
	}

	for len(randomHobbies) < n {
		tempHobby := hobbies[rand.Intn(len(hobbies))]
		if !contains(randomHobbies, tempHobby) {
			randomHobbies = append(randomHobbies, tempHobby)
		}
	}

	return oxfordJoin(randomHobbies)
}

func contains(slice []string, haystack string) bool {
	rVal := false

	for _, v := range slice {
		if v == haystack {
			rVal = true
		}
	}

	return rVal
}

func oxfordJoin(slice []string) string {
	outStr := ""

	if len(slice) == 1 {
		outStr = slice[0]
	} else if len(slice) == 2 {
		outStr = strings.Join(slice, " and ")
	} else if len(slice) > 2 {
		outStr = strings.Join(slice[0:(len(slice)-1)], ", ") + ", and " + slice[len(slice)-1]
	}

	return outStr
}
