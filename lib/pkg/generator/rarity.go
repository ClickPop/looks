package generator

import (
	"math/rand"
	"time"

	"github.com/clickpop/looks/internal/utils"
	"github.com/clickpop/looks/pkg/config"
)

func getRarityDenominator(r config.ConfigRarity) int {
	denominator := 0
	for _, v := range r.Chances {
		denominator += int(v)
	}

	return denominator
}

func getRarityMinimum(r config.ConfigRarity, minRarity string) int {
	minimum := 0
	hasMinKey := false
	minKey := 0

	if minRarity != "" {
		for k, v := range r.Order {
			if v == minRarity {
				minKey = k - 1
				if minKey >= 0 {
					hasMinKey = true
				}
			}
		}
		if hasMinKey {
			for i := 0; i <= minKey; i++ {
				minimum += r.Chances[r.Order[i]]
			}
		}
	}

	return minimum
}

func getRarityLevel(r config.ConfigRarity, minRarity string) []string {
	denominator := getRarityDenominator(r)
	minimum := getRarityMinimum(r, minRarity)

	seed := time.Now().Unix() + int64(time.Now().Nanosecond())
	rand.Seed(seed)

	random := rand.Intn(denominator-minimum) + minimum

	rarity := []string{r.Order[0]}

	currentThreshold := 0

	for i, v := range r.Order {
		if random >= currentThreshold {
			if i != 0 {
				rarity = append(rarity, v)
			}
		}
		currentThreshold += int(r.Chances[v])
	}
	rarity = utils.ReverseStringSlice(rarity)
	return rarity
}

func tagCheck(piece config.PieceAttribute, tags map[string]bool, tagsConfig config.TagConfigSettings) bool {
	for _, tag := range piece.Tags {
		for _, excludeTag := range tagsConfig.Exclusive[tag] {
			if _, ok := tags[excludeTag]; ok {
				return false
			}
		}

		for _, includeTag := range tagsConfig.Inclusive[tag] {
			if _, ok := tags[includeTag]; !ok {
				return false
			}
		}
	}

	return true
}

func handleRarity(pieceTypes map[string]config.PieceAttribute, rarityData config.ConfigRarity, outputOpt config.OutputObject, tagConfig config.TagConfigSettings, tags map[string]bool) (string, config.PieceAttribute) {
	rarityLevel := getRarityLevel(rarityData, outputOpt.MinimumRarity)
	var possiblePieces []string
	for _, v := range rarityLevel {
		for key, piece := range pieceTypes {
			if v == piece.Rarity && tagCheck(piece, tags, tagConfig) {
				possiblePieces = append(possiblePieces, key)
			}
		}
		if len(possiblePieces) > 0 {
			break
		}
	}

	seed := time.Now().Unix() + int64(time.Now().Nanosecond())
	rand.Seed(seed)
	random := rand.Intn(len(possiblePieces))
	choice := possiblePieces[random]
	return choice, pieceTypes[choice]
}
