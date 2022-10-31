package generator

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"

	conf "github.com/clickpop/looks/pkg/config"
	"github.com/clickpop/looks/pkg/logging"
)

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
	MaxValue    float64     `json:"max_value,omitempty"`
}
type GeneratedAsset struct {
	Name  string
	Image *bytes.Buffer
	Meta  *bytes.Buffer
	Hash  string
}

type Job struct {
	id      int
	results chan GeneratedAsset
	errors  chan error
	currImg *image.RGBA
	outFile *os.File
}

var wg sync.WaitGroup

func Generate(config *conf.Config, hashCheckCb func(config *conf.Config), assets chan<- GeneratedAsset, logChan chan<- string, errChan chan<- error) {
	image_count := int(config.Output.ImageCount)
    wg.Add(image_count);

	jobs := make(chan int, image_count)

	num_workers := int(config.Settings.MaxWorkers)
	if image_count < num_workers {
		num_workers = image_count
	}
	for w := 0; w < num_workers; w++ {
		go handleJob(w, config, jobs, assets, errChan, logChan)
	}
	logChan <- fmt.Sprintf("spinning up %d workers", num_workers)
	for i := 0; i < image_count; i++ {
		jobs <- i
	}
	close(jobs)
    wg.Wait()
	close(assets)
	close(errChan)
	close(logChan)
}

func handleJob(worker int, config *conf.Config, jobs <-chan int, assets chan<- GeneratedAsset, errChan chan<- error, logChan chan<- string) {
	stats := buildStats(config)
    for job := range jobs {
		logChan <- logging.FormatLog(worker, job, "start asset build")
		buildAsset(config, worker, job, stats, assets, logChan, errChan)
        wg.Done()
	}
}

func computeLayers(images <-chan *ImageWithPosition, size int) []*image.Image {
	layers := make([]*image.Image, size)
	for img := range images {
		layers[img.Position] = img.Image
	}
	return layers
}

func computeMeta(metadata <-chan *[]byte) []byte {
	for data := range metadata {
		return *data
	}
	return nil
}

func buildAsset(config *conf.Config, worker int, job int, stats map[string]int, assets chan<- GeneratedAsset, logChan chan<- string, errChan chan<- error) {
	size := len(config.Settings.PieceOrder)
	files := make(chan *FileWithPosition, size)
	images := make(chan *ImageWithPosition, size)
	metadata := make(chan *[]byte)
	logChan <- logging.FormatLog(worker, job, "selecting pieces")
	pieces, pieceMeta := loadPieces(config)
	logChan <- logging.FormatLog(worker, job, "loading files")
	go loadFiles(pieces, files, errChan)
	logChan <- logging.FormatLog(worker, job, "decoding images")
	go getImages(files, images, errChan)
	if config.Output.IncludeMeta {
		logChan <- logging.FormatLog(worker, job, "generating meta")
		go generateMeta(pieceMeta, metadata, errChan, config, job)
	}
	layers := computeLayers(images, size)
	logChan <- logging.FormatLog(worker, job, "building image")
	img := buildImage(layers)
	imageOut := new(bytes.Buffer)
	metaOut := new(bytes.Buffer)
	meta := computeMeta(metadata)
	png.Encode(imageOut, img)
	metaOut.Write(meta)
	assets <- GeneratedAsset{Image: imageOut, Meta: metaOut, Name: fmt.Sprint(job)}
}
