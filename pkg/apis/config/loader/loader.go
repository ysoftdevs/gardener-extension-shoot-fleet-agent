package loader

import (
	"io/ioutil"

	"github.com/javamachr/gardener-extension-shoot-fleet-agent/pkg/apis/config"
	"github.com/javamachr/gardener-extension-shoot-fleet-agent/pkg/apis/config/install"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
)

var (
	//codec used to encode
	Codec runtime.Codec
	//scheme to register
	Scheme *runtime.Scheme
)

func init() {
	Scheme = runtime.NewScheme()
	install.Install(Scheme)
	yamlSerializer := json.NewYAMLSerializer(json.DefaultMetaFactory, Scheme, Scheme)
	Codec = versioning.NewDefaultingCodecForScheme(
		Scheme,
		yamlSerializer,
		yamlSerializer,
		schema.GroupVersion{Version: "v1alpha1"},
		runtime.InternalGroupVersioner,
	)
}

// LoadFromFile takes a filename and de-serializes the contents into FleetAgentConfig object.
func LoadFromFile(filename string) (*config.FleetAgentConfig, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return Load(bytes)
}

// Load takes a byte slice and de-serializes the contents into FleetAgentConfig object.
// Encapsulates de-serialization without assuming the source is a file.
func Load(data []byte) (*config.FleetAgentConfig, error) {
	cfg := &config.FleetAgentConfig{}

	if len(data) == 0 {
		return cfg, nil
	}

	decoded, _, err := Codec.Decode(data, &schema.GroupVersionKind{Version: "v1alpha1", Kind: "Config"}, cfg)
	if err != nil {
		return nil, err
	}

	return decoded.(*config.FleetAgentConfig), nil
}
