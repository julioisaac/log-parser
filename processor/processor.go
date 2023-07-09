package processor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type ctxLine struct {
	ctx  context.Context
	line string
}

type ProcessLnFn func(ctx context.Context, value string)
type Response func() interface{}

type ProcessCfg struct {
	SkipWriter  bool
	ReportName  string
	ProcessLnFn ProcessLnFn
	Response    Response
}

type Processor struct {
	reader *bufio.Reader
	pc     []ProcessCfg
	ch     []chan ctxLine
	wg     *sync.WaitGroup
}

func NewProcessor(file *os.File, pCfg []ProcessCfg) *Processor {
	numChannels := len(pCfg)
	channels := make([]chan ctxLine, numChannels)
	for i := 0; i < numChannels; i++ {
		channels[i] = make(chan ctxLine)
	}
	return &Processor{
		pc:     pCfg,
		reader: bufio.NewReader(file),
		ch:     channels,
		wg:     &sync.WaitGroup{},
	}
}

func (p *Processor) Start() {

	for i := 0; i < len(p.ch); i++ {
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
}

func (p *Processor) processHandler(chIdx int, wg *sync.WaitGroup) {
	defer wg.Done()
	for ctxLn := range p.ch[chIdx] {
		func(ctx context.Context, ln string) {
			wg.Add(2)
			p.pc[chIdx].ProcessLnFn(ctx, ln)

		}(ctxLn.ctx, ctxLn.line)
	}
}

func (p *Processor) readAndProcessLn() bool {
	ctx := context.Background()
	ln, err := p.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF || err == context.Canceled {
			return false
		}

		return true
	}
	for i := 0; i < len(p.ch); i++ {
		p.ch[i] <- ctxLine{
			ctx:  ctx,
			line: ln,
		}
	}

	return true
}

func (p *Processor) Stop() {
	for i := 0; i < len(p.ch); i++ {
		close(p.ch[i])
	}
}

func (p *Processor) Writer() {
	for i := 0; i < len(p.ch); i++ {
		if !p.pc[i].SkipWriter {
			p.reportWriter(p.pc[i].Response(), p.pc[i].ReportName)
		}
	}
}

func (p *Processor) reportWriter(data interface{}, filename string) error {
	file, err := os.Create(filename + ".json")
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encodedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}

	_, err = file.Write(encodedData)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %v", err)
	}

	return nil

}
