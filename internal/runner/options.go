package runner

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// ScanOptions struct for cmd tool options
type ScanOptions struct {
	Command          string   `json:"command"`
	ReverseGrepRegex []string `json:"reversegrep"`
	GrepRegex        []string `json:"grep"`
	ConfigFile       string
	OutputFile       string `json:"output"`
	Silent           bool
	Interval         time.Duration
}

// ParseOptions parses the command line options for application
func ParseOptions() *ScanOptions {

	options := new(ScanOptions)

	flag.StringVar(&options.OutputFile, "o", "", "Output File Name, Default: Stdout")
	flag.StringVar(&options.ConfigFile, "c", "", "Config File Name")
	filter := flag.String("g", "", "grep filter, for more than one regex use the config file")
	match := flag.String("r", "", "reverse grep filter, for more than one regex use the config file")
	flag.BoolVar(&options.Silent, "s", false, "Silent mode")
	flag.DurationVar(&options.Interval, "i", time.Second*3, "Execute time interval, e.g. 5s")

	flag.Parse()

	options.ReverseGrepRegex = append(options.ReverseGrepRegex, *filter)
	options.GrepRegex = append(options.GrepRegex, *match)
	options.Command = flag.Arg(0)

	if flag.Arg(0) == "" && options.ConfigFile == "" {
		fmt.Printf("Command not found \n\nTry:\n%s Options Command\n\n", os.Args[0])
		flag.Usage()
		fmt.Printf("\n")
		os.Exit(1)
	} else if options.ConfigFile != "" {
		return parseConfig(*options, options.ConfigFile)
	}

	return options
}

func parseConfig(options ScanOptions, filename string) *ScanOptions {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &options)
	if err != nil {
		log.Fatal("error parsing the configuration file ", err)
	}

	return &options
}
