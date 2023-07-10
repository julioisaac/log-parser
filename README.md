```
                   _
                  | |
  __ _ _   _  __ _| | _____
 / _` | | | |/ _` | |/ / _ \
| (_| | |_| | (_| |   <  __/
 \__, |\__,_|\__,_|_|\_\___|
    | |
    |_|
                                                                   
log parser
```
Although it makes reference to Quake, the script can read any log file and deliver it to a function that will process it as needed.

# Approach
The approach used for reading and processing the file was to implement go routines that receive the lines as they are read and process them concurrently using the functions supplied as inputs.
The go routines in the coding example are creating match-report and death-report.

## Prerequisites (run on Docker)

Before you start using it, you will need to have [Docker](https://www.docker.com/) installed.

## How to Run

* **Install dependencies**

```bash
$ make install
```

* **Running**

```bash
$ make run
```

## How to Test

* **Check tests and linter**

```bash
$ make checks
```

#### Others useful Make commands for local environment

| Command          | Description                                                |
|------------------|------------------------------------------------------------|
| `deps`           | downloads all the dependencies needed to build the project |
| `test`           | runs all the unit tests                                    |
| `docker/build`   | build the docker image                                     |
| `docker/start`   | runs the project by docker                                 | 
| `docker/stop`    | stop docker container                                      |


## Implementation

```processor.ProcessCfg```
```go
type ProcessCfg struct {
	SkipWriter  bool
	ReportName  string
	ProcessLnFn func(ctx context.Context, line string)
	Response    func() interface{}
}
```
The structure was created to make the pass of required parameters easier.

```ProcessLnFn```
```go 
func(ctx context.Context, line string)
```
This is the signature that must be followed by the process function in order to receive the lines to be processed.


```processor.NewFileProcessor```
```go 
func NewFileProcessor(name string, file *os.File, pCfg []ProcessCfg) *FileProcessor
```
The method was designed to read any file and transfer it to another function that would process each line received and write the output.

```go
// Example
...
    anotherReports := []processor.ProcessCfg{
        {
            ReportName:  "another-report",
            ProcessLnFn: anotherReport.Process,
            Response:    anotherReport.GetReport,
        },
    }

    p := processor.NewFileProcessor("another-file", file, anotherReports)
    p.Start()
...
```

## Built With

+ golang 1.20

## About

