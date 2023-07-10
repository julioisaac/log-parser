package deathsReport

import (
	"reflect"
	"testing"
)

func TestDeathReportHandler_initGameFn(t *testing.T) {
	tests := []struct {
		name                string
		gameCounter         int
		expectedGameID      string
		expectedKillCounter int
		expectedDeaths      DeathMap
	}{
		{
			name:                "Valid initialization of game",
			gameCounter:         0,
			expectedGameID:      "game_1",
			expectedKillCounter: 0,
			expectedDeaths:      DeathMap{},
		},
		{
			name:                "Valid initialization of another game",
			gameCounter:         1,
			expectedGameID:      "game_2",
			expectedKillCounter: 0,
			expectedDeaths:      DeathMap{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DeathReportHandler{
				GameCounter: tt.gameCounter,
			}
			p.initGameFn("")

			if p.GameID != tt.expectedGameID {
				t.Errorf("Unexpected game ID. Got %s, want %s", p.GameID, tt.expectedGameID)
			}

			if p.KillCounter != tt.expectedKillCounter {
				t.Errorf("Unexpected kill counter. Got %d, want %d", p.KillCounter, tt.expectedKillCounter)
			}

			if len(p.Deaths) != len(tt.expectedDeaths) {
				t.Errorf("Unexpected number of deaths. Got %d, want %d", len(p.Deaths), len(tt.expectedDeaths))
			}

			for player, expectedTimes := range tt.expectedDeaths {
				if times, ok := p.Deaths[player]; ok {
					if times != expectedTimes {
						t.Errorf("Unexpected number of deaths for player %s. Got %d, want %d", player, times, expectedTimes)
					}
				} else {
					t.Errorf("Player %s not found in deaths map", player)
				}
			}
		})
	}
}

func TestDeathReportHandler_shutdownGameFn(t *testing.T) {
	tests := []struct {
		name           string
		gameID         string
		killCounter    int
		deaths         DeathMap
		expectedReport MatchWithDeathsCause
	}{
		{
			name:        "Valid shutdown game with non-zero kill counter",
			gameID:      "game_1",
			killCounter: 10,
			deaths:      DeathMap{"Player1": 5, "Player2": 3, "Player3": 2},
			expectedReport: MatchWithDeathsCause{
				TotalKills: 10,
				Deaths:     DeathMap{"Player1": 5, "Player2": 3, "Player3": 2}.Sort(),
			},
		},
		{
			name:        "Valid shutdown game with zero kill counter",
			gameID:      "game_2",
			killCounter: 0,
			deaths:      DeathMap{},
			expectedReport: MatchWithDeathsCause{
				TotalKills: 0,
				Deaths:     DeathMap{}.Sort(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DeathReportHandler{
				Report:      make(DeathCauseReport),
				GameID:      tt.gameID,
				KillCounter: tt.killCounter,
				Deaths:      tt.deaths,
			}
			p.shutdownGameFn("")

			report, ok := p.Report[tt.gameID]
			if !ok {
				t.Errorf("Game report for game ID %s not found", tt.gameID)
			}

			if report.TotalKills != tt.expectedReport.TotalKills {
				t.Errorf("Unexpected total kills. Got %d, want %d", report.TotalKills, tt.expectedReport.TotalKills)
			}

			if len(report.Deaths) != len(tt.expectedReport.Deaths) {
				t.Errorf("Unexpected number of deaths. Got %d, want %d", len(report.Deaths), len(tt.expectedReport.Deaths))
			}

			if !reflect.DeepEqual(report.Deaths, tt.expectedReport.Deaths) {
				t.Errorf("Unexpected death. Got %+v, want %+v", report.Deaths, tt.expectedReport.Deaths)
			}
		})
	}
}

func TestDeathReportHandler_deathProcess(t *testing.T) {
	tests := []struct {
		name          string
		inputLine     string
		givenDeaths   DeathMap
		expectedKills int
	}{
		{
			name:      "Valid death process with existing obit",
			inputLine: "20:55 Kill: 2 3 7: Player1 killed Player2 by MOD_ROCKET",
			givenDeaths: DeathMap{
				"MOD_ROCKET": 1,
			},
			expectedKills: 2,
		},
		{
			name:      "Valid death process with new obit",
			inputLine: "21:00 Kill: 4 3 7: Player3 killed Player2 by MOD_RAILGUN",
			givenDeaths: DeathMap{
				"MOD_ROCKET":  1,
				"MOD_RAILGUN": 1,
			},
			expectedKills: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DeathReportHandler{
				Deaths:      tt.givenDeaths,
				KillCounter: tt.givenDeaths.Sum(),
			}
			p.deathProcess(tt.inputLine)

			if len(p.Deaths) != len(tt.givenDeaths) {
				t.Errorf("Unexpected number of deaths. Got %d, want %d", len(p.Deaths), len(tt.givenDeaths))
			}

			for obit, expectedKills := range tt.givenDeaths {
				if kills, ok := p.Deaths[obit]; ok {
					if kills != expectedKills {
						t.Errorf("Unexpected number of kills for obit %s. Got %d, want %d", obit, kills, expectedKills)
					}
				} else {
					t.Errorf("Obit %s not found in deaths map", obit)
				}
			}

			if p.KillCounter != tt.expectedKills {
				t.Errorf("Unexpected number of kills counter. Got %d, want %d", p.KillCounter, tt.expectedKills)
			}
		})
	}
}
