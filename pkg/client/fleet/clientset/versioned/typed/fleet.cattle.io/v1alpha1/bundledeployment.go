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

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	scheme "github.com/ysoftdevs/gardener-extension-shoot-fleet-agent/pkg/client/fleet/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BundleDeploymentsGetter has a method to return a BundleDeploymentInterface.
// A group's client should implement this interface.
type BundleDeploymentsGetter interface {
	BundleDeployments(namespace string) BundleDeploymentInterface
}

// BundleDeploymentInterface has methods to work with BundleDeployment resources.
type BundleDeploymentInterface interface {
	Create(ctx context.Context, bundleDeployment *v1alpha1.BundleDeployment, opts v1.CreateOptions) (*v1alpha1.BundleDeployment, error)
	Update(ctx context.Context, bundleDeployment *v1alpha1.BundleDeployment, opts v1.UpdateOptions) (*v1alpha1.BundleDeployment, error)
	UpdateStatus(ctx context.Context, bundleDeployment *v1alpha1.BundleDeployment, opts v1.UpdateOptions) (*v1alpha1.BundleDeployment, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.BundleDeployment, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.BundleDeploymentList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.BundleDeployment, err error)
	BundleDeploymentExpansion
}

// bundleDeployments implements BundleDeploymentInterface
type bundleDeployments struct {
	client rest.Interface
	ns     string
}

// newBundleDeployments returns a BundleDeployments
func newBundleDeployments(c *FleetV1alpha1Client, namespace string) *bundleDeployments {
	return &bundleDeployments{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the bundleDeployment, and returns the corresponding bundleDeployment object, and an error if there is any.
func (c *bundleDeployments) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.BundleDeployment, err error) {
	result = &v1alpha1.BundleDeployment{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bundledeployments").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BundleDeployments that match those selectors.
func (c *bundleDeployments) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BundleDeploymentList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.BundleDeploymentList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bundledeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested bundleDeployments.
func (c *bundleDeployments) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("bundledeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a bundleDeployment and creates it.  Returns the server's representation of the bundleDeployment, and an error, if there is any.
func (c *bundleDeployments) Create(ctx context.Context, bundleDeployment *v1alpha1.BundleDeployment, opts v1.CreateOptions) (result *v1alpha1.BundleDeployment, err error) {
	result = &v1alpha1.BundleDeployment{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("bundledeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(bundleDeployment).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a bundleDeployment and updates it. Returns the server's representation of the bundleDeployment, and an error, if there is any.
func (c *bundleDeployments) Update(ctx context.Context, bundleDeployment *v1alpha1.BundleDeployment, opts v1.UpdateOptions) (result *v1alpha1.BundleDeployment, err error) {
	result = &v1alpha1.BundleDeployment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("bundledeployments").
		Name(bundleDeployment.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(bundleDeployment).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *bundleDeployments) UpdateStatus(ctx context.Context, bundleDeployment *v1alpha1.BundleDeployment, opts v1.UpdateOptions) (result *v1alpha1.BundleDeployment, err error) {
	result = &v1alpha1.BundleDeployment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("bundledeployments").
		Name(bundleDeployment.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(bundleDeployment).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the bundleDeployment and deletes it. Returns an error if one occurs.
func (c *bundleDeployments) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bundledeployments").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *bundleDeployments) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bundledeployments").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched bundleDeployment.
func (c *bundleDeployments) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.BundleDeployment, err error) {
	result = &v1alpha1.BundleDeployment{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("bundledeployments").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
