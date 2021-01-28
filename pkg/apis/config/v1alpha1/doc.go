// +k8s:deepcopy-gen=package
// +k8s:conversion-gen=github.com/javamachr/gardener-extension-shoot-fleet-agent/pkg/apis/config
// +k8s:openapi-gen=true
// +k8s:defaulter-gen=TypeMeta

//go:generate gen-crd-api-reference-docs -api-dir . -config ../../../../hack/api-reference/config.json -template-dir ../../../../vendor/github.com/gardener/gardener/hack/api-reference/template -out-file ../../../../hack/api-reference/config.md

// Package v1alpha1 contains the Azure provider configuration API resources.
// +groupName=shoot-fleet-agent-service.extensions.config.gardener.cloud
package v1alpha1
