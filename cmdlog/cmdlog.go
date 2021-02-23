package commandwatch

import (
	"bufio"
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
	lastPrintIndex   int
}

// New creates new object
func New(cmd string, args []string) *Cmdlog {
	return &Cmdlog{
		command: cmd,
		args:    args,
	}
}

// Insert to unique output values
func (e *Cmdlog) insertNew(output string) {

	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		found := false

		// range keyword doesn't work on empty slice, first item appended directly
		if len(e.results) == 0 {
			e.results = append(e.results, line)
			continue
		}
		for _, item := range e.results {
			if line == item {
				found = true
			}
		}
		if found == false {
			e.results = append(e.results, line)

		}

	}

}

// Filter output using regex
func (e *Cmdlog) reverseGrep(output string) string {
	// loop through all filters and clean the string matching the regex filter

	var str []string

	str = append(str, output)
	// iterate str and each time string is cleaned using regex put it in a new place inside the slice
	for i, filter := range e.reverseGrepRegex {
		re := regexp.MustCompile(filter)
		outstr := re.ReplaceAllString(str[i], "")
		scanner := bufio.NewScanner(strings.NewReader(outstr))
		strBuilder := strings.Builder{}
		for scanner.Scan() {
			if scanner.Text() != "" {
				strBuilder.WriteString(scanner.Text() + "\n")
			}
		}
		str = append(str, strBuilder.String())

	}
	i := len(str)
	return str[i-1]

}

func (e *Cmdlog) grep(output string) string {
	// loop through all filters and clean the string matching the regex filter

	var str []string
	// iterate str and each time string is cleaned using regex put it in a new place inside the slice
	for _, filter := range e.grepRegex {
		re := regexp.MustCompile(filter)
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

// Run continuesly execute the command and write the new output values to output buffer
func (e *Cmdlog) Run(sleepTime time.Duration) {

	for {
		cmd := exec.Command(e.command, e.args...)

		// run the command and output
		output, err := cmd.Output()
		if err != nil {
			log.Fatalf("Output failed with %s\n", err)
		}
		out := string(output)

		if e.reverseGrepRegex != nil {
			out = e.reverseGrep(out)
		}

		if e.grepRegex != nil {
			out = e.grep(out)
		}

		e.insertNew(out)

		// The sleep value depends on how frequent your command outputs the results
		time.Sleep(sleepTime)
	}

}

// Results returns the results slice from the specified index, pass 0 for returning entire slice
func (e *Cmdlog) Results(i int) []string {
	if i >= 0 && len(e.results) > 0 && i < len(e.results) {
		return e.results[i:]
	}
	return e.results
}

// ResultsSize returns number of elements in results slice
func (e *Cmdlog) ResultsSize() int {
	return len(e.results)
}
