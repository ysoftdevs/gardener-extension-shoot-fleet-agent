# [Gardener Extension for Fleet agent installation](https://gardener.cloud)



[![CircleCI](https://circleci.com/gh/javamachr/gardener-extension-shoot-fleet-agent.svg?style=shield)](https://circleci.com/gh/javamachr/gardener-extension-shoot-fleet-agent)
[![Go Report Card](https://goreportcard.com/badge/github.com/ysoftdevs/gardener-extension-shoot-fleet-agent)](https://goreportcard.com/report/github.com/ysoftdevs/gardener-extension-shoot-fleet-agent)

Project Gardener implements the automated management and operation of [Kubernetes](https://kubernetes.io/) clusters as a service. Its main principle is to leverage Kubernetes concepts for all of its tasks.

Recently, most of the vendor specific logic has been developed [in-tree](https://github.com/gardener/gardener). However, the project has grown to a size where it is very hard to extend, maintain, and test. With [GEP-1](https://github.com/gardener/gardener/blob/master/docs/proposals/01-extensibility.md) we have proposed how the architecture can be changed in a way to support external controllers that contain their very own vendor specifics. This way, we can keep Gardener core clean and independent.

## Configuration

Example configuration for this extension controller with default config and 2 project specific configs:

```yaml
apiVersion: shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1
kind: Configuration
defaultConfig:
  kubeconfig: #base64encoded kubeconfig of cluster running Fleet manager
  labels: #extra labels to apply to Cluster registration
    env: dev
projectConfig:  
  myproject:
    kubeconfig: #base64encoded kubeconfig of cluster running Fleet manager
    labels: 
      project: myproject
      somelabel: samevalue
    namespace: clusters
  myotherproject:
    kubeconfig: #base64encoded kubeconfig of cluster running Fleet manager
    labels: 
      project: mypotherproject
    namespace: fleetclusters
```
Unless the shoot is in one of the myproject or myother project it will get defaultConfiguration.
Configuration supplied in ControllerRegistration will be applied as default however it can be overriden via configuration in Extension for given Shoot.
## Extension-Resources

Example extension resource:

```yaml
apiVersion: extensions.gardener.cloud/v1alpha1
kind: Extension
metadata:
  name: "extension-shoot-fleet-agent"
  namespace: shoot--project--abc
spec:
  type: shoot-fleet-agent
```

When an extension resource is reconciled, the extension controller will register Shoot cluster in Fleet management cluster(configured in kubeconfig in Configuration object above.

Please note, this extension controller relies on existing properly configured [Fleet multi-cluster deployment](https://fleet.rancher.io/multi-cluster-install/) configured above.
By default configuration from ControllerRegistration is used but if Extension specifies different config for particular Shoot that configuration overrides the default settings.

## Shooted seeds registration in Fleet
All shooted seeds will be ignored and not registered in Fleet. If you need to register shooted seeds please do so manually.

## Fleet CRD cluster default labels
By default all clusters registered in fleet will have following labels:

|label|value|
|---|---|
|corebundle|true|
|region| region of shoot |
|cluster| name of shoot|
|seed|name of shoots seed|

Any labels defined in configuration are added to those default labels.

## How to start using or developing this extension controller locally

You can run the controller locally on your machine by executing `make start`. Please make sure to have the kubeconfig to the cluster you want to connect to ready in the `./dev/kubeconfig` file.
Static code checks and tests can be executed by running `VERIFY=true make all`. We are using Go modules for Golang package dependency management and [Ginkgo](https://github.com/onsi/ginkgo)/[Gomega](https://github.com/onsi/gomega) for testing.

## Feedback and Support

Feedback and contributions are always welcome. Please report bugs or suggestions as [GitHub issues](https://github.com/ysoftdevs/gardener-extension-shoot-fleet-agent/issues) or join our [Slack channel #gardener](https://kubernetes.slack.com/messages/gardener) (please invite yourself to the Kubernetes workspace [here](http://slack.k8s.io)).

## Learn more!

Please find further resources about out project here:

* [Our landing page gardener.cloud](https://gardener.cloud/)
* ["Gardener, the Kubernetes Botanist" blog on kubernetes.io](https://kubernetes.io/blog/2018/05/17/gardener/)
* ["Gardener Project Update" blog on kubernetes.io](https://kubernetes.io/blog/2019/12/02/gardener-project-update/)
* [Gardener Extensions Golang library](https://godoc.org/github.com/gardener/gardener/extensions/pkg)
* [GEP-1 (Gardener Enhancement Proposal) on extensibility](https://github.com/gardener/gardener/blob/master/docs/proposals/01-extensibility.md)
* [Extensibility API documentation](https://github.com/gardener/gardener/tree/master/docs/extensions)
