package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	conf "github.com/clickpop/looks/pkg/config"
	"github.com/clickpop/looks/pkg/generator"
	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Command to generate images/meta",
		Long:  "Generate images/metadata based on supplied files/config",
		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()
			logChan := make(chan string)
			errChan := make(chan error)
			assets := make(chan generator.GeneratedAsset, cfg.Output.ImageCount)
			heading := generator.BuildCsvHeading(cfg)
			csv := bytes.NewBuffer(*heading)
            go func(logChan chan string) {
				for msg := range logChan {
					log.Println(msg)
				}
			}(logChan)
			go func(errChan chan error) {
				for err := range errChan {
					log.Panicln(err)
				}
			}(errChan)
			generator.Generate(cfg, nil, assets, logChan, errChan)
            for asset := range assets {
                log.Printf("Storing files for asset %s\n", asset.Name)
                err := generator.StoreFile(cfg, asset)
                if err != nil {
                    return err
                }
                if cfg.Output.MetaFormat == conf.CSV {
                    csv.Write(asset.Meta.Bytes())
                }
            }
            if cfg.Output.MetaFormat == conf.CSV {
                log.Printf("Writing metadata file to %s/metadata.csv", cfg.Output.Directory)
                err := os.WriteFile(fmt.Sprintf("%s/metadata.csv", cfg.Output.Directory), csv.Bytes(), 0666)
                if err != nil {
                    return err
                }
            }
			log.Printf("Generated %d files in %d seconds.\n", cfg.Output.ImageCount, int(time.Since(startTime).Seconds()))
			return nil
		},
	}
)
