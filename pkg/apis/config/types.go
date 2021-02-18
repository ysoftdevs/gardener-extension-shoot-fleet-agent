package config

import (
	healthcheckconfig "github.com/gardener/gardener/extensions/pkg/controller/healthcheck/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FleetAgentConfig configuration resource
type FleetAgentConfig struct {
	metav1.TypeMeta

	// DefaultConfiguration holds default config applied if no project config found
	DefaultConfiguration ProjectConfig

	// ProjectConfiguration holds configuration overrides for each project
	ProjectConfiguration map[string]ProjectConfig

	HealthCheckConfig *healthcheckconfig.HealthCheckConfig
}

// ProjectConfig holds configuration for single project
type ProjectConfig struct {
	// Kubeconfig contains base64 encoded kubeconfig
	Kubeconfig string

	// labels to use in Fleet Cluster registration
	Labels map[string]string

	//namespace to store clusters registrations in Fleet managers cluster
	Namespace string
}
