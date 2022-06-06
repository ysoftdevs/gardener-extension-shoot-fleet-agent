############# builder
FROM golang:1.17.10 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-shoot-fleet-agent
COPY . .
RUN make install

############# gardener-extension-shoot-fleet-agent
FROM alpine:3.15.4 AS gardener-extension-shoot-fleet-agent

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-shoot-fleet-agent /gardener-extension-shoot-fleet-agent
ENTRYPOINT ["/gardener-extension-shoot-fleet-agent"]
