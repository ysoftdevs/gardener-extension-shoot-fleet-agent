package config

import (
	healthcheckconfig "github.com/gardener/gardener/extensions/pkg/controller/healthcheck/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	componentbaseconfig "k8s.io/component-base/config"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FleetAgentConfig configuration resource
type FleetAgentConfig struct {
	metav1.TypeMeta

	// ClientConnection specifies the kubeconfig file and client connection
	// settings for the proxy server to use when communicating with the apiserver.
	ClientConnection *componentbaseconfig.ClientConnectionConfiguration

	// labels to use in Fleet Cluster registration
	Labels map[string]string

	//namespace to store clusters registrations in Fleet managers cluster
	Namespace string

	HealthCheckConfig *healthcheckconfig.HealthCheckConfig
}
