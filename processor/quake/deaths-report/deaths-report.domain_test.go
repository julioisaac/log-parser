package deathsReport

import (
	"testing"
)

func TestDeathMap_Sort(t *testing.T) {
	tests := []struct {
		name string
		dm   DeathMap
		want []Death
	}{
		{
			name: "MultipleDeathCause",
			dm: DeathMap{
				"MOD_TELEFRAG":     25,
				"MOD_TRIGGER_HURT": 17,
				"MOD_RAILGUN":      37,
			},
			want: []Death{
				{Cause: "MOD_RAILGUN", Count: 37},
				{Cause: "MOD_TELEFRAG", Count: 25},
				{Cause: "MOD_TRIGGER_HURT", Count: 17},
			},
		},
		{
			name: "JustOnCause",
			dm: DeathMap{
				"MOD_TRIGGER_HURT": 17,
			},
			want: []Death{
				{Cause: "MOD_TRIGGER_HURT", Count: 17},
			},
		},
		{
			name: "NoDeaths",
			dm:   DeathMap{},
			want: []Death{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dm.Sort()

			if len(result) != len(tt.want) {
				t.Errorf("Expected %d deaths, got %d deaths", len(tt.want), len(result))
				return
			}

			for i, d := range result {
				if d.Cause != tt.want[i].Cause ||
					d.Count != tt.want[i].Count {
					t.Errorf("Mismatch at position %d. Expected %+v, got %+v", i, tt.want[i], d)
				}
			}
		})
	}
}
