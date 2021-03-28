package runner

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// ScanOptions struct for cmd tool options
type ScanOptions struct {
	Command          string   `yaml:"command"`
	Args             []string `yaml:"args"`
	ReverseGrepRegex []string `yaml:"reversegrep"`
	GrepRegex        []string `yaml:"grep"`
	ConfigFile       string
	OutputFile       string `yaml:"output"`
	Silent           bool
	Interval         time.Duration
}

// ParseOptions parses the command line options for application
func ParseOptions() *ScanOptions {

	options := new(ScanOptions)

	flag.StringVar(&options.OutputFile, "o", "", "Output File Name, Default: Stdout")
	flag.StringVar(&options.ConfigFile, "c", "", "Config File Name")
	match := flag.String("g", "", "grep filter, skip everything except regex matches. For more than one regex use the config file")
	filter := flag.String("r", "", "reverse grep filter, print everything execpt regex matches. For more than one regex filter use the config file")
	flag.BoolVar(&options.Silent, "s", false, "Silent mode")
	flag.DurationVar(&options.Interval, "i", time.Second*3, "Execute time interval, e.g. 5s")

	flag.Parse()

	options.ReverseGrepRegex = append(options.ReverseGrepRegex, *filter)
	options.GrepRegex = append(options.GrepRegex, *match)
	options.Command = flag.Arg(0)

	if len(flag.Args()) > 1 {
		for i, v := range flag.Args() {
			if i == 0 { // skip flag.Args(0) -> name of the command itself
				continue
			}
			options.Args = append(options.Args, v) // append all supplied command line args
		}

	}

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
	configFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal("error parsing the configuration file ", err)
	}
	err = yaml.Unmarshal(byteValue, &options)
	if err != nil {
		log.Fatal("error parsing the configuration file ", err)
	}
	return &options
}
