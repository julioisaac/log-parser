package cmd

import (
	"log"
	"log-parser/config"
	"log-parser/processor"
	"log-parser/processor/quake/deaths-report"
	"log-parser/processor/quake/match-report"
	"os"
)

func Run() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	file, err := config.LoadFile(os.Getenv("LOG_FILE_PATH"))
	if err != nil {
		log.Fatalf("Failed to load log file: %v", err)
	}

	reports := []processor.ProcessCfg{
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

	p := processor.NewProcessor(file, reports)
	p.Start()

}
