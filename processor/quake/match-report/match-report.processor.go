package matchReport

import (
	"context"
	"fmt"
	commons2 "log-parser/processor/quake/commons"
)

var matchReport MatchReport
var players Players
var kills KillsMap
var gameId string
var gameCounter, killCounter int

var execute = map[commons2.TypeOfProcess]func(string){
	commons2.InitGame:     initGameFn,
	commons2.ShutdownGame: shutdownGameFn,
	commons2.Player:       playerProcess,
	commons2.Kill:         killProcess,
}

func Process(_ context.Context, line string) {
	typeOfProcess := commons2.GetTypeOfProcess(line)
	if commons2.Ignore == typeOfProcess {
		return
	}

	execute[typeOfProcess](line)
}

func GetReport() interface{} {
	return matchReport
}

func initGameFn(_ string) {
	if matchReport == nil {
		matchReport = MatchReport{}
		gameCounter = 0
	}
	kills = map[string]int{}
	players = Players{}
	gameCounter++
	gameId = fmt.Sprintf("game_%d", gameCounter)
}

func shutdownGameFn(_ string) {
	matchReport[gameId] = MatchWithKills{
		TotalKills: killCounter,
		Players:    players,
		Kills:      kills.Sort(),
	}
	killCounter = 0
}

func playerProcess(line string) {
	matches := commons2.PlayerRegex.FindStringSubmatch(line)
	if !commons2.ContainsString(players, matches[commons2.PlayerName]) {
		players = append(players, matches[commons2.PlayerName])
	}
}

func killProcess(line string) {
	killCounter++
	matches := commons2.KillRegex.FindStringSubmatch(line)
	if commons2.WorldRef == matches[commons2.KillerName] {
		val, ok := kills[matches[commons2.VictimName]]
		if ok {
			kills[matches[commons2.VictimName]] = val - 1
		} else {
			kills[matches[commons2.VictimName]] = -1
		}
		return
	}

	val, ok := kills[matches[commons2.KillerName]]
	if ok {
		kills[matches[commons2.KillerName]] = val + 1
	} else {
		kills[matches[commons2.KillerName]] = 1
	}
}
