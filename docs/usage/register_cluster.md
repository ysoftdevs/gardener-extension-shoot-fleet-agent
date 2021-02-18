# Register Shoot cluster in Fleet manager

## Introduction
Gardener takes care of provisioning clusters. It doesn't install anything into created clusters.
This extension enables App instalation via [Fleet](https://fleet.rancher.io) by registering newly created Shoot clusters into Fleet manager.

### Service Scope
This service enables users to register Shoot cluster in Fleet.
```yaml
kind: Shoot
...
spec:
  extensions:
  - type: shoot-fleet-agent
    providerConfig:
      apiVersion: service.fleet-agent.extensions.gardener.cloud/v1alpha1
      kind:
      defaultConfig:
        clientConnection:
          kubeconfig: base64 encoded kubeconfig
        labels:
          env: test
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
