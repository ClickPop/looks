package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/clickpop/looks/internal/utils"
	"github.com/clickpop/looks/pkg/config"
)

// func getRarityDenominator(r config.ConfigRarity) int {
// 	denominator := 0
// 	for _, v := range r.Chances {
// 		denominator += int(v)
// 	}
//
// 	return denominator
// }
//
// func getRarityMinimum(r config.ConfigRarity, minRarity string) int {
// 	minimum := 0
// 	hasMinKey := false
// 	minKey := 0
//
// 	if minRarity != "" {
// 		for k, v := range r.Order {
// 			if v == minRarity {
// 				minKey = k - 1
// 				if minKey >= 0 {
// 					hasMinKey = true
// 				}
// 			}
// 		}
// 		if hasMinKey {
// 			for i := 0; i <= minKey; i++ {
// 				minimum += r.Chances[r.Order[i]]
// 			}
// 		}
// 	}
//
// 	return minimum
// }
//
// func getRarityLevel(r config.ConfigRarity, minRarity string) []string {
// 	denominator := getRarityDenominator(r)
// 	minimum := getRarityMinimum(r, minRarity)
//
// 	seed := time.Now().Unix() + int64(time.Now().Nanosecond())
// 	rand.Seed(seed)
//
// 	random := rand.Intn(denominator-minimum) + minimum
//
// 	rarity := []string{r.Order[0]}
//
// 	currentThreshold := 0
//
// 	for i, v := range r.Order {
// 		if random >= currentThreshold {
// 			if i != 0 {
// 				rarity = append(rarity, v)
// 			}
// 		}
// 		currentThreshold += int(r.Chances[v])
// 	}
// 	rarity = utils.ReverseStringSlice(rarity)
// 	return rarity
// }

func tagCheck(pieceTags []string, tags map[string]bool, tagsConfig config.TagConfigSettings) bool {
	for _, tag := range pieceTags {
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

func filterVariants(pieceVariants []string, variant string, configVariants config.TagConfigSettings) []string {
  acceptableVariants := make(map[string]bool)
  filtered := make([]string, 0)
  include := configVariants.Inclusive[variant]
  exclude := configVariants.Exclusive[variant]
  for _, v := range utils.Filter(pieceVariants, func(val string) bool {
    return (len(include) < 1 || utils.Contains(include, val)) && (len(exclude) < 1 || !utils.Contains(exclude, val))
  }) {
    acceptableVariants[v] = true
  }

  for variant := range acceptableVariants {
    filtered = append(filtered, variant)
  }
  return filtered
}

func handleRarity(pieceTypes map[string]config.PieceAttribute, config *config.Config, tags map[string]bool, variant string, position int) (string, config.PieceAttribute, string, string) {
  variants := make(map[string][]string)
  rarityData := config.Settings.Rarity
  possiblePieces := make(map[string]float64)
  sum := float64(0)
  for key, piece := range pieceTypes {
    if tagCheck(piece.Tags, tags, config.Settings.Tags) {
      if len(piece.Variants) < 1 {
        switch piece.Rarity.(type) {
        case string:
          rarity := rarityData.Chances[piece.Rarity.(string)]
          sum += rarity
          possiblePieces[key] = rarity
        case float64:
          rarity := piece.Rarity.(float64)
          sum += rarity
          possiblePieces[key] = rarity
        }
      } else {
        filteredVariants := filterVariants(piece.Variants, variant, config.Settings.Variants)
        if (len(filteredVariants) > 0) {
        switch piece.Rarity.(type) {
        case string:
          rarity := rarityData.Chances[piece.Rarity.(string)]
          sum += rarity
          possiblePieces[key] = rarity
        case float64:
          rarity := piece.Rarity.(float64)
          sum += rarity
          possiblePieces[key] = rarity
        }
          variants[key] = filteredVariants
        }
      }
    }
  }

	seed := time.Now().Unix() + int64(time.Now().Nanosecond())
	rand.Seed(seed)
  random := rand.Float64() * sum
  var choice string
  randSum := float64(0)
  for piece, rarity := range possiblePieces {
    randSum += rarity
    if random < randSum {
      choice = piece
      break;
    }
  }
  piece := pieceTypes[choice]
  friendlyName := piece.FriendlyName
  if friendlyName == "" {
    friendlyName = utils.TransformName(choice)
  }
  variantCount := len(variants[choice])
  if variantCount > 0 && position > 0 {
    randomVariant := rand.Intn(variantCount)
    variant := variants[choice][randomVariant]
    choice = fmt.Sprintf("%s_%s", choice, variant)
    friendlyName = fmt.Sprintf("%s %s", variant, piece.FriendlyName)
  } else if variantCount == 1 && position == 0 {
    variant = piece.Variants[0]
  }
	return choice, pieceTypes[choice], variant, friendlyName
}
