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
	options    *scanOptions
	om         *cmdlog.OutputMonitor
	lastupdate time.Time
}

// New creates a new runner
func New(cmdOptions *scanOptions) (*Runner, error) {
	args := []string{}
	runner := &Runner{
		options:    cmdOptions,
		om:         cmdlog.New(cmdOptions.command, args),
		lastupdate: time.Now(),
	}
	if cmdOptions.filter != "" {
		runner.om.SetFilter(cmdOptions.filter)
	}
	return runner, nil
}

// Execute command and listen for changes
func (r *Runner) Execute() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go r.om.Run(r.options.interval)
	if r.options.silent == false {
		go r.om.PrintOutput()
	}
	wg.Wait()
}

// Output Print date written to output slice
func (r *Runner) Output() {
	for {
		fmt.Println(r.lastupdate.After(*r.om.LastChanged()))
		if r.lastupdate.After(*r.om.LastChanged()) {
			r.om.PrintOutput()
		}
		time.Sleep(r.options.interval)
	}
}

func (r *Runner) write() {
	f, err := os.OpenFile(r.options.outputFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	if r.lastupdate.After(*r.om.LastChanged()) {
		if _, err := f.WriteString("appended"); err != nil {
			log.Println(err)
		}
	}

}
