package matchesReport

import (
	"context"
	"fmt"
	"log-parser/processor/quake/commons"
)

type MatchReportHandler struct {
	Report      MatchReport
	Players     Players
	Kills       KillsMap
	GameID      string
	GameCounter int
	KillCounter int
}

var dispatchTable = map[commons.TypeOfProcess]func(*MatchReportHandler, string){
	commons.InitGame:     (*MatchReportHandler).initGameFn,
	commons.ShutdownGame: (*MatchReportHandler).shutdownGameFn,
	commons.Player:       (*MatchReportHandler).playerProcess,
	commons.Kill:         (*MatchReportHandler).killProcess,
}

func NewMatchReportHandler() *MatchReportHandler {
	return &MatchReportHandler{
		Report:      MatchReport{},
		Kills:       KillsMap{},
		GameCounter: 0,
	}
}

func (p *MatchReportHandler) GetReport() interface{} {
	return p.Report
}

func (p *MatchReportHandler) Process(_ context.Context, line string) {
	typeOfProcess := commons.GetTypeOfProcess(line)
	if typeOfProcess == commons.Ignore {
		return
	}

	dispatchTable[typeOfProcess](p, line)
}

func (p *MatchReportHandler) initGameFn(_ string) {
	p.GameCounter++
	p.GameID = fmt.Sprintf("game_%d", p.GameCounter)
	p.KillCounter = 0
	p.Kills = KillsMap{}
	p.Players = Players{}
}

func (p *MatchReportHandler) shutdownGameFn(_ string) {
	p.Report[p.GameID] = MatchWithKills{
		TotalKills: p.KillCounter,
		Players:    p.Players,
		Kills:      p.Kills.Sort(),
	}
	p.KillCounter = 0
}

func (p *MatchReportHandler) playerProcess(line string) {
	matches := commons.PlayerRegex.FindStringSubmatch(line)
	playerName := matches[commons.PlayerName]
	if !commons.ContainsString(p.Players, playerName) {
		p.Players = append(p.Players, playerName)
	}
}

func (p *MatchReportHandler) killProcess(line string) {
	p.KillCounter++
	matches := commons.KillRegex.FindStringSubmatch(line)
	killerName := matches[commons.KillerName]
	victimName := matches[commons.VictimName]

	if killerName == commons.WorldRef {
		if val, ok := p.Kills[victimName]; ok {
			p.Kills[victimName] = val - 1
		} else {
			p.Kills[victimName] = -1
		}
		return
	}

	if val, ok := p.Kills[killerName]; ok {
		p.Kills[killerName] = val + 1
	} else {
		p.Kills[killerName] = 1
	}
}
