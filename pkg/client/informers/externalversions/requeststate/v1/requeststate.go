/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	requeststatev1 "github.com/dmartinol/crd-watcher/pkg/apis/requeststate/v1"
	versioned "github.com/dmartinol/crd-watcher/pkg/client/clientset/versioned"
	internalinterfaces "github.com/dmartinol/crd-watcher/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/dmartinol/crd-watcher/pkg/client/listers/requeststate/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// RequestStateInformer provides access to a shared informer and lister for
// RequestStates.
type RequestStateInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.RequestStateLister
}

type requestStateInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewRequestStateInformer constructs a new informer for RequestState type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewRequestStateInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredRequestStateInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredRequestStateInformer constructs a new informer for RequestState type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredRequestStateInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ComV1().RequestStates(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ComV1().RequestStates(namespace).Watch(context.TODO(), options)
			},
		},
		&requeststatev1.RequestState{},
		resyncPeriod,
		indexers,
	)
}

func (f *requestStateInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredRequestStateInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *requestStateInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&requeststatev1.RequestState{}, f.defaultInformer)
}

func (f *requestStateInformer) Lister() v1.RequestStateLister {
	return v1.NewRequestStateLister(f.Informer().GetIndexer())
}
