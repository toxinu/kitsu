package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// Config is kitsu configuration
type Config struct {
	Source          string
	Destination     string
	MaxVersions     int
	CleanerInterval int
	Excludes        []regexp.Regexp
}

type stringslice []string

func (s *stringslice) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *stringslice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// usage print usage details
func usage() {
	fmt.Printf("Usage: %s [OPTIONS] <source> <destination>\n", os.Args[0])
	flag.PrintDefaults()
}

// New return kitsu configuration
func New() *Config {
	var (
		err             error
		re              *regexp.Regexp
		config          *Config
		excludePatterns stringslice
	)
	config = &Config{}

	flag.Usage = usage
	flag.Var(&excludePatterns, "exclude", "Exclude pattern as regexp (repeatable).")
	maxVersions := flag.Int("max-versions", 20, "Maximum file versions to keep (will delete olders).")
	cleanerInterval := flag.Int("cleaner-interval", 60, "Cleaner interval in seconds. (minimum 5)")
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	config.Source = flag.Arg(0)
	config.Destination = flag.Arg(1)
	config.MaxVersions = *maxVersions
	config.CleanerInterval = *cleanerInterval

	if config.Source == "" || config.Destination == "" {
		log.Fatal("Error: kitsu <source> <destination>")
	}

	// Resolve absolute path for source
	config.Source, err = filepath.Abs(config.Source)
	if err != nil {
		log.Fatal(err)
	}

	// Resolve absolute path for destination
	config.Destination, err = filepath.Abs(config.Destination)
	if err != nil {
		log.Fatal(err)
	}

	// Check if source exists
	_, err = os.Stat(config.Source)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(10)
	}

	// Compile exclude patterns
	for _, pattern := range excludePatterns {
		log.Printf("Exclude: %s\n", pattern)
		re, err = regexp.Compile(pattern)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(20)
		}
		config.Excludes = append(config.Excludes, *re)
	}

	return config
}

func (c Config) Check() {
	if _, err := os.Stat(c.Destination); os.IsNotExist(err) {
		log.Printf("Creating destination directory\n")
		sourceDirectory, _ := os.Stat(c.Source)
		os.MkdirAll(c.Destination, sourceDirectory.Mode().Perm())
	}
}
