package player

import (
	"flag"
	"strings"

	clientset "github.com/dmartinol/crd-watcher/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

type RunAs int64

const (
	DirectorMode RunAs = iota
	WorkerMode
)

func (r RunAs) String() string {
	switch r {
	case DirectorMode:
		return "Director"
	case WorkerMode:
		return "Worker"
	}
	return "unknown"
}

func RunAsFromString(runtimeMode string) RunAs {
	switch strings.ToLower(runtimeMode) {
	case "director":
		return DirectorMode
	case "worker":
		return WorkerMode
	}
	return DirectorMode
}

type Config struct {
	runAs       RunAs
	namespace   string
	config      *rest.Config
	stateclient *clientset.Clientset
}

func NewConfig() *Config {
	c := Config{
		runAs: DirectorMode,
	}

	runMode := flag.String("mode", "director", "Run mode, one of director, worker")
	flag.StringVar(&c.namespace, "namespace", "default", "Target namespace")

	var kubeconfig string
	var master string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "Absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "Master url")
	flag.Parse()

	c.runAs = RunAsFromString(*runMode)
	klog.Infof("Running as %s", c.runAs)

	// creates the connection
	var err error
	c.config, err = clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	c.stateclient, err = clientset.NewForConfig(c.config)
	if err != nil {
		klog.Fatal(err)
	}
	return &c
}

func (c *Config) RunAsDirector() bool {
	return c.runAs == DirectorMode
}

func (c *Config) RunAsWorker() bool {
	return c.runAs == WorkerMode
}

func (c *Config) Namespace() string {
	return c.namespace
}
