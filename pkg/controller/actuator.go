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
const KubeconfigSecretName = "kubecfg"
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
	fleetManager, err := NewManagerForConfig(fleetClientConfig, "clusters") //TODO get from config
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

	cfg := &config.Config{}
	if ex.Spec.ProviderConfig != nil { //here we parse providerconfig
		if _, _, err := a.decoder.Decode(ex.Spec.ProviderConfig.Raw, nil, cfg); err != nil {
			return fmt.Errorf("failed to decode provider config: %+v", err)
		}
	}

	a.registerClusterInFleetManager(ctx, namespace, cluster)
	return a.updateStatus(ctx, ex)
}

// Delete the Extension resource.
func (a *actuator) Delete(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	a.logger.Info("Component is being deleted", "component", "fleet-agent-management", "namespace", namespace)
	return nil
}

// Restore the Extension resource.
func (a *actuator) Restore(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	a.logger.Info("Component is being restored", "component", "fleet-agent-management")
	return a.Reconcile(ctx, ex)
}

// Migrate the Extension resource.
func (a *actuator) Migrate(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	a.logger.Info("Component is being migrated", "component", "fleet-agent-management")

	return a.Delete(ctx, ex)
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

func (a *actuator) registerClusterInFleetManager(ctx context.Context, namespace string, cluster *extensions.Cluster) {
	a.logger.Info("Starting with already registered check")
	registered, err := a.fleetManager.GetCluster(ctx, cluster.Shoot.Name)
	if !errors.IsNotFound(err) {
		a.logger.Info("Cluster already registered - skipping registration", "clientId", registered.Spec.ClientID)
		return
	} else {
		a.logger.Info("Cluster registration not found.")
	}
	a.logger.Info("Starting cluster registration process")
	secret := &corev1.Secret{}

	labels := make(map[string]string)
	labels["corebundle"] = "true"
	labels["region"] = cluster.Shoot.Spec.Region
	labels["cluster"] = cluster.Shoot.Name
	if a.serviceConfig.FleetAgentConfig.Labels != nil && len(a.serviceConfig.FleetAgentConfig.Labels) > 0 { //adds labels from configuration
		for key, value := range a.serviceConfig.Labels {
			labels[key] = value
		}
	}
	a.logger.Info("Looking up Secret with KubeConfig for given Shoot.", "namespace", namespace, "secretName", KubeconfigSecretName)

	if err := a.client.Get(ctx, kutil.Key(namespace, KubeconfigSecretName), secret); err == nil {
		secretData := make(map[string][]byte)
		secretData["value"] = secret.Data[KubeconfigKey]
		a.logger.Info("Loaded kubeconfig from secret", "kubeconfig", secret, "namespace", namespace)

		const fleetRegisterNamespace = "clusters"
		kubeconfigSecret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kubecfg-" + cluster.Shoot.Name,
				Namespace: fleetRegisterNamespace,
			},
			Data: secretData,
		}

		clusterRegistration := fleetv1alpha1.Cluster{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cluster.Shoot.Name,
				Namespace: fleetRegisterNamespace,
				Labels:    labels,
			},
			Spec: fleetv1alpha1.ClusterSpec{
				KubeConfigSecret: "kubecfg-" + cluster.Shoot.Name,
			},
		}
		a.logger.Info("Creating kubeconfig secret for Fleet registration.")
		if _, err = a.fleetManager.CreateKubeconfigSecret(ctx, &kubeconfigSecret); err != nil {
			a.logger.Error(err, "Failed to create secret with kubeconfig for Fleet registration")
		}
		a.logger.Info("Creating Cluster registration for Fleet registration.")
		if _, err = a.fleetManager.CreateCluster(ctx, &clusterRegistration); err != nil {
			a.logger.Error(err, "Failed to create Cluster for Fleet registration")
		}
		a.logger.Info("Registered shoot cluster in Fleet Manager ", "registration", clusterRegistration)
	} else {
		a.logger.Error(err, "Failed to find Secret with kubeconfig for Fleet registration.")
	}
}

func (a *actuator) updateStatus(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	return controller.TryUpdateStatus(ctx, retry.DefaultBackoff, a.client, ex, func() error {
		return nil
	})
}
