package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/clickpop/looks/internal/utils"
	conf "github.com/clickpop/looks/pkg/config"
)

func buildDescription(c *conf.Config, meta OpenSeaMeta) (string, string) {
	switch {
	case c.Descriptions.SimpleFragments != nil && len(c.Descriptions.SimpleFragments) > 0:
		return buildSimpleDescription(c, meta)
	case c.Descriptions.StatFragments != nil:
		return buildStatDescription(c, meta)
	}
	return "", ""
}

func buildSimpleDescription(c *conf.Config, meta OpenSeaMeta) (string, string) {
	fragments := make([]string, 0)
	fragmentMap := make(map[string]bool, c.Descriptions.FragmentCount)
	for len(fragmentMap) < c.Descriptions.FragmentCount {
		rand.Seed(time.Now().Unix() + int64(time.Now().Nanosecond()))
		fragment := c.Descriptions.SimpleFragments[rand.Intn(len(c.Descriptions.SimpleFragments))]
		fragmentMap[fragment] = true
	}

	for k := range fragmentMap {
		fragments = append(fragments, k)
	}

	return fmt.Sprintf(c.Descriptions.Template, utils.OxfordJoin(fragments)), ""
}

func buildStatDescription(c *conf.Config, meta OpenSeaMeta) (string, string) {
	stats := make(map[string]int)
	namesToKeys := make(map[string]string)
	namesToKeys["fallback"] = "fallback"
	for k, v := range c.Settings.Stats {
		name := v.Name
		if name == "" {
			name = utils.TransformName(k)
		}
		stats[name] = 0
		namesToKeys[name] = k
	}

	for _, v := range meta.Attributes {
		if _, isStat := stats[v.TraitType]; isStat {
			stats[v.TraitType] += v.Value.(int)
		}
	}
	
	primaryStat := getPrimaryStat(stats, c.Descriptions.FallbackPrimaryStat)
	randomDescriptor := getRandomDescriptor(c.Descriptions.StatFragments[namesToKeys[primaryStat]].Descriptors)
	randomHobbies := getRandomHobbies(c.Descriptions.StatFragments[namesToKeys[primaryStat]].Hobbies, c.Descriptions.FragmentCount)
	currentType := c.Descriptions.StatFragments[namesToKeys[primaryStat]].Name

	return fmt.Sprintf(c.Descriptions.Template, currentType, randomDescriptor, randomHobbies), currentType
}

func getRandomDescriptor(descriptors []string) string {
	rand.Seed(time.Now().Unix() + int64(time.Now().Nanosecond()))
	return descriptors[rand.Intn(len(descriptors))]
}

func getRandomHobbies(hobbies []string, n int) string {
	rand.Seed(time.Now().Unix() + int64(time.Now().Nanosecond()))
	var randomHobbies []string

	if len(hobbies) < n {
		n = len(hobbies)
	}

	for len(randomHobbies) < n {
		tempHobby := hobbies[rand.Intn(len(hobbies))]
		if !utils.Contains(randomHobbies, tempHobby) {
			randomHobbies = append(randomHobbies, tempHobby)
		}
	}

	return utils.OxfordJoin(randomHobbies)
}
