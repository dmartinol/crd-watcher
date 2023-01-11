package player

import (
	"fmt"
	"time"

	v1 "github.com/dmartinol/crd-watcher/pkg/apis/requeststate/v1"
	clientset "github.com/dmartinol/crd-watcher/pkg/client/clientset/versioned"
	requestscheme "github.com/dmartinol/crd-watcher/pkg/client/clientset/versioned/scheme"
	informers "github.com/dmartinol/crd-watcher/pkg/client/informers/externalversions/requeststate/v1"
	listers "github.com/dmartinol/crd-watcher/pkg/client/listers/requeststate/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	SuccessSynced         = "Synced"
	ErrResourceExists     = "ErrResourceExists"
	MessageResourceSynced = "RequestState synced successfully"
)

// Watcher is the controller implementation for RequestState resources
type Watcher struct {
	stateclient clientset.Interface
	stateLister listers.RequestStateLister
	stateSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	// recorder record.EventRecorder

	fn func(state v1.RequestState)
}

func NewWatcher(
	stateclient clientset.Interface,
	stateInformer informers.RequestStateInformer,
	fn func(state v1.RequestState)) *Watcher {

	utilruntime.Must(scheme.AddToScheme(requestscheme.Scheme))
	// klog.V(4).Info("Creating event broadcaster")
	// eventBroadcaster := record.NewBroadcaster()
	// eventBroadcaster.StartStructuredLogging(0)
	// eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	// recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	watcher := &Watcher{
		stateclient: stateclient,
		stateLister: stateInformer.Lister(),
		stateSynced: stateInformer.Informer().HasSynced,
		workqueue:   workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "RequestStates"),
		// recorder:    recorder,
		fn: fn,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when RequestState resources change
	stateInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: watcher.enqueueRequestState,
		UpdateFunc: func(old, new interface{}) {
			watcher.enqueueRequestState(new)
		},
	})

	return watcher
}

func (c *Watcher) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting RequestState watcher")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.stateSynced); !ok {
		return fmt.Errorf("Failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Foo resources
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func (c *Watcher) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Watcher) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("Expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("Error syncing '%s': %s, requeuing", key, err.Error())
		}

		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Watcher) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Invalid resource key: %s", key))
		return nil
	}

	state, err := c.stateLister.RequestStates(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("requeststate '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	klog.Infof("Received %s %s for UID %s, job %s and state %s", state.Kind, state.Name, state.Spec.RequestUid, state.Spec.Job, state.Spec.State)
	klog.Infof("Current history is %s", state.Status.History)

	c.fn(*state)

	// c.recorder.Event(state, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Watcher) enqueueRequestState(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}
