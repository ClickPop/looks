package generator

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	conf "github.com/clickpop/looks/pkg/config"
)

func computeMetadata(finalMeta OpenSeaMeta, stats map[string]conf.ConfigStat, pieces []*PieceMetadata) (OpenSeaMeta, map[string]conf.ConfigStat) {
	for _, piece := range pieces {
    finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: piece.Type, Value: piece.Piece})
    sort.Slice(finalMeta.Attributes, func(i, j int) bool {
      return finalMeta.Attributes[i].TraitType < finalMeta.Attributes[j].TraitType
    })
    for _, v := range piece.Attributes {
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

  return finalMeta, stats
}

func generateMeta(pieces []*PieceMetadata, metadata chan<- *[]byte, errChan chan<- error, config *conf.Config, jobId int) {
	var finalMeta OpenSeaMeta
	finalMeta.Attributes = make([]OpenSeaAttribute, 0)
	stats := make(map[string]conf.ConfigStat)
	for _, v := range config.Settings.Stats {
		attr := v
		attr.Value = 0
		stats[attr.Name] = attr
	}
	var headings []string
	if config.Output.MetaFormat == conf.CSV {
		headings = getCSVCols(config)
	}

	finalMeta, stats = computeMetadata(finalMeta, stats, pieces)

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
    if name != "" {
      finalMeta.Attributes = append(finalMeta.Attributes, OpenSeaAttribute{TraitType: "Type", Value: name})
    }
	}
	finalMeta.Name = fmt.Sprint(jobId)
	switch config.Output.MetaFormat {
	case conf.JSON:
		jsonData, err := json.MarshalIndent(finalMeta, "", "  ")
		if err != nil {
			errChan <- err
		}
		metadata <- &jsonData
	case conf.CSV:
		metadata <- buildCsvRow(headings, finalMeta)
	}
}

func getCSVCols(config *conf.Config) []string {
	headings := make([]string, 0)
	headings = append(headings, "Name")
	headings = append(headings, "Description")
	for k := range config.Attributes {
		headings = append(headings, strings.Title(k))
	}
	for k := range config.Settings.Stats {
		headings = append(headings, strings.Title(k))
	}
	return headings
}

func buildCsvHeading(config *conf.Config) *[]byte {
	headings := getCSVCols(config)
	slice := []byte(strings.Join(headings, ","))
	return &slice
}

func buildCsvRow(headings []string, meta OpenSeaMeta) *[]byte {
	rowMap := make(map[string]interface{})
	rowMap["Name"] = meta.Name
	rowMap["Description"] = meta.Description
	for _, attribute := range meta.Attributes {
		rowMap[strings.Title(attribute.TraitType)] = attribute.Value
	}
	row := make([]string, 0)
	for _, col := range headings {
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
	slice := []byte(strings.Join(row, ","))
	return &slice
}
