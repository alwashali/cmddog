package runner

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	cmdlog "github.com/alwashali/cmdlog/cmdlog"
)

// Runner is a client for running the enumeration process.
type Runner struct {
	options *ScanOptions
	cl      *cmdlog.Cmdlog
}

// New creates a new runner
func New(cmdOptions *ScanOptions) (*Runner, error) {
	args := []string{}

	runner := &Runner{
		options: cmdOptions,
		// create cmdlog instance pasing command and args
		cl: cmdlog.New(cmdOptions.Command, args),
	}
	if cmdOptions.ReverseGrepRegex[0] != "" {
		runner.cl.SetReverseGrepRegex(cmdOptions.ReverseGrepRegex[0])
	}
	if cmdOptions.GrepRegex[0] != "" {
		runner.cl.SetGrepRegex(cmdOptions.GrepRegex[0])
	}
	return runner, nil
}

// Execute Periodically run command
func (r *Runner) Execute() {
	wg := new(sync.WaitGroup)
	wg.Add(3)

	// wg 1
	go r.cl.Run(r.options.Interval)

	if r.options.Silent != true {
		//wg 2
		go func() {
			i := 0
			size := 0
			for {
				size = r.cl.ResultsSize()
				if size > i {
					for _, value := range r.cl.Results(i) {
						fmt.Println(value)
					}
					i = size
				}

				time.Sleep(r.options.Interval)

			}

		}()
	}

	if r.options.OutputFile != "" {
		//wg 3
		go func() {

			filename := r.options.OutputFile
			stat, err := os.Stat(filename)
			if err == nil {
				if stat.IsDir() {
					log.Panicf("%s is a directory", filename)
				}
			}

			f, err := os.OpenFile(filename,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()

			i := 0
			size := 0

			for {
				size = r.cl.ResultsSize()

				if size > i {
					for _, line := range r.cl.Results(i) {
						if _, err := f.WriteString(line + "\n"); err != nil {
							log.Println(err)
						}
					}
					i = size
				}

				time.Sleep(r.options.Interval)

			}

		}()
	}

	wg.Wait()
}
