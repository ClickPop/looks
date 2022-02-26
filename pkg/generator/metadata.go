package generator

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	conf "github.com/clickpop/looks/pkg/config"
)

func generateMeta(metadata Metadata, config *conf.Config, genMeta bool, i int) ([]byte, error) {
	var finalMeta OpenSeaMeta
		finalMeta.Attributes = make([]OpenSeaAttribute, 0)
		stats := make(map[string]conf.ConfigStat)
		for _, v := range config.Settings.Stats {
			attr := v
			attr.Value = 0
			stats[attr.Name] = attr
		}
		for j := 0; j < len(metadata.PieceMeta); j++ {
			currMeta := metadata.PieceMeta[j]
			finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: currMeta.Type, Value: currMeta.Piece})
			sort.Slice(finalMeta.Attributes, func(i, j int) bool {
				return finalMeta.Attributes[i].TraitType < finalMeta.Attributes[j].TraitType
			})
			for _, v := range currMeta.Attributes {
				attr := stats[v.Name]
				attr.Value += v.Value
				if attr.Value >= v.Maximum {
					attr.Value = v.Maximum
				} else if attr.Value <= v.Minimum {
					attr.Value = v.Minimum
				}
				stats[v.Name] = attr
			}
		}
		for k, v := range stats {
			finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: k, Value: v.Value, DisplayType: "number", MaxValue: v.Maximum})
		}
		for k, v := range config.Settings.Attributes {
			val := v.Value
			name := v.Name
			attrType := v.Type
			if name == "" {
				name = k
			}

			switch v.Type {
			case "timestamp":
				attrType = "date"
				val = time.Now().Unix()
			}
			finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: name, DisplayType: attrType, Value: val})
		}
		if genMeta {
			description, name := buildDescription(config, finalMeta)
			finalMeta.Description = description
			finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: "Type", Value: name})
		}
		finalMeta.Name = fmt.Sprint(i)
		jsonData, err := json.MarshalIndent(finalMeta, "", "  ")
		if err != nil {
			return nil, err
		}
		return jsonData, nil
}