package commons

import (
	"regexp"
	"strings"
)

const (
	killRgx   = `(\d{1,3}:\d{2}) Kill: (\d+) (\d+) (\d+): (.*?) killed (.*?) by (\w+)(?: (\w+))?`
	playerRgx = `(\d{1,3}:\d{2}) ClientUserinfoChanged: (\d+) n\\(.*?)\\t`
)

var KillRegex = regexp.MustCompile(killRgx)
var PlayerRegex = regexp.MustCompile(playerRgx)

const (
	initGameRef     = "InitGame:"
	shutdownGameRef = "ShutdownGame:"
	killRef         = "Kill:"
	playerRef       = "ClientUserinfoChanged:"
	separatorRef    = "-------"
	WorldRef        = "<world>"
)

type TypeOfProcess string

const (
	InitGame     TypeOfProcess = "InitGame"
	ShutdownGame TypeOfProcess = "ShutdownGame"
	Player       TypeOfProcess = "Player"
	Kill         TypeOfProcess = "Kill"
	Ignore       TypeOfProcess = "Ignore"
)

func ContainsString(arr []string, str string) bool {
	for _, item := range arr {
		if strings.Contains(item, str) {
			return true
		}
	}
	return false
}

func GetTypeOfProcess(line string) TypeOfProcess {
	switch {
	case strings.Contains(line, separatorRef):
		return Ignore
	case strings.Contains(line, initGameRef):
		return InitGame
	case strings.Contains(line, shutdownGameRef):
		return ShutdownGame
	case strings.Contains(line, playerRef):
		return Player
	case strings.Contains(line, killRef):
		return Kill
	default:
		return Ignore
	}
}
