package cmd

import (
	"log"
	"sync"

	"github.com/clickpop/looks/pkg/generator"
	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Command to generate images/meta",
		Long:  "Generate images/metadata based on supplied files/config",
		RunE: func(cmd *cobra.Command, args []string) error {
      wg := new(sync.WaitGroup)
      wg.Add(1)
      logChan := make(chan string)
      errChan := make(chan error)
      go func(logChan <-chan string) {
        for msg := range logChan{
          log.Println(msg)
        }
        wg.Done()
      }(logChan)
      go func(errChan <-chan error) {
        for err := range errChan {
          log.Panicln(err)
        }
      }(errChan)
			_, err := generator.Generate(cfg, nil, logChan, errChan)
			if err != nil {
				return err
			}
      wg.Wait()
			return nil
		},
	}
)
