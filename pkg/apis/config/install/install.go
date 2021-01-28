package install

import (
	"github.com/javamachr/gardener-extension-shoot-fleet-agent/pkg/apis/config"
	v1alpha1 "github.com/javamachr/gardener-extension-shoot-fleet-agent/pkg/apis/config/v1alpha1"
	v1alpha1fleet "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	schemeBuilder = runtime.NewSchemeBuilder(
		v1alpha1.AddToScheme,
		v1alpha1fleet.AddToScheme,
		config.AddToScheme,
		setVersionPriority,
	)

	// AddToScheme adds all APIs to the scheme.
	AddToScheme = schemeBuilder.AddToScheme
)

func setVersionPriority(scheme *runtime.Scheme) error {
	return scheme.SetVersionPriority(v1alpha1.SchemeGroupVersion)
}

// Install installs all APIs in the scheme.
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(AddToScheme(scheme))
}
