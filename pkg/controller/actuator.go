// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"reflect"

	"github.com/gardener/gardener/pkg/extensions"
	"github.com/go-logr/logr"
	fleetv1alpha1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	"github.com/javamachr/gardener-extension-shoot-fleet-agent/pkg/controller/config"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
)

// ActuatorName is the name of the Fleet agent actuator.
const ActuatorName = "shoot-fleet-agent-actuator"

// KubeconfigSecretName name of secret that holds kubeconfig for Shoot
const KubeconfigSecretName = "kubecfg"

// KubeconfigKey key in KubeconfigSecretName secret that holds kubeconfig for Shoot
const KubeconfigKey = "kubeconfig"

// NewActuator returns an actuator responsible for Extension resources.
func NewActuator(config config.Config) extension.Actuator {
	fleetKubeConfig, _ := b64.StdEncoding.DecodeString(config.FleetAgentConfig.ClientConnection.Kubeconfig)
	var kubeconfigPath string
	var err error
	if kubeconfigPath, err = writeKubeconfigToTempFile(fleetKubeConfig); err != nil {
		panic(err)
	}
	fleetClientConfig, _ := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	var fleetNamespace = "clusters"
	if len(config.Namespace) != 0 {
		fleetNamespace = config.Namespace
	}
	fleetManager, err := NewManagerForConfig(fleetClientConfig, fleetNamespace)
	if err != nil {
		panic(err)
	}

	return &actuator{
		logger:        log.Log.WithName(ActuatorName),
		serviceConfig: config,
		fleetManager:  fleetManager,
	}
}

type actuator struct {
	client       client.Client
	config       *rest.Config
	decoder      runtime.Decoder
	fleetManager *FleetManager

	serviceConfig config.Config

	logger logr.Logger
}

// Reconcile the Extension resource.
func (a *actuator) Reconcile(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	a.logger.Info("Component is being reconciled", "component", "fleet-agent-management", "namespace", namespace)
	cluster, err := controller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}
	shootsConfigOverride := &config.Config{}
	if ex.Spec.ProviderConfig != nil { //parse providerConfig defaults override for this Shoot
		if _, _, err := a.decoder.Decode(ex.Spec.ProviderConfig.Raw, nil, shootsConfigOverride); err != nil {
			return fmt.Errorf("failed to decode provider config: %+v", err)
		}
	}
	a.ReconcileClusterInFleetManager(ctx, namespace, cluster, shootsConfigOverride)
	return a.updateStatus(ctx, ex)
}

// Delete the Extension resource.
func (a *actuator) Delete(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	cluster, err := controller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}
	a.logger.Info("Component is being deleted", "component", "fleet-agent-management", "namespace", namespace, "cluster", buildCrdName(cluster))
	err = a.fleetManager.DeleteKubeconfigSecret(ctx, buildCrdName(cluster))
	if err != nil {
		a.logger.Error(err, "Failed to delete kubeconfig secret for Shoot cluster.", "cluster", buildCrdName(cluster))
	}
	err = a.fleetManager.DeleteCluster(ctx, buildCrdName(cluster))
	if err != nil {
		a.logger.Error(err, "Failed to delete Cluster registration for Shoot cluster.", "cluster", buildCrdName(cluster))
	}
	return nil
}

// Restore the Extension resource.
func (a *actuator) Restore(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	//NOOP as there are no resources by this controller in Seed
	return nil
}

// Migrate the Extension resource.
func (a *actuator) Migrate(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	//NOOP as there are no resources by this controller in Seed
	return nil
}

// InjectConfig injects the rest config to this actuator.
func (a *actuator) InjectConfig(config *rest.Config) error {
	a.config = config
	return nil
}

// InjectClient injects the controller runtime client into the reconciler.
func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

// InjectScheme injects the given scheme into the reconciler.
func (a *actuator) InjectScheme(scheme *runtime.Scheme) error {
	a.decoder = serializer.NewCodecFactory(scheme).UniversalDecoder()
	return nil
}

// ReconcileClusterInFleetManager reconciles cluster registration in remote fleet manager
func (a *actuator) ReconcileClusterInFleetManager(ctx context.Context, namespace string, cluster *extensions.Cluster, override *config.Config) {
	a.logger.Info("Starting with already registered check")
	labels := prepareLabels(cluster, a.serviceConfig, override)
	registered, err := a.fleetManager.GetCluster(ctx, cluster.Shoot.Name)
	if !errors.IsNotFound(err) {
		if reflect.DeepEqual(registered.Labels, labels) {
			a.logger.Info("Cluster already registered - skipping registration", "clientId", registered.Spec.ClientID)
		} else {
			a.logger.Info("Updating labels of already registered cluster.", "clientId", registered.Spec.ClientID)
			a.updateClusterLabelsInFleet(ctx, registered, labels)
		}
		return
	}
	a.registerNewClusterInFleet(ctx, namespace, cluster, labels)
}

func (a *actuator) updateClusterLabelsInFleet(ctx context.Context, clusterRegistration *fleetv1alpha1.Cluster, labels map[string]string) {
	clusterRegistration.Labels = labels
	_, err := a.fleetManager.UpdateCluster(ctx, clusterRegistration)
	if err != nil {
		a.logger.Error(err, "Failed to update cluster labels in Fleet registration.", "clusterName", clusterRegistration.Name)
	}
}

func (a *actuator) registerNewClusterInFleet(ctx context.Context, namespace string, cluster *extensions.Cluster, labels map[string]string) {
	a.logger.Info("Looking up Secret with KubeConfig for given Shoot.", "namespace", namespace, "secretName", KubeconfigSecretName)
	secret := &corev1.Secret{}
	if err := a.client.Get(ctx, kutil.Key(namespace, KubeconfigSecretName), secret); err == nil {
		secretData := make(map[string][]byte)
		secretData["value"] = secret.Data[KubeconfigKey]
		a.logger.Info("Loaded kubeconfig from secret", "kubeconfig", secret, "namespace", namespace)

		const fleetRegisterNamespace = "clusters"
		kubeconfigSecret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kubecfg-" + buildCrdName(cluster),
				Namespace: fleetRegisterNamespace,
			},
			Data: secretData,
		}

		clusterRegistration := fleetv1alpha1.Cluster{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      buildCrdName(cluster),
				Namespace: fleetRegisterNamespace,
				Labels:    labels,
			},
			Spec: fleetv1alpha1.ClusterSpec{
				KubeConfigSecret: "kubecfg-" + buildCrdName(cluster),
			},
		}
		if _, err = a.fleetManager.CreateKubeconfigSecret(ctx, &kubeconfigSecret); err != nil {
			a.logger.Error(err, "Failed to create secret with kubeconfig for Fleet registration")
		}
		if _, err = a.fleetManager.CreateCluster(ctx, &clusterRegistration); err != nil {
			a.logger.Error(err, "Failed to create Cluster for Fleet registration")
		}
		a.logger.Info("Registered shoot cluster in Fleet Manager ", "registration", clusterRegistration)
	} else {
		a.logger.Error(err, "Failed to find Secret with kubeconfig for Fleet registration.")
	}
}

func prepareLabels(cluster *extensions.Cluster, serviceConfig config.Config, override *config.Config) map[string]string {
	labels := make(map[string]string)
	labels["corebundle"] = "true"
	labels["region"] = cluster.Shoot.Spec.Region
	labels["cluster"] = cluster.Shoot.Name
	if len(override.Labels) > 0 { //adds labels from Shoot configuration
		for key, value := range override.Labels {
			labels[key] = value
		}
	} else {
		if len(serviceConfig.FleetAgentConfig.Labels) > 0 { //adds labels from default configuration
			for key, value := range serviceConfig.Labels {
				labels[key] = value
			}
		}
	}
	return labels
}

// buildCrdName creates a unique name for cluster registration resources in Fleet manager cluster
func buildCrdName(cluster *extensions.Cluster) string {
	return cluster.Seed.Name + "" + cluster.Shoot.Name
}

func (a *actuator) updateStatus(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	return controller.TryUpdateStatus(ctx, retry.DefaultBackoff, a.client, ex, func() error {
		return nil
	})
}
