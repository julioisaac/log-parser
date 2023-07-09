package deathsReport

import (
	"context"
	"fmt"
	commons2 "log-parser/processor/quake/commons"
)

var deathCauseReport DeathCauseReport
var deaths DeathMap
var gameId string
var gameCounter, killCounter int

var execute = map[commons2.TypeOfProcess]func(string){
	commons2.InitGame:     initGameFn,
	commons2.ShutdownGame: shutdownGameFn,
	commons2.Player:       playerProcess,
	commons2.Kill:         deathProcess,
}

func GetReport() interface{} {
	return deathCauseReport
}

func Process(ctx context.Context, line string) {
	typeOfProcess := commons2.GetTypeOfProcess(line)
	if commons2.Ignore == typeOfProcess {
		return
	}

	execute[typeOfProcess](line)
}

func initGameFn(_ string) {
	if deathCauseReport == nil {
		deathCauseReport = DeathCauseReport{}
		gameCounter = 0
	}
	deaths = map[string]int{}
	gameCounter++
	gameId = fmt.Sprintf("game_%d", gameCounter)
}

func shutdownGameFn(_ string) {
	deathCauseReport[gameId] = MatchWithDeathsCause{
		TotalKills: killCounter,
		Deaths:     deaths.Sort(),
	}
	killCounter = 0
}

func playerProcess(_ string) {
}

func deathProcess(line string) {
	killCounter++
	matches := commons2.KillRegex.FindStringSubmatch(line)

	val, ok := deaths[matches[commons2.Obit]]
	if ok {
		deaths[matches[commons2.Obit]] = val + 1
	} else {
		deaths[matches[commons2.Obit]] = 1
	}
}
