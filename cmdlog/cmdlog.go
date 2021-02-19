package commandwatch

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// OutputMonitor options
type OutputMonitor struct {
	command      string
	args         []string
	output       []string
	regexfilters []string
	lastchanged  time.Time
}

// New creates new object
func New(cmd string, args []string) *OutputMonitor {
	return &OutputMonitor{
		command: cmd,
		args:    args,
	}
}

// Insert to unique output values
func (e *OutputMonitor) insert(line string) {
	if line != "" {
		e.output = append(e.output, line)
	}
}

// Filter output using regex
func (e *OutputMonitor) filter(output string) string {
	// loop through all filters and clean the string matching the regex filter
	if len(e.regexfilters) > 0 {
		var str []string
		str = append(str, output)
		// iterate str and each time string is cleaned using regex put it in a new place inside the slice
		for i, filter := range e.regexfilters {
			re := regexp.MustCompile(filter)
			str = append(str, re.ReplaceAllString(str[i], ""))

		}
		i := len(str)
		return str[i-1]
	}
	return output
}

// SetFilter for cleaning the output before store
// Use with caution -> your regex should match only unwanted match in the output
func (e *OutputMonitor) SetFilter(filter string) {
	e.regexfilters = append(e.regexfilters, filter)
}

//check if there is any new lines in the output
func (e *OutputMonitor) check(output string) {
	for _, item := range e.output {
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			if item != scanner.Text() {
				e.insert(scanner.Text())
				e.lastchanged = time.Now()

			}
		}

	}
}

// Run continuesly execute the command and write the new output values to output buffer
func (e *OutputMonitor) Run(sleepTime time.Duration) {

	for {
		cmd := exec.Command(e.command, e.args...)
		// run the command and output
		output, err := cmd.Output()
		if err != nil {
			log.Fatalf("Output failed with %s\n", err)
		}
		out := string(output)
		out = e.filter(out)
		e.check(out)
		// The sleep value depends on how frequent your command outputs the results
		time.Sleep(sleepTime)
	}

}

// Output returns the output slice
func (e *OutputMonitor) Output() []string {
	return e.output
}

// PrintOutput returns the output written
func (e *OutputMonitor) PrintOutput() {
	for line := range e.output {
		fmt.Println(line)
	}
}

// LastChanged returns last time a new value is added to output slice
func (e *OutputMonitor) LastChanged() *time.Time {
	return &e.lastchanged
}
