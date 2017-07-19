package cmd

import (
	"os"

	"github.com/toxinu/kitsu/cleaner"
	"github.com/toxinu/kitsu/config"
	"github.com/toxinu/kitsu/watcher"
)

type Cmd struct {
	Config  *config.Config
	Cleaner *cleaner.Cleaner
	Watcher *watcher.Watcher
}

func New() *Cmd {
	c := &Cmd{}
	c.Config = config.New()
	c.Cleaner = cleaner.New(c.Config)
	c.Watcher = watcher.New(c.Config)
	return c
}

func (c Cmd) Run() {
	c.Config.Check()
	c.Cleaner.Start()

	exitCode := c.Watcher.Run()

	c.Cleaner.Stop()
	os.Exit(exitCode)
}
