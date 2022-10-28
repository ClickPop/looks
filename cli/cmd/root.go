package cmd

import (
	"log"
	"os"
	"strings"

	conf "github.com/clickpop/looks-cli/pkg/config"
	"github.com/clickpop/looks/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	cfg         *config.Config = &config.Config{}
	cfgPathname string
	cfgFilename string
	cfgFiletype string
	rootCmd     = &cobra.Command{
		Use:   "looks",
		Short: "Looks generator CLI",
		Long:  "Looks is a tool for artists to have easier access to generative art",
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a skeleton config file",
		Long:  "Builds a skeleton config file with sensible defaults where possible",
		Run: func(cmd *cobra.Command, args []string) {
			storeConfig()
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.json", "Path to a config file to be used by the Looks generator")
	initCmd.PersistentFlags().StringVar(&cfgFilename, "name", "config", "Filename to use for generated config file")
	initCmd.PersistentFlags().StringVar(&cfgPathname, "path", "./", "Path to use for generated config file")
	initCmd.PersistentFlags().StringVar(&cfgFiletype, "type", "json", "Filetype to use for generated config file. Currently supported types are: json, yaml")
	generateCmd.PersistentFlags().BoolVar(&cfg.Output.IncludeMeta, "meta", true, "If generator should build meta")
	generateCmd.PersistentFlags().IntVar(&cfg.Output.ImageCount, "count", 100, "Number of assets to create")
	generateCmd.PersistentFlags().IntVar(&cfg.Settings.MaxWorkers, "workers", 3, "Number of workers to spin up. WARNING: Setting this higher than default will use more resources and might make the program unstable")
	viper.BindPFlag("output.image-count", generateCmd.PersistentFlags().Lookup("count"))
	viper.BindPFlag("output.include-meta", generateCmd.PersistentFlags().Lookup("meta"))
	viper.BindPFlag("output.num-workers", generateCmd.PersistentFlags().Lookup("workers"))
	initConfig()
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok || strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "cannot find") {
			for _, arg := range os.Args {
				if arg == "init" {
					break
				}
			}
		} else {
			log.Fatal(err)
		}
	}
	viper.Unmarshal(&cfg)
	conf.ValidateConfig(cfg)
}

func storeConfig() {
	viper.Reset()
	viper.AddConfigPath(cfgPathname)
	viper.SetConfigType(cfgFiletype)
	viper.SetConfigName(cfgFilename)
	viper.SetDefault("input", &config.InputObject{Local: config.InputLocalObject{Pathname: "pieces", Filename: "%s-%s.png"}})
	viper.SetDefault("output", &config.OutputObject{Directory: "generated", ImageCount: 100, IncludeMeta: true, MetaFormat: config.JSON})
	viper.SetDefault("settings", &config.ConfigSettings{MaxWorkers: 3})
	viper.SetDefault("attributes", map[string]config.ConfigPiece{})
	viper.SetDefault("descriptions", config.ConfigDescriptions{})
	err := viper.SafeWriteConfig()
	if err != nil {
		log.Fatal(err)
	}
}
