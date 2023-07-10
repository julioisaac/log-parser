package processor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log-parser/logger"
	"os"
	"sync"
)

type ctxLine struct {
	ctx  context.Context
	line string
}

type ProcessCfg struct {
	SkipWriter  bool
	ReportName  string
	ProcessLnFn func(ctx context.Context, line string)
	Response    func() interface{}
}

type FileProcessor struct {
	name   string
	reader *bufio.Reader
	pc     []ProcessCfg
	n      int
	ch     []chan ctxLine
	wg     *sync.WaitGroup
}

func NewFileProcessor(name string, file *os.File, pCfg []ProcessCfg) *FileProcessor {
	numChannels := len(pCfg)
	channels := make([]chan ctxLine, numChannels)
	for i := 0; i < numChannels; i++ {
		channels[i] = make(chan ctxLine)
	}
	return &FileProcessor{
		name:   name,
		pc:     pCfg,
		n:      numChannels,
		reader: bufio.NewReader(file),
		ch:     channels,
		wg:     &sync.WaitGroup{},
	}
}

func (p *FileProcessor) Start() {
	logger.Log.Infof("[%s] file processor started", p.name)

	logger.Log.Infof("[%s] initializing %d goroutines", p.name, p.n)
	for i := 0; i < p.n; i++ {
		p.wg.Add(1)
		go p.processHandler(i, p.wg)
	}

	for {
		keepRunning := p.readAndProcessLn()
		if !keepRunning {
			p.Writer()
			p.Stop()
			break
		}
	}
	logger.Log.Infof("[%s] file processor stopped", p.name)

}

func (p *FileProcessor) processHandler(chIdx int, wg *sync.WaitGroup) {
	defer wg.Done()
	for ctxLn := range p.ch[chIdx] {
		func(ctx context.Context, ln string) {
			wg.Add(2)
			p.pc[chIdx].ProcessLnFn(ctx, ln)

		}(ctxLn.ctx, ctxLn.line)
	}
}

func (p *FileProcessor) readAndProcessLn() bool {
	ctx := context.Background()
	ln, err := p.reader.ReadString('\n')
	logger.Log.Tracef("[%s] reading line: %s", p.name, ln)

	if err != nil {
		if err == io.EOF || err == context.Canceled {
			return false
		}

		return true
	}
	for i := 0; i < p.n; i++ {
		p.ch[i] <- ctxLine{
			ctx:  ctx,
			line: ln,
		}
	}

	return true
}

func (p *FileProcessor) Stop() {
	for i := 0; i < p.n; i++ {
		close(p.ch[i])
	}
}

func (p *FileProcessor) Writer() {
	for i := 0; i < p.n; i++ {
		if !p.pc[i].SkipWriter {
			err := p.writeJson(p.pc[i].Response(), p.pc[i].ReportName)
			if err != nil {
				logger.Log.WithError(err).Error("error writing json")
			}
		}
	}
}

const jsonExt = "json"

func (p *FileProcessor) writeJson(data interface{}, filename string) error {
	logger.Log.Debugf("[%s] writing file '%s'", p.name, filename)
	file, err := os.Create(fmt.Sprintf("%s.%s", filename, jsonExt))
	if err != nil {
		return err
	}
	defer file.Close()

	encodedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(encodedData)
	if err != nil {
		return err
	}
	logger.Log.Debugf("%s", encodedData)

	return nil

}
