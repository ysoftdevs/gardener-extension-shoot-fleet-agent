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
	v1alpha1 "github.com/javamachr/gardener-extension-shoot-fleet-agent/pkg/client/fleet/clientset/versioned/typed/fleet.cattle.io/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeFleetV1alpha1 struct {
	*testing.Fake
}

func (c *FakeFleetV1alpha1) Bundles(namespace string) v1alpha1.BundleInterface {
	return &FakeBundles{c, namespace}
}

func (c *FakeFleetV1alpha1) BundleDeployments(namespace string) v1alpha1.BundleDeploymentInterface {
	return &FakeBundleDeployments{c, namespace}
}

func (c *FakeFleetV1alpha1) BundleNamespaceMappings(namespace string) v1alpha1.BundleNamespaceMappingInterface {
	return &FakeBundleNamespaceMappings{c, namespace}
}

func (c *FakeFleetV1alpha1) Clusters(namespace string) v1alpha1.ClusterInterface {
	return &FakeClusters{c, namespace}
}

func (c *FakeFleetV1alpha1) ClusterGroups(namespace string) v1alpha1.ClusterGroupInterface {
	return &FakeClusterGroups{c, namespace}
}

func (c *FakeFleetV1alpha1) ClusterRegistrations(namespace string) v1alpha1.ClusterRegistrationInterface {
	return &FakeClusterRegistrations{c, namespace}
}

func (c *FakeFleetV1alpha1) ClusterRegistrationTokens(namespace string) v1alpha1.ClusterRegistrationTokenInterface {
	return &FakeClusterRegistrationTokens{c, namespace}
}

func (c *FakeFleetV1alpha1) Contents() v1alpha1.ContentInterface {
	return &FakeContents{c}
}

func (c *FakeFleetV1alpha1) GitRepos(namespace string) v1alpha1.GitRepoInterface {
	return &FakeGitRepos{c, namespace}
}

func (c *FakeFleetV1alpha1) GitRepoRestrictions(namespace string) v1alpha1.GitRepoRestrictionInterface {
	return &FakeGitRepoRestrictions{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeFleetV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
