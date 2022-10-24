package config

type ConfigFormat string
type MetaFormat string
type DescriptionFormat string

const (
	JSON MetaFormat = "json"
	CSV  MetaFormat = "csv"
)

const (
	JSON_CONFIG ConfigFormat = "json"
	YAML_CONFIG ConfigFormat = "yaml"
	TOML_CONFIG ConfigFormat = "toml"
)

type Config struct {
	Input        InputObject            `json:"input" yaml:"input" toml:"input" mapstructure:"input"`
	Output       OutputObject           `json:"output" yaml:"output" toml:"output" mapstructure:"output"`
	Settings     ConfigSettings         `json:"settings" yaml:"settings" toml:"settings" mapstructure:"settings"`
	Attributes   map[string]ConfigPiece `json:"attributes" yaml:"attributes" toml:"attributes" mapstructure:"attributes"`
	Descriptions ConfigDescriptions     `json:"descriptions" yaml:"descriptions" toml:"descriptions" mapstructure:"descriptions"`
}
type InputObject struct {
	Local InputLocalObject `json:"local" yaml:"local" toml:"local" mapstructure:"local"`
}

type InputLocalObject struct {
	Filename string `json:"filename" yaml:"filename" toml:"filename" mapstructure:"filename"`
	Pathname string `json:"pathname" yaml:"pathname" toml:"pathname" mapstructure:"pathname"`
}

type OutputObject struct {
	Local         OutputLocalObject `json:"local" yaml:"local" toml:"local" mapstructure:"local"`
	Internal      bool              `json:"internal" yaml:"internal" toml:"internal" mapstructure:"internal"`
	ImageCount    float64           `json:"image-count" yaml:"image-count" toml:"image-count" mapstructure:"image-count"`
	IncludeMeta   bool              `json:"include-meta" yaml:"include-meta" toml:"include-meta" mapstructure:"include-meta"`
	MetaFormat    MetaFormat        `json:"meta-format" yaml:"meta-format" toml:"meta-format" mapstructure:"meta-format"`
	MinimumRarity string            `json:"minimum-rarity" yaml:"minimum-rarity" toml:"minimum-rarity" mapstructure:"minimum-rarity"`
}

type OutputLocalObject struct {
	Directory string `json:"directory" yaml:"directory" toml:"directory" mapstructure:"directory"`
}

type ConfigSettings struct {
	PieceOrder []string                   `json:"piece-order" yaml:"piece-order" toml:"piece-order" mapstructure:"piece-order"`
	Stats      map[string]ConfigStat      `json:"stats" yaml:"stats" toml:"stats" mapstructure:"stats"`
	Attributes map[string]ConfigAttribute `json:"attributes" yaml:"attributes" toml:"attributes" mapstructure:"attributes"`
	Rarity     ConfigRarity               `json:"rarity" yaml:"rarity" toml:"rarity" mapstructure:"rarity"`
	MaxWorkers float64                    `json:"max-workers" yaml:"max-workers" toml:"max-workers" mapstructure:"max-workers"`
	Tags       TagConfigSettings          `json:"tags" yaml:"tags" toml:"tags" mapstructure:"tags"`
}

type TagConfigSettings struct {
	Inclusive map[string][]string `json:"inclusive" yaml:"inclusive" toml:"inclusive" mapstructure:"inclusive"`
	Exclusive map[string][]string `json:"exclusive" yaml:"exclusive" toml:"exclusive" mapstructure:"exclusive"`
}

type ConfigDescriptions struct {
	Template            string                            `json:"template" yaml:"template" toml:"template" mapstructure:"template"`
	FallbackPrimaryStat string                            `json:"fallback-primary-stat" yaml:"fallback-primary-stat" toml:"fallback-primary-stat" mapstructure:"fallback-primary-stat"`
	FragmentCount       int                               `json:"fragment-count" yaml:"fragment-count" toml:"fragment-count" mapstructure:"fragment-count"`
	StatFragments       map[string]ConfigDescriptionTypes `json:"stat-fragments" yaml:"stat-fragments" toml:"stat-fragments" mapstructure:"stat-fragments"`
	SimpleFragments     []string                          `json:"simple-fragments" yaml:"simple-fragments" toml:"simple-fragments" mapstructure:"simple-fragments"`
}

type ConfigDescriptionTypes struct {
	Name        string   `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
	Descriptors []string `json:"descriptors" yaml:"descriptors" toml:"descriptors" mapstructure:"descriptors"`
	Hobbies     []string `json:"hobbies" yaml:"hobbies" toml:"hobbies" mapstructure:"hobbies"`
}

type ConfigStat struct {
	Name    string `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
	Minimum int    `json:"minimum" yaml:"minimum" toml:"minimum" mapstructure:"minimum"`
	Maximum int    `json:"maximum" yaml:"maximum" toml:"maximum" mapstructure:"maximum"`
	Value   int
}

type ConfigAttribute struct {
	Name  string      `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
	Type  string      `json:"type" yaml:"type" toml:"type" mapstructure:"type"`
	Value interface{} `json:"value" yaml:"value" toml:"value" mapstructure:"value"`
}

type ConfigRarity struct {
	Order   []string       `json:"order" yaml:"order" toml:"order" mapstructure:"order"`
	Chances map[string]int `json:"chances" yaml:"chances" toml:"chances" mapstructure:"chances"`
}

type PieceAttribute struct {
	Rarity       string         `json:"rarity" yaml:"rarity" toml:"rarity" mapstructure:"rarity"`
	Stats        map[string]int `json:"stats" yaml:"stats" toml:"stats" mapstructure:"stats"`
	Tags         []string       `json:"tags" yaml:"tags" toml:"tags" mapstructure:"tags"`
	FriendlyName string         `json:"friendly-name" yaml:"friendly-name" toml:"friendly-name" mapstructure:"friendly-name"`
}

type ConfigPiece struct {
	FriendlyName string                    `json:"friendly-name" yaml:"friendly-name" toml:"friendly-name" mapstructure:"friendly-name"`
	Pieces       map[string]PieceAttribute `json:"pieces" yaml:"pieces" toml:"pieces" mapstructure:"pieces"`
}
