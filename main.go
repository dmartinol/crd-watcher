package main

import (
	"github.com/dmartinol/crd-watcher/pkg/player"
)

func main() {
	config := player.NewConfig()

	if config.RunAsDirector() {
		director := player.NewDirectorForConfig(config)
		director.Run()
	} else if config.RunAsWorker() {
		worker := player.NewWorkerForConfig(config)
		worker.Run()
	}
}
