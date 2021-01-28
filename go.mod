module github.com/javamachr/gardener-extension-shoot-fleet-agent

go 1.15

require (
	github.com/ahmetb/gen-crd-api-reference-docs v0.2.0
	github.com/gardener/etcd-druid v0.3.0 // indirect
	github.com/gardener/gardener v1.15.1-0.20210115062544-6dc08568692a
	github.com/go-logr/logr v0.3.0
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/golang/mock v1.4.4-0.20200731163441-8734ec565a4d
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rancher/fleet/pkg/apis v0.0.0-20200909045814-3675caaa7070
	github.com/rogpeppe/go-internal v1.7.0 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/ulikunitz/xz v0.5.7 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210123111255-9b0068b26619 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/tools v0.1.0 // indirect
	k8s.io/api v0.19.6
	k8s.io/apimachinery v0.19.6
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/code-generator v0.19.6
	k8s.io/component-base v0.19.6
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
	sigs.k8s.io/controller-runtime v0.7.1
)

replace (
	k8s.io/api => k8s.io/api v0.19.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.6
	k8s.io/apiserver => k8s.io/apiserver v0.19.6
	k8s.io/client-go => k8s.io/client-go v0.19.6
	k8s.io/code-generator => k8s.io/code-generator v0.19.6
	k8s.io/component-base => k8s.io/component-base v0.19.6
	k8s.io/helm => k8s.io/helm v2.13.1+incompatible
)
