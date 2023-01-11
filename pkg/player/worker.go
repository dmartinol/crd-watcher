package player

import (
	"time"

	v1 "github.com/dmartinol/crd-watcher/pkg/apis/requeststate/v1"
)

type Worker struct {
	config  *Config
	watcher *Watcher
}

func NewWorkerForConfig(config *Config) *Worker {
	worker := &Worker{config: config}

	return worker
}

func (w *Worker) Run() {
	w.watcher = StartWatcher(w.config, func(state v1.RequestState) {
		if state.Spec.State == STATE_REQUESTED {
			time.Sleep(2 * time.Second)
			UpdateRequestState(w.config, state, "", STATE_STARTED)
		} else if state.Spec.State == STATE_STARTED {
			time.Sleep(2 * time.Second)
			UpdateRequestState(w.config, state, "", STATE_COMPLETED)
		} else if state.Spec.State == STATE_COMPLETED {
			// NOP
		}
	})
}
