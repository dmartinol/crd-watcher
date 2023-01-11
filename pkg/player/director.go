package player

import (
	"time"

	v1 "github.com/dmartinol/crd-watcher/pkg/apis/requeststate/v1"
)

type Director struct {
	config  *Config
	watcher *Watcher
}

var job1 = "JOB1"
var job2 = "JOB2"

func NewDirectorForConfig(config *Config) *Director {
	director := &Director{config: config}

	return director
}

func (d *Director) Run() {
	ok := make(chan bool)
	defer close(ok)
	go d.directorSequence(ok)

	CreateNewRequestState(d.config, job1, d.config.namespace)
	<-ok
}

func (d *Director) directorSequence(ok chan<- bool) {
	d.watcher = StartWatcher(d.config, func(state v1.RequestState) {
		if state.Spec.State == STATE_REQUESTED {
			// NOP
		} else if state.Spec.State == STATE_STARTED {
			// NOP
		} else if state.Spec.State == STATE_COMPLETED {
			if state.Spec.Job == job1 {
				time.Sleep(2 * time.Second)
				UpdateRequestState(d.config, state, job2, STATE_REQUESTED)
			} else {
				time.Sleep(2 * time.Second)
				ok <- true
			}
		}
	})
}
