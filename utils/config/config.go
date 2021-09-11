package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/clickpop/looks/errors"
)

type Config struct {
	Input 					 InputObject 																	`json:"input,omitempty"`
	Output 					 OutputObject 																`json:"output,omitempty"`
	Settings 				 ConfigSettings 															`json:"setting,omitempty"`
	Attributes       map[string]map[string]map[string]interface{} `json:"attributes,omitempty"`
	DescriptionData  map[string]DescriptionData                   `json:"description-data,omitempty"`
}
type InputObject struct {
	Local InputLocalObject `json:"local,omitempty"`
}

type InputLocalObject struct {
	Filename string `json:"filename,omitempty"`
	Pathname string `json:"pathname,omitempty"`
}

type OutputObject struct {
	Local    	 OutputLocalObject `json:"local,omitempty"`
	Internal 	 bool							 `json:"internal,omitempty"`
	ImageCount float64 				 	 `json:"image-count,omitempty"`
}

type OutputLocalObject struct {
	Directory string `json:"directory,omitempty"`
}

type ConfigSettings struct {
	PieceOrder []string `json:"piece-order,omitempty"`
	MaxWorkers float64 	`json:"max-workers,omitempty"`
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