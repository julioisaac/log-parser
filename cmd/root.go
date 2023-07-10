package cmd

import (
	"log"
	"log-parser/config"
	"log-parser/processor"
	"log-parser/processor/quake/deaths-report"
	"log-parser/processor/quake/match-report"
)

func Run() {

	file, err := config.LoadFile(config.GetString("LOG_FILE_PATH"))
	if err != nil {
		log.Fatalf("Failed to load log file: %v", err)
	}

	quakeReports := []processor.ProcessCfg{
		{
			ReportName:  "match-report",
			ProcessLnFn: matchReport.Process,
			Response:    matchReport.GetReport,
		},
		{
			ReportName:  "death-report",
			ProcessLnFn: deathsReport.Process,
			Response:    deathsReport.GetReport,
		},
	}

	p := processor.NewFileProcessor("quake", file, quakeReports)
	p.Start()

}
