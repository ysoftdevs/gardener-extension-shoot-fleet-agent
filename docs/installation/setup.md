# Gardener Fleet agent management

## Introduction
Gardener comes with an extension that enables shoot owners to register their cluster in Fleet.

## Extension Installation
The `shoot-fleet-agent` extension can be deployed and configured via Gardener's native resource [ControllerRegistration](https://github.com/gardener/gardener/blob/master/docs/extensions/controllerregistration.md).

### Prerequisites
To let the `shoot-fleet-agent` operate properly, you need to have:
- a working cluster with Fleet multicluster setup enabled
- have kubeconfig with read/write access to cluster.fleet.cattle.io and secret resrouces in some namespace

### ControllerRegistration
An example of a `ControllerRegistration` for the `shoot-fleet-agent` can be found here: https://github.com/javamachr/gardener-extension-shoot-fleet-agent/blob/master/example/controller-registration.yaml

### Configuration
The `ControllerRegistration` contains a Helm chart which eventually deploy the `shoot-fleet-agent` to seed clusters. 

```yaml
apiVersion: core.gardener.cloud/v1beta1
kind: ControllerRegistration
...
  values:
    defaultConfig:
        clientConnection:
          kubeconfig: abcd
        labels:
```

If the `shoot-fleet-agent` should be enabled for every shoot cluster in your Gardener managed environment, you need to globally enable it in the `ControllerRegistration`:
```yaml
apiVersion: core.gardener.cloud/v1beta1
kind: ControllerRegistration
...
  resources:
  - globallyEnabled: true
    kind: Extension
    type: shoot-fleet-agent
```

Alternatively, you're given the option to only enable the service for certain shoots:
```yaml
kind: Shoot
apiVersion: core.gardener.cloud/v1beta1
...
spec:
  extensions:
  - type: shoot-fleet-agent
...
```

<style>
#body-inner blockquote {
    border: 0;
    padding: 10px;
    margin-top: 40px;
    margin-bottom: 40px;
    border-radius: 4px;
    background-color: rgba(0,0,0,0.05);
    box-shadow: 0 3px 6px rgba(0,0,0,0.16), 0 3px 6px rgba(0,0,0,0.23);
    position:relative;
    padding-left:60px;
}
#body-inner blockquote:before {
    content: "!";
    font-weight: bold;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    background-color: #00a273;
    color: white;
    vertical-align: middle;
    margin: auto;
    width: 36px;
    font-size: 30px;
    text-align: center;
}
</style>
