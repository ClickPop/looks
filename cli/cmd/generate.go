package cmd

import (
	"bytes"
	"fmt"
	"log"
	"time"
    "sync"

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
	        wg := sync.WaitGroup{}
            go func(logChan chan string) {
                wg.Add(1)
				for msg := range logChan {
					log.Println(msg)
				}
                wg.Done()
                close(logChan)
			}(logChan)
			go func(errChan chan error) {
				for err := range errChan {
					log.Panicln(err)
				}
                close(errChan)
			}(errChan)
            go func(cfg *conf.Config, asset chan generator.GeneratedAsset, csv *bytes.Buffer) {
                wg.Add(1)
                for asset := range assets {
                    generator.StoreFile(cfg, &asset)
                    if cfg.Output.MetaFormat == conf.CSV {
                        csv.Write(asset.Meta.Bytes())
                    }
                }
                wg.Done()
            }(cfg, assets, csv)
			generator.Generate(cfg, nil, assets, logChan, errChan)
            wg.Wait()
			fmt.Printf("Generated %d files in %d seconds.\n", cfg.Output.ImageCount, int(time.Since(startTime).Seconds()))
			return nil
		},
	}
)
