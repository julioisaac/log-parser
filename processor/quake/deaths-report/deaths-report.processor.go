package deathsReport

import (
	"context"
	"fmt"
	"log-parser/processor/quake/commons"
)

type DeathReportHandler struct {
	Report      DeathCauseReport
	Deaths      DeathMap
	GameID      string
	GameCounter int
	KillCounter int
}

var dispatchTable = map[commons.TypeOfProcess]func(*DeathReportHandler, string){
	commons.InitGame:     (*DeathReportHandler).initGameFn,
	commons.ShutdownGame: (*DeathReportHandler).shutdownGameFn,
	commons.Kill:         (*DeathReportHandler).deathProcess,
}

func NewDeathReportHandler() *DeathReportHandler {
	return &DeathReportHandler{
		Report:      DeathCauseReport{},
		Deaths:      DeathMap{},
		GameCounter: 0,
	}
}

func (p *DeathReportHandler) GetReport() interface{} {
	return p.Report
}

func (p *DeathReportHandler) Process(_ context.Context, line string) {
	typeOfProcess := commons.GetTypeOfProcess(line)
	if typeOfProcess == commons.Ignore || typeOfProcess == commons.Player {
		return
	}

	dispatchTable[typeOfProcess](p, line)
}

func (p *DeathReportHandler) initGameFn(_ string) {
	p.GameCounter++
	p.GameID = fmt.Sprintf("game_%d", p.GameCounter)
	p.KillCounter = 0
	p.Deaths = DeathMap{}
}

func (p *DeathReportHandler) shutdownGameFn(_ string) {
	p.Report[p.GameID] = MatchWithDeathsCause{
		TotalKills: p.KillCounter,
		Deaths:     p.Deaths.Sort(),
	}
}

func (p *DeathReportHandler) deathProcess(line string) {
	p.KillCounter++
	matches := commons.KillRegex.FindStringSubmatch(line)

	val, ok := p.Deaths[matches[commons.Obit]]
	if ok {
		p.Deaths[matches[commons.Obit]] = val + 1
	} else {
		p.Deaths[matches[commons.Obit]] = 1
	}
}
