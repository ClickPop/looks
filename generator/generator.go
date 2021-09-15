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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/clickpop/looks/utils"
	conf "github.com/clickpop/looks/utils/config"
)

var wg sync.WaitGroup

type Metadata struct {
	Piece        string
	Type         string
	Attributes   map[string]conf.ConfigStat
	Rarity       string
	FriendlyName string
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

func Generate(config *conf.Config, hashCheckCb func(config *conf.Config)) ([]GeneratedRat, error) {
	startTime := time.Now()

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

	if image_count == 0 {
		image_count = 10
		config.Output.ImageCount = float64(image_count)
	}
	wg.Add(image_count)

	jobs := make(chan int, image_count)
	results := make(chan GeneratedRat, image_count)
	errChan := make(chan error)

	num_workers := config.Settings.MaxWorkers

	if num_workers == 0 {
		num_workers = 3
		config.Settings.MaxWorkers = num_workers
	}

	log.Printf("Spinning up %d workers", int(num_workers))
	for w := 0; w < int(num_workers); w++ {
		go makeFile(config, jobs, results, errChan)
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
	log.Printf("Generated %d files in directory %s in %d seconds.\n", image_count, outputDir, int(time.Since(startTime).Seconds()))
	if outputDir != "" && hashCheckCb == nil {
		err = checkHashes(outputDir)
		if err != nil {
			return nil, err
		}
	}
	return assets, nil
}

func makeFile(config *conf.Config, jobs <-chan int, results chan<- GeneratedRat, errChan chan<- error) {
	stats := make(map[string]int)
	for k := range config.Settings.Stats {
		stats[k] = 0
	}

	for job := range jobs {
		i := job
		log.Printf("Loading files for image #%d\n", i)
		files, metadata, err := loadFiles(config)
		if err != nil {
			errChan <- err
		}
		log.Printf("Decoding data for image #%d\n", i)
		images, err := getImages(files)
		if err != nil {
			errChan <- err
		}
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
		finalMeta.Attributes = make([]OpenSeaAttribute, 0)
		attributes := make(map[string]conf.ConfigStat)
		for _, v := range config.Settings.Stats {
			attr := v
			attr.Value = 0
			attributes[attr.Name] = attr
		}
		for j := 0; j < len(metadata); j++ {
			currMeta := metadata[j]
			finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: currMeta.Type, Value: currMeta.Piece})
			sort.Slice(finalMeta.Attributes, func(i, j int) bool {
				return finalMeta.Attributes[i].TraitType < finalMeta.Attributes[j].TraitType
			})
			for _, v := range currMeta.Attributes {
				attr := attributes[v.Name]
				attr.Value += v.Value
				if attr.Value >= v.Maximum {
					attr.Value = v.Maximum
				} else if attr.Value <= v.Minimum {
					attr.Value = v.Minimum
				}
				attributes[v.Name] = attr
			}
		}
		for k, v := range attributes {
			finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: k, Value: v.Value, DisplayType: "number", MaxValue: v.Maximum})
		}
		finalMeta.Description = buildDescription(config, finalMeta)
		finalMeta.Name = fmt.Sprint(i)
		jsonData, err := json.MarshalIndent(finalMeta, "", "  ")
		if err != nil {
			errChan <- err
		}
		if config.Output.Internal {
			png.Encode(imageOut, img)
			metaOut.Write(jsonData)
		} else {
			imageOut = nil
			metaOut = nil
		}
		if config.Output.Local.Directory != "" {
			out, err := os.Create(fmt.Sprintf("%s/%d.png", config.Output.Local.Directory, i))
			if err != nil {
				errChan <- err
			}
			png.Encode(out, img)
			out.Close()
			log.Printf("Image #%d.png created\n", i)
			err = os.WriteFile(fmt.Sprintf("./%s/%d.json", config.Output.Local.Directory, i), jsonData, 0666)
			if err != nil {
				errChan <- err
			}
			log.Printf("Metadata #%d.json created\n", i)
		}
		if results != nil {
			results <- GeneratedRat{Image: imageOut, Meta: metaOut}
		}
		wg.Done()
		config.Output.ImageCount -= 1
	}
	if config.Output.ImageCount == 0 {
		close(results)
		close(errChan)
	}
}

func loadFiles(config *conf.Config) ([]*bytes.Reader, []Metadata, error) {
	fileNames := config.Settings.PieceOrder
	var files []*bytes.Reader
	var metadata []Metadata
	for i := 0; i < len(fileNames); i++ {
		file := fileNames[i]
		pieceTypeFriendlyName := config.Attributes[file].FriendlyName
		if pieceTypeFriendlyName == "" {
			pieceTypeFriendlyName = utils.TransformName(file)
		}
		piece, meta := handleRarity(config.Attributes[file].Pieces, config.Settings.Rarity)

		if piece != "nil" {
			pieceFriendlyName := config.Attributes[file].Pieces[piece].FriendlyName
			if pieceFriendlyName == "" {
				pieceFriendlyName = utils.TransformName(piece)
			}
			filename := fmt.Sprintf(config.Input.Local.Filename, file, piece)
			data, err := os.ReadFile(fmt.Sprintf("%s/%s", config.Input.Local.Pathname, filename))
			if err != nil {
				return nil, nil, err
			}
			files = append(files, bytes.NewReader(data))
			stats := make(map[string]conf.ConfigStat)
			for k, v := range meta.Stats {
				stat := config.Settings.Stats[k]
				name := stat.Name
				if name == "" {
					name = utils.TransformName(k)
				}
				stat.Value = v
				stats[k] = stat
			}
			metadata = append(metadata, Metadata{Type: pieceTypeFriendlyName, Piece: pieceFriendlyName, Attributes: stats, Rarity: meta.Rarity, FriendlyName: meta.FriendlyName})
		}
	}

	return files, metadata, nil
}

func getRarityDenominator(r conf.ConfigRarity) int {
	denominator := 0
	for _, v := range r.Chances {
		denominator += int(v)
	}

	return denominator
}

func getRarityLevel(r conf.ConfigRarity) []string {
	denominator := getRarityDenominator(r)

	seed := time.Now().Unix() + int64(time.Now().Nanosecond())
	rand.Seed(seed)

	random := rand.Intn(denominator)

	rarity := []string{r.Order[0]}

	currentThreshold := 0

	for i, v := range r.Order {
		currentThreshold += int(r.Chances[v])
		if random <= currentThreshold {
			if i != 0 {
				rarity = append(rarity, v)
			}
			break
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(rarity)))
	return rarity
}

func handleRarity(pieceTypes map[string]conf.ConfigAttribute, rarityData conf.ConfigRarity) (string, conf.ConfigAttribute) {
	rarityLevel := getRarityLevel(rarityData)
	var possiblePieces []string
	for _, v := range rarityLevel {
		for key, piece := range pieceTypes {
			if v == piece.Rarity {
				possiblePieces = append(possiblePieces, key)
			}
		}
		if len(possiblePieces) > 0 {
			break
		}
	}

	seed := time.Now().Unix() + int64(time.Now().Nanosecond())
	rand.Seed(seed)

	random := rand.Intn(len(possiblePieces))
	choice := possiblePieces[random]
	return choice, pieceTypes[choice]
}

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

func checkHashes(outputDir string) error {
	hashes := make(map[string]string)
	log.Println("Checking hashes for collisions...")
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		os.Mkdir(outputDir, 0777)
	} else if err != nil {
		return err
	}
	err = filepath.WalkDir(outputDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if strings.Contains(d.Name(), ".png") {
			data, _ := os.ReadFile(path)
			reader := bytes.NewReader(data)
			hasher := md5.New()
			_, err := io.Copy(hasher, reader)
			if err != nil {
				return err
			}
			hash := hex.EncodeToString(hasher.Sum(nil))
			if hashes[hash] != "" {
				log.Println("COLLISION", hashes[hash], path)
				return fmt.Errorf("COLLISION: %s & %s", hashes[hash], path)
			}
			hashes[hash] = path
		}
		return err
	})
	if err != nil {
		return err
	}
	log.Println("All hashes unique")
	return nil
}

func buildDescription(c *conf.Config, meta OpenSeaMeta) string {
	stats := make(map[string]int)
	for k := range c.Settings.Stats {
		stats[k] = 0
	}

	for _, v := range meta.Attributes {
		if _, isStat := stats[v.TraitType]; isStat {
			stats[v.TraitType] += v.Value.(int)
		}
	}

	primaryStat := getPrimaryStat(stats, c.Descriptions.FallbackPrimaryStat)
	randomDescriptor := getRandomDescriptor(c.Descriptions.Types[primaryStat].Descriptors)
	randomHobbies := getRandomHobbies(c.Descriptions.Types[primaryStat].Hobbies, c.Descriptions.HobbiesCount)
	currentType := c.Descriptions.Types[primaryStat].Name

	return fmt.Sprintf(c.Descriptions.Template, currentType, randomDescriptor, randomHobbies)
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
	complex := false
	punct := ","

	for _, v := range slice {
		if strings.Contains(v, ",") {
			complex = true
			punct = ";"
		}
	}

	if len(slice) == 1 {
		outStr = slice[0]
	} else if len(slice) == 2 && !complex {
		outStr = strings.Join(slice, " and ")
	} else if len(slice) > 2 || (len(slice) == 2 && complex) {
		outStr = strings.Join(slice[0:(len(slice)-1)], fmt.Sprintf("%s ", punct)) + fmt.Sprintf("%s and ", punct) + slice[len(slice)-1]
	}

	return outStr
}

func getPrimaryStat(stats map[string]int, fallbackPrimaryStat string) string {
	max := -(int(^uint(0) >> 1)) - 1
	primaryStat := fallbackPrimaryStat

	for stat, v := range stats {
		if v > max {
			primaryStat = stat
			max = v
		} else if v == max {
			primaryStat = fallbackPrimaryStat
		}
	}

	return primaryStat
}