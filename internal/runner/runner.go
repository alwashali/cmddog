package runner

import (
	"sync"
	"time"

	cmdlog "github.com/alwashali/cmdlog/cmdlog"
)

// Runner is a client for running the enumeration process.
type Runner struct {
	options *ScanOptions
	om      *cmdlog.Cmdlog
}

// New creates a new runner
func New(cmdOptions *ScanOptions) (*Runner, error) {
	args := []string{}

	runner := &Runner{
		options: cmdOptions,
		// create cmdlog instance pasing command and args
		om: cmdlog.New(cmdOptions.Command, args),
	}
	if cmdOptions.ReverseGrepRegex[0] != "" {
		runner.om.SetReverseGrepRegex(cmdOptions.ReverseGrepRegex[0])
	}
	if cmdOptions.GrepRegex[0] != "" {
		runner.om.SetGrepRegex(cmdOptions.GrepRegex[0])
	}
	return runner, nil
}

// Execute Periodically run command
func (r *Runner) Execute() {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go r.om.Run(r.options.Interval)
	if r.options.silent != true {
		go func() {
			for {
				r.om.PrintNew()
				time.Sleep(r.options.Interval)
			}
		}()
	}
	wg.Wait()
}

// func (r *Runner) write() {
// 	f, err := os.OpenFile(r.options.outputFile,
// 		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer f.Close()

// 	if r.lastupdate.After(*r.om.LastChanged()) {
// 		if _, err := f.WriteString("appended"); err != nil {
// 			log.Println(err)
// 		}
// 	}

// }
