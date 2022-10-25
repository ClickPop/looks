package generator

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"

	conf "github.com/clickpop/looks/pkg/config"
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
	MaxValue    int         `json:"max_value,omitempty"`
}
type GeneratedAsset struct {
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

type AssetsResult struct {
	Assets []GeneratedAsset
	CSV    io.Reader
}

func Generate(config *conf.Config, hashCheckCb func(config *conf.Config)) (AssetsResult, error) {
	image_count := int(config.Output.ImageCount)

	jobs := make(chan int, image_count)
	results := make(chan GeneratedAsset, image_count)
	errChan := make(chan error)

	num_workers := config.Settings.MaxWorkers
	for w := 0; w < int(num_workers); w++ {
		go handleJob(config, jobs, results, errChan)
	}
	for i := 0; i < image_count; i++ {
		jobs <- i
	}
	close(jobs)
	var assets []GeneratedAsset
	heading := buildCsvHeading(config)
	csv := bytes.NewBuffer(*heading)
	for {
		select {
		case asset, open := <-results:
			assets = append(assets, asset)
			if config.Output.MetaFormat == conf.CSV {
				csv.Write(asset.Meta.Bytes())
			}
			if !open {
				return AssetsResult{Assets: assets, CSV: csv}, nil
			}
		case err := <-errChan:
			if err.Error() == "abort" {
				return AssetsResult{}, nil
			}
		}
	}
}

func handleJob(config *conf.Config, jobs <-chan int, results chan<- GeneratedAsset, errChan chan<- error) {
	stats := buildStats(config)
	for job := range jobs {
    asset, err := buildAsset(config, job, stats)
		if err != nil {
			errChan <- err
		}
		if results != nil {
			results <- asset
		}
	}
  close(results)
  close(errChan)
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

func buildAsset(config *conf.Config, jobId int, stats map[string]int) (GeneratedAsset, error) {
	size := len(config.Settings.PieceOrder)
	files := make(chan *FileWithPosition, size)
	images := make(chan *ImageWithPosition, size)
	metadata := make(chan *[]byte)
	errChan := make(chan error)
	pieces, pieceMeta := loadPieces(config)
	go loadFiles(pieces, files, errChan)
	go getImages(files, images, errChan)
	if config.Output.IncludeMeta {
		go generateMeta(pieceMeta, metadata, errChan, config, jobId)
	}
	layers := computeLayers(images, size)
	img := buildImage(layers)
	imageOut := new(bytes.Buffer)
	metaOut := new(bytes.Buffer)
	meta := computeMeta(metadata)
	if config.Output.Internal {
		png.Encode(imageOut, img)
		metaOut.Write(meta)
	} else {
		imageOut = nil
		metaOut = nil
	}
	if config.Output.Local.Directory != "" {
		err := storeFile(config, img, meta, jobId)
		if err != nil {
			return GeneratedAsset{}, nil
		}
	}
	return GeneratedAsset{Image: imageOut, Meta: metaOut}, nil
}
