package matchReport

import (
	"testing"
)

func TestKillsMap_Sort(t *testing.T) {
	tests := []struct {
		name string
		km   KillsMap
		want []MatchKill
	}{
		{
			name: "MultiplePlayers",
			km: KillsMap{
				"Isgalamido":   5,
				"Zeh":          3,
				"Dono da bola": 7,
			},
			want: []MatchKill{
				{Ranking: 1, Player: "Dono da bola", NumberOfTimes: 7},
				{Ranking: 2, Player: "Isgalamido", NumberOfTimes: 5},
				{Ranking: 3, Player: "Zeh", NumberOfTimes: 3},
			},
		},
		{
			name: "SinglePlayer",
			km: KillsMap{
				"Zeh": -1,
			},
			want: []MatchKill{
				{Ranking: 1, Player: "Zeh", NumberOfTimes: -1},
			},
		},
		{
			name: "NoPlayers",
			km:   KillsMap{},
			want: []MatchKill{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.km.Sort()

			if len(result) != len(tt.want) {
				t.Errorf("Expected %d kills, got %d kills", len(tt.want), len(result))
				return
			}

			for i, k := range result {
				if k.Ranking != tt.want[i].Ranking ||
					k.Player != tt.want[i].Player ||
					k.NumberOfTimes != tt.want[i].NumberOfTimes {
					t.Errorf("Mismatch at position %d. Expected %+v, got %+v", i, tt.want[i], k)
				}
			}
		})
	}
}
