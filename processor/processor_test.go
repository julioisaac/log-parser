package processor

import (
	"encoding/json"
	"io"
	"log"
	"log-parser/config"
	deathsReport "log-parser/processor/quake/deaths-report"
	matchReport "log-parser/processor/quake/match-report"
	"os"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	type args struct {
		name string
		file *os.File
		pCfg []ProcessCfg
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Given the test1 log should return json teste1 as expected with just the Player Zeh in Kills field",
			args: args{
				name: "quake",
				file: func() *os.File {
					givenLogTest1, err := config.LoadFile("../tests/test1-given.log")
					if err != nil {
						log.Fatalf("Failed to load test1 log file: %v", err)
					}
					return givenLogTest1
				}(),
				pCfg: func() []ProcessCfg {
					return []ProcessCfg{
						{
							ReportName:  "match-report",
							SkipWriter:  true,
							ProcessLnFn: matchReport.Process,
							Response:    matchReport.GetReport,
						},
						{
							ReportName:  "deaths-report",
							SkipWriter:  true,
							ProcessLnFn: deathsReport.Process,
							Response:    deathsReport.GetReport,
						},
					}
				}(),
			},
			want: func() []string {
				expectedReports := []string{}
				expectedTest1, err := config.LoadFile("../tests/test1-mr-expected.json")
				if err != nil {
					log.Fatalf("Failed to load test1 json file: %v", err)
				}
				if test1, err := io.ReadAll(expectedTest1); err == nil {
					expectedReports = append(expectedReports, string(test1))
				}
				expectedTest2, err := config.LoadFile("../tests/test1-dr-expected.json")
				if err != nil {
					log.Fatalf("Failed to load test1 json file: %v", err)
				}
				if test2, err := io.ReadAll(expectedTest2); err == nil {
					expectedReports = append(expectedReports, string(test2))
				}
				return expectedReports
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewFileProcessor(tt.args.name, tt.args.file, tt.args.pCfg)
			p.Start()

			for i, cfg := range tt.args.pCfg {
				gotJson, err := json.MarshalIndent(cfg.Response(), "", "  ")
				if (err != nil) != tt.wantErr {
					t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if string(gotJson) != tt.want[i] {
					t.Errorf("Process() got = %v, want %v", string(gotJson), tt.want[i])
				}
			}
		})
	}
}
