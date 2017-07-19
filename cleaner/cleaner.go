package cleaner

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/toxinu/kitsu/config"
)

var fileNameRegexp = regexp.MustCompile(`(^.*)\.\d{14}\.\d{4,8}$`)

// Cleaner will delete file versions
type Cleaner struct {
	Config   *config.Config
	Disabled bool
	Interval int
	Exit     chan bool
}

// New return a cleaner
func New(cfg *config.Config) *Cleaner {
	clr := &Cleaner{
		Config:   cfg,
		Interval: cfg.CleanerInterval,
		Exit:     make(chan bool),
	}
	if cfg.CleanerInterval < 5 {
		log.Println("Cleaner interval won't go below 5 seconds.")
		clr.Interval = 5
	}
	if cfg.MaxVersions < 1 || cfg.CleanerInterval < 1 {
		log.Println("Cleaner disabled.")
		clr.Disabled = true
	}
	return clr
}

func (c *Cleaner) cleanMaxFile() error {
	var (
		pathToDelete       string
		filePath           string
		currentFilePath    string
		paths              []string
		currentFileCounter int
	)

	err := filepath.Walk(
		c.Config.Destination,
		func(path string, info os.FileInfo, err error) error {
			filePath = fileNameRegexp.ReplaceAllString(path, `$1`)
			if filePath != currentFilePath {
				currentFilePath = filePath
				currentFileCounter = 0
				paths = []string{}
			}

			currentFileCounter++
			paths = append(paths, path)

			if currentFileCounter > c.Config.MaxVersions {
				pathToDelete, paths = paths[0], paths[1:len(paths)]
				log.Printf("Deleting \"%s\"...\n", pathToDelete)
				err := os.Remove(pathToDelete)
				if err != nil {
					return err
				}
			}

			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		})
	return err
}

// Start run cleaner
func (c *Cleaner) Start() {
	if c.Disabled {
		return
	}

	timer := time.NewTimer(time.Duration(c.Interval) * time.Second)
	log.Printf("Cleaner will run every %d seconds with a maximum of %d file versions\n", c.Interval, c.Config.MaxVersions)

	go func() {
		for {
			select {
			case <-timer.C:
				log.Println("Cleaning file versions...")
				_ = c.cleanMaxFile()
				timer = time.NewTimer(
					time.Duration(c.Interval) * time.Second)
			case <-c.Exit:
				timer.Stop()
				return
			}

		}
	}()
}

// Stop shutdown cleaner
func (c *Cleaner) Stop() {
	if c.Disabled {
		return
	}
	c.Exit <- true
	close(c.Exit)
}
