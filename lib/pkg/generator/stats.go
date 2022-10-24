package generator

import "github.com/clickpop/looks/pkg/config"

func getPrimaryStat(stats map[string]int, fallbackPrimaryStat string) string {
	max := -(int(^uint(0) >> 1)) - 1
	primaryStat := fallbackPrimaryStat

	for stat, v := range stats {
		if v > max {
			primaryStat = stat
			max = v
		} else if v == max {
			primaryStat = fallbackPrimaryStat
		}
	}
	return primaryStat
}

func buildStats(config *config.Config) map[string]int {
	stats := make(map[string]int)
	for k := range config.Settings.Stats {
		stats[k] = 0
	}
	return stats
}
