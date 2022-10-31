package generator

import (
	"bytes"
	"fmt"
	"os"

	"github.com/clickpop/looks/internal/utils"
	conf "github.com/clickpop/looks/pkg/config"
)

type FileWithPosition struct {
	File     *bytes.Reader
	Position int
}

func loadFiles(pieces []string, files chan<- *FileWithPosition, errChan chan<- error) {
	for position, piece := range pieces {
		file, err := os.ReadFile(piece)
		if err != nil {
			errChan <- err
			return
		}
		files <- &FileWithPosition{File: bytes.NewReader(file), Position: position}
	}
	close(files)
}

func loadPieces(config *conf.Config) ([]string, []*PieceMetadata) {
	tags := make(map[string]bool)
	variant := ""
	pieces := make([]string, len(config.Settings.PieceOrder))
	metadata := make([]*PieceMetadata, len(config.Settings.PieceOrder))
	for position, pieceType := range config.Settings.PieceOrder {
		pieceTypeFriendlyName := config.Attributes[pieceType].FriendlyName
		if pieceTypeFriendlyName == "" {
			pieceTypeFriendlyName = utils.TransformName(pieceType)
		}
		piece, meta, v, pieceFriendlyName := handleRarity(config.Attributes[pieceType].Pieces, config, tags, variant, position)
		variant = v
		if piece != "nil" {
			fileName := fmt.Sprintf(config.Input.Local.Filename, pieceType, piece)
			filePath := fmt.Sprintf("%s/%s", config.Input.Local.Pathname, fileName)
			stats := make(map[string]conf.ConfigStat)
			for k, v := range meta.Stats {
				stat := config.Settings.Stats[k]
				name := stat.Name
				if name == "" {
					name = utils.TransformName(k)
				}
				stat.Value = v
				stats[k] = stat
			}
			for _, tag := range meta.Tags {
				if _, ok := tags[tag]; !ok {
					tags[tag] = true
				}
			}
			var rarity string
			switch meta.Rarity.(type) {
			case string:
				rarity = meta.Rarity.(string)
			case float64:
				rarity = fmt.Sprintf("%f", meta.Rarity.(float64))
			}
			metadata[position] = &PieceMetadata{Type: pieceTypeFriendlyName, Piece: pieceFriendlyName, Attributes: stats, Rarity: rarity, FriendlyName: meta.FriendlyName}
			pieces[position] = filePath
		}
	}
	return pieces, metadata
}
