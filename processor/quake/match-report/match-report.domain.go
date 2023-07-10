package matchReport

import "sort"

type MatchKill struct {
	Ranking       int    `json:"ranking"`
	Player        string `json:"Player"`
	NumberOfTimes int    `json:"numberOfTimes"`
}

type Players []string
type MatchReport map[string]MatchWithKills

type MatchWithKills struct {
	TotalKills int         `json:"total_kills"`
	Players    Players     `json:"players"`
	Kills      []MatchKill `json:"kills"`
}

type KillsMap map[string]int

func (km KillsMap) Sort() []MatchKill {
	sorted := []MatchKill{}
	keys := make([]string, 0, len(km))
	for k := range km {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return km[keys[i]] > km[keys[j]]
	})

	for i, key := range keys {
		sorted = append(sorted, MatchKill{
			Ranking:       i + 1,
			Player:        key,
			NumberOfTimes: km[key],
		})
	}

	return sorted
}
