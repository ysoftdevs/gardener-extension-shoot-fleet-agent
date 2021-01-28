package v1alpha1

import (
	healthcheckconfig "github.com/gardener/gardener/extensions/pkg/controller/healthcheck/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	componentbaseconfigv1alpha1 "k8s.io/component-base/config/v1alpha1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FleetAgentConfig configuration resource
type FleetAgentConfig struct {
	metav1.TypeMeta `json:",inline"`

	// ClientConnection specifies the kubeconfig file and client connection
	// settings for the proxy server to use when communicating with the apiserver.
	// +optional
	ClientConnection *componentbaseconfigv1alpha1.ClientConnectionConfiguration `json:"clientConnection,omitempty"`

	// labels to use in Fleet Cluster registration
	Labels map[string]string `json:"labels,omitempty"`

	//namespace to store clusters registrations in Fleet managers cluster
	Namespace string `json:"namespace,omitempty"`

	HealthCheckConfig *healthcheckconfig.HealthCheckConfig `json:"healthCheckConfig,omitempty"`
}
