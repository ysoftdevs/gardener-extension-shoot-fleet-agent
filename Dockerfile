############# builder
FROM eu.gcr.io/gardener-project/3rd/golang:1.15.5 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-shoot-fleet-agent
COPY . .
RUN make install

############# gardener-extension-shoot-fleet-agent
FROM eu.gcr.io/gardener-project/3rd/alpine:3.12.3 AS gardener-extension-shoot-fleet-agent

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-shoot-fleet-agent /gardener-extension-shoot-fleet-agent
ENTRYPOINT ["/gardener-extension-shoot-fleet-agent"]
