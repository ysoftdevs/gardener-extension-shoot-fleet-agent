/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

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

	v1alpha1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterRegistrationTokens implements ClusterRegistrationTokenInterface
type FakeClusterRegistrationTokens struct {
	Fake *FakeFleetV1alpha1
	ns   string
}

var clusterregistrationtokensResource = schema.GroupVersionResource{Group: "fleet.cattle.io", Version: "v1alpha1", Resource: "clusterregistrationtokens"}

var clusterregistrationtokensKind = schema.GroupVersionKind{Group: "fleet.cattle.io", Version: "v1alpha1", Kind: "ClusterRegistrationToken"}

// Get takes name of the clusterRegistrationToken, and returns the corresponding clusterRegistrationToken object, and an error if there is any.
func (c *FakeClusterRegistrationTokens) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ClusterRegistrationToken, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(clusterregistrationtokensResource, c.ns, name), &v1alpha1.ClusterRegistrationToken{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegistrationToken), err
}

// List takes label and field selectors, and returns the list of ClusterRegistrationTokens that match those selectors.
func (c *FakeClusterRegistrationTokens) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ClusterRegistrationTokenList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(clusterregistrationtokensResource, clusterregistrationtokensKind, c.ns, opts), &v1alpha1.ClusterRegistrationTokenList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ClusterRegistrationTokenList{ListMeta: obj.(*v1alpha1.ClusterRegistrationTokenList).ListMeta}
	for _, item := range obj.(*v1alpha1.ClusterRegistrationTokenList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterRegistrationTokens.
func (c *FakeClusterRegistrationTokens) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(clusterregistrationtokensResource, c.ns, opts))

}

// Create takes the representation of a clusterRegistrationToken and creates it.  Returns the server's representation of the clusterRegistrationToken, and an error, if there is any.
func (c *FakeClusterRegistrationTokens) Create(ctx context.Context, clusterRegistrationToken *v1alpha1.ClusterRegistrationToken, opts v1.CreateOptions) (result *v1alpha1.ClusterRegistrationToken, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(clusterregistrationtokensResource, c.ns, clusterRegistrationToken), &v1alpha1.ClusterRegistrationToken{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegistrationToken), err
}

// Update takes the representation of a clusterRegistrationToken and updates it. Returns the server's representation of the clusterRegistrationToken, and an error, if there is any.
func (c *FakeClusterRegistrationTokens) Update(ctx context.Context, clusterRegistrationToken *v1alpha1.ClusterRegistrationToken, opts v1.UpdateOptions) (result *v1alpha1.ClusterRegistrationToken, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(clusterregistrationtokensResource, c.ns, clusterRegistrationToken), &v1alpha1.ClusterRegistrationToken{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegistrationToken), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeClusterRegistrationTokens) UpdateStatus(ctx context.Context, clusterRegistrationToken *v1alpha1.ClusterRegistrationToken, opts v1.UpdateOptions) (*v1alpha1.ClusterRegistrationToken, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(clusterregistrationtokensResource, "status", c.ns, clusterRegistrationToken), &v1alpha1.ClusterRegistrationToken{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegistrationToken), err
}

// Delete takes name of the clusterRegistrationToken and deletes it. Returns an error if one occurs.
func (c *FakeClusterRegistrationTokens) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(clusterregistrationtokensResource, c.ns, name), &v1alpha1.ClusterRegistrationToken{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterRegistrationTokens) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(clusterregistrationtokensResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ClusterRegistrationTokenList{})
	return err
}

// Patch applies the patch and returns the patched clusterRegistrationToken.
func (c *FakeClusterRegistrationTokens) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterRegistrationToken, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(clusterregistrationtokensResource, c.ns, name, pt, data, subresources...), &v1alpha1.ClusterRegistrationToken{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegistrationToken), err
}
