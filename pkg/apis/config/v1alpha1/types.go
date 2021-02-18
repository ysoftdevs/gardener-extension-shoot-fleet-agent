package v1alpha1

import (
	healthcheckconfig "github.com/gardener/gardener/extensions/pkg/controller/healthcheck/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FleetAgentConfig configuration resource
type FleetAgentConfig struct {
	metav1.TypeMeta `json:",inline"`

	// DefaultConfiguration holds default config applied if no project config found
	DefaultConfiguration ProjectConfig `json:"defaultConfig,omitempty"`

	// ProjectConfiguration holds configuration overrides for each project
	ProjectConfiguration map[string]ProjectConfig `json:"projectConfig,omitempty"`

	HealthCheckConfig *healthcheckconfig.HealthCheckConfig `json:"healthCheckConfig,omitempty"`
}

// ProjectConfig holds configuration for single project
type ProjectConfig struct {
	// Kubeconfig contains base64 encoded kubeconfig
	Kubeconfig string `json:"kubeconfig,omitempty"`

	// labels to use in Fleet Cluster registration
	Labels map[string]string `json:"labels,omitempty"`

	//namespace to store clusters registrations in Fleet managers cluster
	Namespace string `json:"namespace,omitempty"`
}
