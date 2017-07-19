package watcher

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/toxinu/kitsu/config"
)

type Watcher struct {
	Config *config.Config
}

func New(cfg *config.Config) *Watcher {
	w := &Watcher{Config: cfg}
	return w
}

func (w *Watcher) Run() int {
	return w.watch()
}

func (w *Watcher) copyFile(destination string, source string) (error, int) {
	if _, err := os.Stat(destination); err == nil {
		return errors.New("file already exists, something wrong"), 0
	}
	in, err := os.Open(source)
	if err != nil {
		return err, 0
	}
	defer in.Close()
	out, err := os.Create(destination)
	if err != nil {
		return err, 0
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err, 0
	}
	return cerr, 0
}

func (w *Watcher) handleEvent(path string, source string, rootDestination string, excludes []regexp.Regexp) (error, int) {
	var (
		fi             os.FileInfo
		relativePath   string
		dirPath        string
		destination    string
		dirDestination string
		err            error
		exitCode       int
	)
	relativePath, err = filepath.Rel(source, path)
	if err != nil {
		return err, 0
	}

	for _, pattern := range excludes {
		if pattern.Match([]byte(relativePath)) == true {
			log.Printf("Ignoring %s", relativePath)
			return nil, 0
		}
	}

	destination = filepath.Join(rootDestination, relativePath)
	dirDestination = filepath.Dir(destination)
	dirPath = filepath.Dir(path)

	if _, err = os.Stat(dirDestination); os.IsNotExist(err) {
		fmt.Printf("Creating directory: %s\n", dirDestination)
		sourceDir, _ := os.Stat(dirPath)
		os.MkdirAll(dirDestination, sourceDir.Mode().Perm())
	}

	fi, err = os.Stat(path)
	if err != nil {
		return err, 0
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		log.Printf("Creating directory %s\n", destination)
		os.Mkdir(destination, mode.Perm())
		return nil, 0
	case mode.IsRegular():
		destination += "." + time.Now().UTC().Format("20060102150405.99999999")
		log.Printf("Copying %s\n", path)
		err, exitCode = w.copyFile(destination, path)
		if err != nil {
			log.Printf("Error while copying file: %s\n", err)
		}
		return nil, exitCode
	}
	return nil, 0
}

func (w *Watcher) watch() int {
	fileSystemChannel := make(chan notify.EventInfo, 1)

	keyboardInterruptChannel := make(chan os.Signal, 1)
	signal.Notify(keyboardInterruptChannel, os.Interrupt)

	log.Printf("Watching %s\n", w.Config.Source)
	if err := notify.Watch(
		w.Config.Source+"...", fileSystemChannel,
		notify.Create, notify.Write); err != nil {
		log.Fatal(err)
	}

	defer notify.Stop(fileSystemChannel)

	for {
		select {
		case eventInfo := <-fileSystemChannel:
			log.Printf("Event %s %s\n", eventInfo.Event().String(), eventInfo.Path())
			err, exitCode := w.handleEvent(
				eventInfo.Path(), w.Config.Source,
				w.Config.Destination, w.Config.Excludes)
			if err != nil {
				log.Printf("Error while copying file: %s\n", err)
				if exitCode > 0 {
					return exitCode
				}
			}
		case _ = <-keyboardInterruptChannel:
			log.Println("Received Interrupt signal.")
			return 2
		}
	}
}
