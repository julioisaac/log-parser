package cmd

import (
	"log"
	"log-parser/config"
	"log-parser/processor"
	"log-parser/processor/quake/deaths-report"
	"log-parser/processor/quake/matches-report"
)

func Run() {

	file, err := config.LoadFile(config.GetString("LOG_FILE_PATH"))
	if err != nil {
		log.Fatalf("Failed to load log file: %v", err)
	}

	matchReportHandler := matchesReport.NewMatchReportHandler()
	deathReportHandler := deathsReport.NewDeathReportHandler()

	quakeReports := []processor.ProcessCfg{
		{
			ReportName:  "matches-report",
			ProcessLnFn: matchReportHandler.Process,
			Response:    matchReportHandler.GetReport,
		},
		{
			ReportName:  "death-report",
			ProcessLnFn: deathReportHandler.Process,
			Response:    deathReportHandler.GetReport,
		},
	}

	p := processor.NewFileProcessor("quake", file, quakeReports)
	p.Start()

}
