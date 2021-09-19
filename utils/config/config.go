package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Input        InputObject            `json:"input,omitempty"`
	Output       OutputObject           `json:"output,omitempty"`
	Settings     ConfigSettings         `json:"settings,omitempty"`
	Attributes   map[string]ConfigPiece `json:"attributes,omitempty"`
	Descriptions ConfigDescriptions     `json:"descriptions,omitempty"`
}
type InputObject struct {
	Local InputLocalObject `json:"local,omitempty"`
}

type InputLocalObject struct {
	Filename string `json:"filename,omitempty"`
	Pathname string `json:"pathname,omitempty"`
}

type OutputObject struct {
	Local         OutputLocalObject `json:"local,omitempty"`
	Internal      bool              `json:"internal,omitempty"`
	ImageCount    float64           `json:"image-count,omitempty"`
	MinimumRarity string            `json:"minimum-rarity,omitempty"`
}

type OutputLocalObject struct {
	Directory string `json:"directory,omitempty"`
}

type ConfigSettings struct {
	PieceOrder []string              `json:"piece-order,omitempty"`
	Stats      map[string]ConfigStat `json:"stats,omitempty"`
	Attributes map[string]ConfigAttribute		 `json:"attributes,omitempty"`
	Rarity     ConfigRarity          `json:"rarity,omitempty"`
	MaxWorkers float64               `json:"max-workers,omitempty"`
}

type ConfigDescriptions struct {
	Template            string                            `json:"template,omitempty"`
	FallbackPrimaryStat string                            `json:"fallback-primary-stat,omitempty"`
	HobbiesCount        int                               `json:"hobbies-count,omitempty"`
	Types               map[string]ConfigDescriptionTypes `json:"types,omitempty"`
}
type ConfigDescriptionTypes struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Descriptors []string `json:"descriptors,omitempty"`
	Hobbies     []string `json:"hobbies,omitempty"`
}

type ConfigStat struct {
	Name    string `json:"name,omitempty"`
	Minimum int    `json:"minimum,omitempty"`
	Maximum int    `json:"maximum,omitempty"`
	Value   int
}

type ConfigAttribute struct {
	Name 	string 			`json:"name,omitempty"`
	Type 	string 			`json:"type,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type ConfigRarity struct {
	Order   []string       `json:"order,omitempty"`
	Chances map[string]int `json:"chances,omitempty"`
}

type PieceAttribute struct {
	Rarity       string         `json:"rarity,omitempty"`
	Stats        map[string]int `json:"stats,omitempty"`
	FriendlyName string         `json:"friendly-name,omitempty"`
}

type ConfigPiece struct {
	FriendlyName string                     `json:"friendly-name,omitempty"`
	Pieces       map[string]PieceAttribute `json:"pieces,omitempty"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	filePath := path
	if path == "" {
		filePath = "./config.json"
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
