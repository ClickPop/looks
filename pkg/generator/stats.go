package generator

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