package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/clickpop/looks/errors"
)

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

type DescriptionData struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Descriptors []string `json:"descriptors,omitempty"`
	Hobbies     []string `json:"hobbies,omitempty"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	filePath := path
	if path == "" {
		filePath = "./config.json"
	}
	data, err := ioutil.ReadFile(filePath)
	errors.HandleError(err)
	err = json.Unmarshal(data, &config)
	errors.HandleError(err)
	return config, errors.HandleError(err)
}