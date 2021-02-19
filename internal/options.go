package runner

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type scanOptions struct {
	outputFile string
	silent     bool
	command    string
	filter     string
	interval   time.Duration
}

// ParseOptions parses the command line options for application
func ParseOptions() *scanOptions {
	options := &scanOptions{}

	flag.StringVar(&options.command, "c", "", "Command to execute and monitor")
	flag.StringVar(&options.outputFile, "o", "", "Output File Name")
	flag.StringVar(&options.filter, "f", "", "Regex Filter")
	flag.BoolVar(&options.silent, "s", false, "Verbose mode")
	flag.DurationVar(&options.interval, "i", time.Second*5, "Execute time interval")

	flag.Parse()
	if options.command == "" {
		fmt.Printf("Command not found \n\n")
		flag.Usage()
		os.Exit(1)
	}
	return options
}
