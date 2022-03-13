package generator

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	conf "github.com/clickpop/looks/pkg/config"
)

var csv [][]string

func generateMeta(metadata Metadata, config *conf.Config, i int) ([]byte, error) {
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
	if config.Output.IncludeMeta {
		description, name := buildDescription(config, finalMeta)
		finalMeta.Description = description
		finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: "Type", Value: name})
	}
	finalMeta.Name = fmt.Sprint(i)
	switch config.Output.MetaFormat {
	case conf.JSON:
		jsonData, err := json.MarshalIndent(finalMeta, "", "  ")
		if err != nil {
			return nil, err
		}
		return jsonData, nil
	case conf.CSV:
		buildCsvRow(finalMeta)
		return nil, nil
	}
	return nil, nil
}

func buildCsvHeading(config *conf.Config) {
	csv = make([][]string, 0)
	headings := make([]string, 0)
	headings = append(headings, "Name")
	headings = append(headings, "Description")
	for k := range config.Attributes {
		headings = append(headings, strings.Title(k))
	}
	for k := range config.Settings.Stats {
		headings = append(headings, strings.Title(k))
	}
	csv = append(csv, headings)
}

func buildCsvRow(meta OpenSeaMeta) {
	rowMap := make(map[string]interface{})
	rowMap["Name"] = meta.Name
	rowMap["Description"] = meta.Description
	for _, attribute := range meta.Attributes {
		rowMap[strings.Title(attribute.TraitType)] = attribute.Value
	}
	row := make([]string, 0)
	for _, col := range csv[0] {
		if rowMap[col] != nil {
			if val, ok := rowMap[col].(string); ok {
				if val != "" {
					row = append(row, val)
				}
			}

			if val, ok := rowMap[col].(int); ok {
				if val != 0 {
					row = append(row, fmt.Sprint(val))
				}
			}
		} else {
			row = append(row, "")
		}
	}
	csv = append(csv, row)
}
