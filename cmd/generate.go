package cmd

import (
	"github.com/clickpop/looks/pkg/generator"
	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Command to generate images/meta",
		Long:  "Generate images/metadata based on supplied files/config",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := generator.Generate(cfg, nil, cfg.Output.IncludeMeta)
			if err != nil {
				return err
			}
			return nil
		},
	}
)
