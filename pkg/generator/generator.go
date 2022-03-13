package generator

import (
	"bytes"
	CSV "encoding/csv"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	conf "github.com/clickpop/looks/pkg/config"
)

var wg sync.WaitGroup

type PieceMetadata struct {
	Piece        string
	Type         string
	Attributes   map[string]conf.ConfigStat
	Rarity       string
	FriendlyName string
}

type Metadata struct {
	Type      string
	PieceMeta []PieceMetadata
}

type OpenSeaMeta struct {
	Image           string             `json:"image,omitempty"`
	ImageData       string             `json:"image_data,omitempty"`
	ExternalURL     string             `json:"external_url,omitempty"`
	Description     string             `json:"description,omitempty"`
	Name            string             `json:"name,omitempty"`
	Attributes      []OpenSeaAttribute `json:"attributes,omitempty"`
	BackgroundColor string             `json:"background_color,omitempty"`
	AnimationURL    string             `json:"animation_url,omitempty"`
	YouTubeURL      string             `json:"youtube_url,omitempty"`
}

type OpenSeaAttribute struct {
	TraitType   string      `json:"trait_type,omitempty"`
	DisplayType string      `json:"display_type,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	MaxValue    int         `json:"max_value,omitempty"`
}
type GeneratedRat struct {
	Image *bytes.Buffer
	Meta  *bytes.Buffer
}

type Job struct {
	id      int
	results chan GeneratedRat
	errors  chan error
	currImg *image.RGBA
	outFile *os.File
}

func Generate(config *conf.Config, hashCheckCb func(config *conf.Config)) ([]GeneratedRat, error) {
	startTime := time.Now()
	buildCsvHeading(config)
	log.Println(csv)
	outputDir := config.Output.Local.Directory
	if (config.Output == conf.OutputObject{}) {
		config.Output.Internal = true
	}

	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		os.Mkdir(outputDir, 0777)
	} else if err != nil {
		return nil, err
	}

	image_count := int(config.Output.ImageCount)

	wg.Add(image_count)

	jobs := make(chan int, image_count)
	results := make(chan GeneratedRat, image_count)
	errChan := make(chan error)

	num_workers := config.Settings.MaxWorkers

	log.Printf("Spinning up %d workers", int(num_workers))
	for w := 0; w < int(num_workers); w++ {
		go handleJob(config, jobs, results, errChan)
	}
	for i := 0; i < image_count; i++ {
		jobs <- i
	}
	close(jobs)

	var assets []GeneratedRat
	for r := range results {
		assets = append(assets, r)
	}
	wg.Wait()
	if config.Output.MetaFormat == conf.CSV {
		metaFile, err := os.Create(fmt.Sprintf("%s/meta.csv", config.Output.Local.Directory))
		if err != nil {
			log.Fatal(err)
		}
		defer metaFile.Close()
		w := CSV.NewWriter(metaFile)
		err = w.WriteAll(csv)
		if err != nil {
			log.Fatal("Error writing csv:", err)
		}
	}
	log.Printf("Generated %d files in directory %s in %d seconds.\n", image_count, outputDir, int(time.Since(startTime).Seconds()))
	if outputDir != "" && hashCheckCb == nil {
		noCollisions := false
		for !noCollisions {
			collisions, err := checkHashes(outputDir)
			if err != nil {
				return nil, err
			}
			re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
			for _, path := range collisions {
				match := re.FindStringSubmatch(path)
				if len(match) > 2 {
					str := match[2]
					val, err := strconv.Atoi(str)
					if err != nil {
						return nil, err
					}
					stats := buildStats(config)
					asset, err := buildAsset(config, val, stats)
					if err != nil {
						return nil, err
					}
					assets = append(assets, asset)
				}
			}
			if len(collisions) < 1 {
				noCollisions = true
			}
		}
	}
	return assets, nil
}

func handleJob(config *conf.Config, jobs <-chan int, results chan<- GeneratedRat, errChan chan<- error) {
	stats := buildStats(config)

	for job := range jobs {
		i := job
		rat, err := buildAsset(config, i, stats)
		if err != nil {
			errChan <- err
		}
		if results != nil {
			results <- rat
		}
		wg.Done()
		config.Output.ImageCount -= 1
	}
	if config.Output.ImageCount == 0 {
		close(results)
		close(errChan)
	}
}

func buildAsset(config *conf.Config, jobId int, stats map[string]int) (GeneratedRat, error) {
	i := jobId
	log.Printf("Loading files for image #%d\n", i)
	files, metadata, err := loadFiles(config)
	if err != nil {
		return GeneratedRat{}, nil
	}
	log.Printf("Decoding data for image #%d\n", i)
	images, err := getImages(files)
	if err != nil {
		return GeneratedRat{}, nil
	}
	imageOut := new(bytes.Buffer)
	metaOut := new(bytes.Buffer)
	img := buildImage(images, i)
	var meta []byte
	if config.Output.IncludeMeta {
		meta, err = generateMeta(metadata, config, i)
	}
	if err != nil {
		return GeneratedRat{}, nil
	}
	if config.Output.Internal {
		png.Encode(imageOut, img)
		metaOut.Write(meta)
	} else {
		imageOut = nil
		metaOut = nil
	}
	if config.Output.Local.Directory != "" {
		err = storeFile(config, img, meta, i)
		if err != nil {
			return GeneratedRat{}, nil
		}
	}
	return GeneratedRat{Image: imageOut, Meta: metaOut}, nil
}
