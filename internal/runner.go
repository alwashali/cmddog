package runner

import (
	"sync"
	"time"

	cmdlog "github.com/alwashali/cmdlog/cmdlog"
)

// Runner is a client for running the enumeration process.
type Runner struct {
	options *scanOptions
	om      *cmdlog.Cmdlog
}

// New creates a new runner
func New(cmdOptions *scanOptions) (*Runner, error) {
	args := []string{}
	runner := &Runner{
		options: cmdOptions,
		om:      cmdlog.New(cmdOptions.command, args),
	}
	if cmdOptions.filter != "" {
		runner.om.SetFilter(cmdOptions.filter)
	}
	return runner, nil
}

// Execute Periodically run command
func (r *Runner) Execute() {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go r.om.Run(r.options.interval)
	if r.options.silent != true {
		go func() {
			for {
				r.om.PrintNew()
				time.Sleep(r.options.interval)
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
