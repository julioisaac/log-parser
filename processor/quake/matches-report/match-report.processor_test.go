package matchesReport

import (
	"log-parser/processor/quake/commons"
	"testing"
)

func TestMatchReportHandler_initGameFn(t *testing.T) {
	tests := []struct {
		name                string
		gameCounter         int
		expectedGameID      string
		expectedKillCounter int
		expectedKills       KillsMap
		expectedPlayers     Players
	}{
		{
			name:                "Valid initialization of game",
			gameCounter:         0,
			expectedGameID:      "game_1",
			expectedKillCounter: 0,
			expectedKills:       KillsMap{},
			expectedPlayers:     Players{},
		},
		{
			name:                "Valid initialization of another game",
			gameCounter:         1,
			expectedGameID:      "game_2",
			expectedKillCounter: 0,
			expectedKills:       KillsMap{},
			expectedPlayers:     Players{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MatchReportHandler{
				GameCounter: tt.gameCounter,
			}
			p.initGameFn("")

			if p.GameID != tt.expectedGameID {
				t.Errorf("Unexpected game ID. Got %s, want %s", p.GameID, tt.expectedGameID)
			}

			if p.KillCounter != tt.expectedKillCounter {
				t.Errorf("Unexpected kill counter. Got %d, want %d", p.KillCounter, tt.expectedKillCounter)
			}

			if len(p.Kills) != len(tt.expectedKills) {
				t.Errorf("Unexpected number of kills. Got %d, want %d", len(p.Kills), len(tt.expectedKills))
			}

			for player, expectedKills := range tt.expectedKills {
				if kills, ok := p.Kills[player]; ok {
					if kills != expectedKills {
						t.Errorf("Unexpected number of kills for player %s. Got %d, want %d", player, kills, expectedKills)
					}
				} else {
					t.Errorf("Player %s not found in kills map", player)
				}
			}

			if len(p.Players) != len(tt.expectedPlayers) {
				t.Errorf("Unexpected number of players. Got %d, want %d", len(p.Players), len(tt.expectedPlayers))
			}

			for _, player := range tt.expectedPlayers {
				if !commons.ContainsString(p.Players, player) {
					t.Errorf("Player %s not found in players list", player)
				}
			}
		})
	}
}

func TestMatchReportHandler_shutdownGameFn(t *testing.T) {
	tests := []struct {
		name                string
		gameID              string
		killCounter         int
		players             Players
		kills               KillsMap
		expectedReport      MatchWithKills
		expectedKillCounter int
	}{
		{
			name:        "Valid shutdown game with non-zero kill counter",
			gameID:      "game_1",
			killCounter: 10,
			players:     Players{"Isgalamido", "Oootsimo", "Zeh"},
			kills:       KillsMap{"Zeh": 2, "Isgalamido": 3, "Oootsimo": 5},
			expectedReport: MatchWithKills{
				TotalKills: 10,
				Players:    Players{"Isgalamido", "Oootsimo", "Zeh"},
				Kills:      KillsMap{"Oootsimo": 5, "Isgalamido": 3, "Zeh": 2}.Sort(),
			},
			expectedKillCounter: 0,
		},
		{
			name:        "Valid shutdown game with zero kill counter",
			gameID:      "game_2",
			killCounter: 0,
			players:     Players{"Zeh", "Oootsimo"},
			kills:       KillsMap{},
			expectedReport: MatchWithKills{
				TotalKills: 0,
				Players:    Players{"Zeh", "Oootsimo"},
				Kills:      KillsMap{}.Sort(),
			},
			expectedKillCounter: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MatchReportHandler{
				Report:      make(MatchReport),
				GameID:      tt.gameID,
				KillCounter: tt.killCounter,
				Players:     tt.players,
				Kills:       tt.kills,
			}
			p.shutdownGameFn("")

			report, ok := p.Report[tt.gameID]
			if !ok {
				t.Errorf("Game report for game ID %s not found", tt.gameID)
			}

			if report.TotalKills != tt.expectedReport.TotalKills {
				t.Errorf("Unexpected total kills. Got %d, want %d", report.TotalKills, tt.expectedReport.TotalKills)
			}

			if len(report.Players) != len(tt.expectedReport.Players) {
				t.Errorf("Unexpected number of players. Got %d, want %d", len(report.Players), len(tt.expectedReport.Players))
			}

			for i, player := range report.Players {
				if player != tt.expectedReport.Players[i] {
					t.Errorf("Unexpected player. Got %s, want %s", player, tt.expectedReport.Players[i])
				}
			}

			if len(report.Kills) != len(tt.expectedReport.Kills) {
				t.Errorf("Unexpected number of kills. Got %d, want %d", len(report.Kills), len(tt.expectedReport.Kills))
			}

			for i, kill := range report.Kills {
				expectedKill := tt.expectedReport.Kills[i]
				if kill.Player != expectedKill.Player || kill.NumberOfTimes != expectedKill.NumberOfTimes {
					t.Errorf("Unexpected kill. Got %+v, want %+v", kill, expectedKill)
				}
			}

			if p.KillCounter != tt.expectedKillCounter {
				t.Errorf("Unexpected kill counter. Got %d, want %d", p.KillCounter, tt.expectedKillCounter)
			}
		})
	}
}

func TestMatchReportHandler_playerProcess(t *testing.T) {
	tests := []struct {
		name            string
		inputLine       string
		existingPlayers []string
		expectedPlayers []string
	}{
		{
			name:            "Valid player process with new player",
			inputLine:       "0:04 ClientUserinfoChanged: 6 n\\Zeh\\t\\0\\model\\sarge/default\\hmodel\\sarge/default\\g_redteam\\\\g_blueteam\\\\c1\\1\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
			existingPlayers: []string{"Isgalamido", "Oootsimo"},
			expectedPlayers: []string{"Isgalamido", "Oootsimo", "Zeh"},
		},
		{
			name:            "Valid player process with existing player",
			inputLine:       "0:05 ClientUserinfoChanged: 6 n\\Oootsimo\\t\\0\\model\\sarge/default\\hmodel\\sarge/default\\g_redteam\\\\g_blueteam\\\\c1\\1\\c2\\5\\hc\\100\\w\\0\\l\\0\\tt\\0\\tl\\0",
			existingPlayers: []string{"Isgalamido", "Oootsimo", "Zeh"},
			expectedPlayers: []string{"Isgalamido", "Oootsimo", "Zeh"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MatchReportHandler{
				Players: tt.existingPlayers,
			}
			p.playerProcess(tt.inputLine)

			if len(p.Players) != len(tt.expectedPlayers) {
				t.Errorf("Unexpected number of players. Got %d, want %d", len(p.Players), len(tt.expectedPlayers))
			}

			for _, expectedPlayer := range tt.expectedPlayers {
				if !commons.ContainsString(p.Players, expectedPlayer) {
					t.Errorf("Player %s not found in players list", expectedPlayer)
				}
			}
		})
	}
}

func TestMatchReportHandler_killProcess(t *testing.T) {
	tests := []struct {
		name           string
		inputLine      string
		givenKills     KillsMap
		expectedPlayer KillsMap
		expectedKilled int
	}{
		{
			name:      "Valid kill process with existing killer and victim",
			inputLine: "20:55 Kill: 2 3 7: Player1 killed Player2 by MOD_ROCKET",
			givenKills: KillsMap{
				"Player1": 1,
			},
			expectedKilled: 2,
		},
		{
			name:      "Valid kill process with new killer and existing victim",
			inputLine: "21:00 Kill: 4 3 7: Player3 killed Player2 by MOD_RAILGUN",
			givenKills: KillsMap{
				"Player1": 1,
				"Player3": 1,
			},
			expectedKilled: 3,
		},
		{
			name:      "Valid kill process with existing killer and new victim",
			inputLine: "21:05 Kill: 2 5 7: Player1 killed Player4 by MOD_SHOTGUN",
			givenKills: KillsMap{
				"Player1": 2,
				"Player3": 1,
			},
			expectedKilled: 4,
		},
		{
			name:      "Valid kill process with killer as world",
			inputLine: "21:10 Kill: 1022 5 22: <world> killed Player4 by MOD_TRIGGER_HURT",
			givenKills: KillsMap{
				"Player1": 2,
				"Player3": 1,
			},
			expectedKilled: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MatchReportHandler{
				Kills:       tt.givenKills,
				KillCounter: tt.givenKills.Sum(),
			}
			p.killProcess(tt.inputLine)

			if len(p.Kills) != len(tt.givenKills) {
				t.Errorf("Unexpected number of kills. Got %d, want %d", len(p.Kills), len(tt.givenKills))
			}

			for player, expectedKills := range tt.givenKills {
				if kills, ok := p.Kills[player]; ok {
					if kills != expectedKills {
						t.Errorf("Unexpected number of kills for player %s. Got %d, want %d", player, kills, expectedKills)
					}
				} else {
					t.Errorf("Player %s not found in kills map", player)
				}
			}

			if p.KillCounter != tt.expectedKilled {
				t.Errorf("Unexpected number of kills counter. Got %d, want %d", p.KillCounter, tt.expectedKilled)
			}
		})
	}
}
