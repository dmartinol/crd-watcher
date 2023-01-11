package player

import (
	"context"
	"encoding/json"
	"log"
	"time"

	v1 "github.com/dmartinol/crd-watcher/pkg/apis/requeststate/v1"
	informers "github.com/dmartinol/crd-watcher/pkg/client/informers/externalversions"
	"github.com/dmartinol/crd-watcher/pkg/signals"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

var STATE_REQUESTED = "REQUESTED"
var STATE_STARTED = "STARTED"
var STATE_COMPLETED = "COMPLETED"

func ReadAllRequestStates(config *Config) []v1.RequestState {
	states, err := config.stateclient.ComV1().RequestStates("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	output := []v1.RequestState{}
	for _, state := range states.Items {
		klog.Infof("Found request state %s in namespace %s", state.Name, state.ObjectMeta.Namespace)
		output = append(output, state)
	}
	return output
}

func CreateNewRequestState(config *Config, job string, namespace string) *v1.RequestState {
	uuid := uuid.New().String()
	state := &v1.RequestState{
		ObjectMeta: metav1.ObjectMeta{
			Name:      uuid,
			Namespace: namespace,
		},
		Spec: v1.RequestStateSpec{
			RequestUid: uuid,
			Job:        job,
			State:      STATE_REQUESTED,
		},
	}
	PrettyPrint("Creating request state to %s", state)
	state, err := config.stateclient.ComV1().RequestStates(namespace).Create(context.TODO(), state, metav1.CreateOptions{})
	if err != nil {
		klog.Fatal(err)
	}
	return state
}

func UpdateRequestState(config *Config, state v1.RequestState, job string, newstate string) *v1.RequestState {
	stateCopy := state.DeepCopy()
	stateCopy.Status.History = append(stateCopy.Status.History, v1.RequestStateHistory{
		Job:       stateCopy.Spec.Job,
		State:     stateCopy.Spec.State,
		Timestamp: time.Now().String(),
	})
	if job != "" {
		stateCopy.Spec.Job = job
	}
	stateCopy.Spec.State = newstate
	PrettyPrint("Updating request state to %s", stateCopy)
	output, err := config.stateclient.ComV1().RequestStates(stateCopy.ObjectMeta.Namespace).Update(context.TODO(), stateCopy, metav1.UpdateOptions{})
	if err != nil {
		klog.Fatal(err)
	}
	return output
}

func StartWatcher(config *Config, fn func(state v1.RequestState)) *Watcher {
	stopCh := signals.SetupSignalHandler()

	stateInformerFactory := informers.NewSharedInformerFactory(config.stateclient, time.Second*30)

	watcher := NewWatcher(config.stateclient,
		stateInformerFactory.Com().V1().RequestStates(), fn)
	klog.Info("Starting controller")
	stateInformerFactory.Start(stopCh)

	if err := watcher.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
	return watcher
}

func PrettyPrint(text string, state *v1.RequestState) {
	marshaled, err := json.MarshalIndent(state, "", "   ")
	if err != nil {
		log.Fatalf("Marshaling error: %s", err)
	}
	klog.Infof(text, string(marshaled))
}
