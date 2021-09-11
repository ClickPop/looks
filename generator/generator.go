package generator

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
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/clickpop/looks/errors"
	"github.com/clickpop/looks/utils"
	conf "github.com/clickpop/looks/utils/config"
)

var wg sync.WaitGroup

type Metadata struct {
	Piece      string             		`json:"piece"`
	Type       string             		`json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type OpenSeaMeta struct {
	Image 					string 						 `json:"image,omitempty"`
	ImageData 			string 						 `json:"image_data,omitempty"`
	ExternalURL 		string 						 `json:"external_url,omitempty"`
	Description 		string 						 `json:"description,omitempty"`
	Name       			string 						 `json:"name,omitempty"`
	Attributes 			[]OpenSeaAttribute `json:"attributes,omitempty"`
	BackgroundColor string 						 `json:"background_color,omitempty"`
	AnimationURL 		string 						 `json:"animation_url,omitempty"`
	YouTubeURL 			string 						 `json:"youtube_url,omitempty"`
}

type OpenSeaAttribute struct {
	TraitType 	string 			`json:"trait_type,omitempty"`
	DisplayType string 			`json:"display_type,omitempty"`
	Value 			interface{} `json:"value,omitempty"`
}

type Job struct {
	id     int
	config conf.Config
}

type GeneratedRat struct {
	Image *bytes.Buffer
	Meta *bytes.Buffer
}

func Generate(config conf.Config) []GeneratedRat {
	startTime := time.Now()

	outputDir := config.Output.Local.Directory

	if (outputDir == "") {
		outputDir = "./pieces"
		config.Output.Local.Directory = outputDir
	}
	
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		os.Mkdir(outputDir, 0777)
	}

	image_count := int(config.Output.ImageCount)

	if (image_count == 0) {
		image_count = 10
		config.Output.ImageCount = float64(image_count)
	}

	jobs := make(chan Job, image_count)
	results := make(chan GeneratedRat, image_count)

	num_workers := config.Settings.MaxWorkers

	if num_workers == 0 {
		num_workers = 3
		config.Settings.MaxWorkers = num_workers
	}

	log.Printf("Spinning up %d workers", int(num_workers))
	for w := 0; w < int(num_workers); w++ {
		wg.Add(1)
		go makeFile(jobs, results)
	}

	for i := 0; i < image_count; i++ {
		jobs <- Job{i, config}
	}
	close(jobs)
	var assets []GeneratedRat
	for r := range results {
		assets = append(assets, r)
	}
	wg.Wait()
	log.Printf("Generated %d files in directory %s in %d seconds.\n", image_count, outputDir, int(time.Since(startTime).Seconds()))
	checkHashes()
	return assets
}

func makeFile(jobs <-chan Job, results chan<- GeneratedRat) {
	for job := range jobs {
		config := job.config
		i := job.id
		log.Printf("Loading files for image #%d\n", i)
		files, metadata, err := loadFiles(&config)
		errors.HandleError(err)
		log.Printf("Decoding data for image #%d\n", i)
		images, err := getImages(files)
		errors.HandleError(err)
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
		imageOut := new(bytes.Buffer)
		metaOut := new(bytes.Buffer)
		var finalMeta OpenSeaMeta
		finalMeta.Attributes = make([]OpenSeaAttribute, 0);
		attributes := make(map[string]int, 0)
		for j := 0; j < len(metadata); j++ {
			currMeta := metadata[j]
			finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: currMeta.Type, Value: currMeta.Piece})
			for k, v := range currMeta.Attributes {
				if k != "rarity" {
					switch t := v.(type) {
						case float64:
							attributes[k] += int(t)
					}
				}
			}
		}
		for k, v := range attributes {
			finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: k, Value: v, DisplayType: "number"})
		}
		finalMeta.Description = buildDescription(&config, finalMeta)
		finalMeta.Name = fmt.Sprint(i)
		jsonData, err := json.MarshalIndent(finalMeta, "", "  ")
		errors.HandleError(err)
		if config.Output.Internal {
			png.Encode(imageOut, img)
			metaOut.Write(jsonData)
		} else {
			imageOut = nil
			metaOut = nil
		}
		if config.Output.Local.Directory != "" {
			out, err := os.Create(fmt.Sprintf("%s/%d.png", config.Output.Local.Directory, i))
			errors.HandleError(err)
			png.Encode(out, img)
			out.Close()
			log.Printf("Image #%d.png created\n", i)
			err = os.WriteFile(fmt.Sprintf("./%s/%d.json", config.Output.Local.Directory, i), jsonData, 0666)
			errors.HandleError(err)
			log.Printf("Metadata #%d.json created\n", i)
		}
		if results != nil {
			results <- GeneratedRat{Image: imageOut, Meta: metaOut}
		}
	}
	if results != nil {
		close(results)
	}
	wg.Done()
}

func loadFiles(config *conf.Config) ([]*bytes.Reader, []Metadata, error) {
	fileNames := config.Settings.PieceOrder
	var files []*bytes.Reader
	var metadata []Metadata
	for i := 0; i < len(fileNames); i++ {
		file := fileNames[i]
		piece, meta := handleRarity(config.Attributes[file])
		if piece != "nil" {
			filename := fmt.Sprintf(config.Input.Local.Filename, file, piece)
			data, err := os.ReadFile(fmt.Sprintf("%s/%s", config.Input.Local.Pathname, filename))
			if err != nil {
				return nil, nil, errors.HandleError(err, errors.Exit)
			}
			files = append(files, bytes.NewReader(data))
			metadata = append(metadata, Metadata{Type: file, Piece: piece, Attributes: meta})
		}
	}
	return files, metadata, nil
}

func handleRarity(pieceTypes map[string]map[string]interface{}) (string, map[string]interface{}) {
	var chances []string
	for key, val := range pieceTypes {
		rarities := val
		for _, v := range rarities {
			for i := 0; i < int(v.(float64)); i++ {
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

func getImages(files []*bytes.Reader) ([]*image.Image, error) {
	var images []*image.Image
	for i := 0; i < len(files); i++ {
		img, err := png.Decode(files[i])
		if err != nil {
			return nil, errors.HandleError(err, errors.Exit)
		}
		images = append(images, &img)
	}
	return images, nil
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

func buildDescription(c *conf.Config, meta OpenSeaMeta) string {
	primaryStat := "default"

	cunning := 0
	cuteness := 0
	rattitude := 0

	for _, v := range meta.Attributes {
		switch v.TraitType {
		case "cuteness":
			cuteness += v.Value.(int)
		case "cunning":
			cunning += v.Value.(int)
		case "rattitude":
			rattitude += v.Value.(int)
		} 
	}

	if cunning > cuteness && cunning > rattitude {
		primaryStat = "cunning"
	} else if cuteness > cunning && cuteness > rattitude {
		primaryStat = "cuteness"
	} else if rattitude > cunning && rattitude > cuteness {
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
		if !utils.Contains(randomHobbies, tempHobby) {
			randomHobbies = append(randomHobbies, tempHobby)
		}
	}

	return oxfordJoin(randomHobbies)
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