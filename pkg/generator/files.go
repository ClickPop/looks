package generator

import (
	"bytes"
	"fmt"
	"os"

	"github.com/clickpop/looks/internal/utils"
	conf "github.com/clickpop/looks/pkg/config"
)

func loadFiles(config *conf.Config) ([]*bytes.Reader, Metadata, error) {
	fileNames := config.Settings.PieceOrder
	var files []*bytes.Reader
	var metadata Metadata
	for i := 0; i < len(fileNames); i++ {
		file := fileNames[i]
		pieceTypeFriendlyName := config.Attributes[file].FriendlyName
		if pieceTypeFriendlyName == "" {
			pieceTypeFriendlyName = utils.TransformName(file)
		}
		piece, meta := handleRarity(config.Attributes[file].Pieces, config.Settings.Rarity, config.Output)

		if piece != "nil" {
			pieceFriendlyName := config.Attributes[file].Pieces[piece].FriendlyName
			if pieceFriendlyName == "" {
				pieceFriendlyName = utils.TransformName(piece)
			}
			filename := fmt.Sprintf(config.Input.Local.Filename, file, piece)
			data, err := os.ReadFile(fmt.Sprintf("%s/%s", config.Input.Local.Pathname, filename))
			if err != nil {
				return nil, Metadata{}, err
			}
			files = append(files, bytes.NewReader(data))
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
			metadata.PieceMeta = append(metadata.PieceMeta, PieceMetadata{Type: pieceTypeFriendlyName, Piece: pieceFriendlyName, Attributes: stats, Rarity: meta.Rarity, FriendlyName: meta.FriendlyName})
		}
	}

	return files, metadata, nil
}
