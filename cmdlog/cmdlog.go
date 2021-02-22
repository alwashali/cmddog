package commandwatch

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Cmdlog options
type Cmdlog struct {
	command          string
	args             []string
	results          []string
	grepRegex        []string
	reverseGrepRegex []string
	lastPrintIndex   int64
}

// New creates new object
func New(cmd string, args []string) *Cmdlog {
	return &Cmdlog{
		command:        cmd,
		args:           args,
		lastPrintIndex: -1,
	}
}

// Insert to unique output values
func (e *Cmdlog) insert(line string) {
	if line != "" {
		e.results = append(e.results, line)
	}
}

// Filter output using regex
func (e *Cmdlog) filter(output string) string {
	// loop through all filters and clean the string matching the regex filter

	var str []string
	str = append(str, output)
	// iterate str and each time string is cleaned using regex put it in a new place inside the slice
	for i, filter := range e.reverseGrepRegex {
		re := regexp.MustCompile(filter)
		str = append(str, re.ReplaceAllString(str[i], ""))

	}
	i := len(str)
	return str[i-1]

}

func (e *Cmdlog) match(output string) string {
	// loop through all filters and clean the string matching the regex filter

	var str []string
	// iterate str and each time string is cleaned using regex put it in a new place inside the slice
	fmt.Println(e.grepRegex)
	for _, mfilter := range e.grepRegex {
		re := regexp.MustCompile(mfilter)
		for _, m := range re.FindAllString(output, -1) {
			str = append(str, m)
		}

	}
	rstring := strings.Builder{}
	for _, s := range str {
		rstring.WriteString(s)
		rstring.WriteString("\n")
	}

	return rstring.String()

}

// SetReverseGrepRegex for cleaning the output before store
// Use with caution -> your regex should match only unwanted text in the output
func (e *Cmdlog) SetReverseGrepRegex(filter string) {
	e.reverseGrepRegex = append(e.reverseGrepRegex, filter)
}

// SetGrepRegex set regex pattern
func (e *Cmdlog) SetGrepRegex(match string) {
	e.grepRegex = append(e.grepRegex, match)
}

//check if there is any new lines(value) in the output
func (e *Cmdlog) check(output string) {

	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		found := false
		for _, item := range e.results {
			if line == item {
				found = true
			}
		}
		if found == false {
			e.insert(line)
		}

	}

}

// Run continuesly execute the command and write the new output values to output buffer
func (e *Cmdlog) Run(sleepTime time.Duration) {
	e.results = append(e.results, " ")
	for {
		cmd := exec.Command(e.command, e.args...)

		// run the command and output
		output, err := cmd.Output()
		if err != nil {
			log.Fatalf("Output failed with %s\n", err)
		}
		out := string(output)
		if e.reverseGrepRegex != nil {
			out = e.filter(out)
		}
		if e.grepRegex != nil {
			out = e.match(out)
		}
		e.check(out)
		// The sleep value depends on how frequent your command outputs the results
		time.Sleep(sleepTime)
	}

}

// Results returns the results slice
func (e *Cmdlog) Results() []string {
	return e.results
}

// PrintNew prints the new added values
func (e *Cmdlog) PrintNew() {
	for i, line := range e.results {
		if i > int(e.lastPrintIndex) {
			fmt.Println(line)
			e.lastPrintIndex = int64(i)
		}
	}
}
