package deathsReport

import "sort"

type Death struct {
	Cause string `json:"cause"`
	Count int    `json:"count"`
}

type Deaths []Death

type MatchWithDeathsCause struct {
	TotalKills int    `json:"total_kills"`
	Deaths     Deaths `json:"deaths"`
}

type DeathCauseReport map[string]MatchWithDeathsCause

type DeathMap map[string]int

func (km DeathMap) Sort() []Death {
	sorted := []Death{}
	keys := make([]string, 0, len(km))
	for k := range km {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return km[keys[i]] > km[keys[j]]
	})

	for _, key := range keys {
		sorted = append(sorted, Death{
			Cause: key,
			Count: km[key],
		})
	}

	return sorted
}

func (km DeathMap) Sum() int {
	sum := 0
	for _, value := range km {
		sum += value
	}
	return sum
}
