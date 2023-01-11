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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	requeststatev1 "github.com/dmartinol/crd-watcher/pkg/apis/requeststate/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRequestStates implements RequestStateInterface
type FakeRequestStates struct {
	Fake *FakeComV1
	ns   string
}

var requeststatesResource = schema.GroupVersionResource{Group: "com.redhat.ecosystem.sample", Version: "v1", Resource: "requeststates"}

var requeststatesKind = schema.GroupVersionKind{Group: "com.redhat.ecosystem.sample", Version: "v1", Kind: "RequestState"}

// Get takes name of the requestState, and returns the corresponding requestState object, and an error if there is any.
func (c *FakeRequestStates) Get(ctx context.Context, name string, options v1.GetOptions) (result *requeststatev1.RequestState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(requeststatesResource, c.ns, name), &requeststatev1.RequestState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*requeststatev1.RequestState), err
}

// List takes label and field selectors, and returns the list of RequestStates that match those selectors.
func (c *FakeRequestStates) List(ctx context.Context, opts v1.ListOptions) (result *requeststatev1.RequestStateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(requeststatesResource, requeststatesKind, c.ns, opts), &requeststatev1.RequestStateList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &requeststatev1.RequestStateList{ListMeta: obj.(*requeststatev1.RequestStateList).ListMeta}
	for _, item := range obj.(*requeststatev1.RequestStateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested requestStates.
func (c *FakeRequestStates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(requeststatesResource, c.ns, opts))

}

// Create takes the representation of a requestState and creates it.  Returns the server's representation of the requestState, and an error, if there is any.
func (c *FakeRequestStates) Create(ctx context.Context, requestState *requeststatev1.RequestState, opts v1.CreateOptions) (result *requeststatev1.RequestState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(requeststatesResource, c.ns, requestState), &requeststatev1.RequestState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*requeststatev1.RequestState), err
}

// Update takes the representation of a requestState and updates it. Returns the server's representation of the requestState, and an error, if there is any.
func (c *FakeRequestStates) Update(ctx context.Context, requestState *requeststatev1.RequestState, opts v1.UpdateOptions) (result *requeststatev1.RequestState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(requeststatesResource, c.ns, requestState), &requeststatev1.RequestState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*requeststatev1.RequestState), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRequestStates) UpdateStatus(ctx context.Context, requestState *requeststatev1.RequestState, opts v1.UpdateOptions) (*requeststatev1.RequestState, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(requeststatesResource, "status", c.ns, requestState), &requeststatev1.RequestState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*requeststatev1.RequestState), err
}

// Delete takes name of the requestState and deletes it. Returns an error if one occurs.
func (c *FakeRequestStates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(requeststatesResource, c.ns, name, opts), &requeststatev1.RequestState{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRequestStates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(requeststatesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &requeststatev1.RequestStateList{})
	return err
}

// Patch applies the patch and returns the patched requestState.
func (c *FakeRequestStates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *requeststatev1.RequestState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(requeststatesResource, c.ns, name, pt, data, subresources...), &requeststatev1.RequestState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*requeststatev1.RequestState), err
}