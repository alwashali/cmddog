package runner

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type scanOptions struct {
	silent     bool
	command    string
	filter     string
	match      string
	configFile string
	outputFile string
	interval   time.Duration
}

// ParseOptions parses the command line options for application
func ParseOptions() *scanOptions {

	options := &scanOptions{}

	flag.StringVar(&options.outputFile, "o", "", "Output File Name")
	flag.StringVar(&options.configFile, "c", "", "Config File Name")
	flag.StringVar(&options.filter, "f", "", "Regex Filter, for more than one regex use the config file")
	flag.StringVar(&options.match, "m", "", "Matching Regex, for more than one regex use the config file")
	flag.BoolVar(&options.silent, "s", false, "Verbose mode")
	flag.DurationVar(&options.interval, "i", time.Second*5, "Execute time interval")

	flag.Parse()

	if flag.Arg(0) == "" && options.configFile == "" {
		fmt.Printf("Command not found \n%s Command Options\n\n", os.Args[0])
		flag.Usage()
		fmt.Printf("\n")
		os.Exit(1)
	} else if options.configFile != "" {
		options = parseConfig()
	}

	return options
}

func parseConfig() *scanOptions {
	return nil
}
