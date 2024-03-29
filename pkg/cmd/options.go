// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"github.com/gardener/gardener/extensions/pkg/controller/cmd"
	extensionshealthcheckcontroller "github.com/gardener/gardener/extensions/pkg/controller/healthcheck"
	healthcheckconfig "github.com/gardener/gardener/extensions/pkg/controller/healthcheck/config"
	"github.com/spf13/pflag"
	apisconfig "github.com/ysoftdevs/gardener-extension-shoot-fleet-agent/pkg/apis/config"
	"github.com/ysoftdevs/gardener-extension-shoot-fleet-agent/pkg/apis/config/v1alpha1"
	"github.com/ysoftdevs/gardener-extension-shoot-fleet-agent/pkg/controller"
	controllerconfig "github.com/ysoftdevs/gardener-extension-shoot-fleet-agent/pkg/controller/config"
	healthcheckcontroller "github.com/ysoftdevs/gardener-extension-shoot-fleet-agent/pkg/controller/healthcheck"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	scheme  *runtime.Scheme
	decoder runtime.Decoder
)

func init() {
	scheme = runtime.NewScheme()
	utilruntime.Must(apisconfig.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	decoder = serializer.NewCodecFactory(scheme).UniversalDecoder()
}

// FleetServiceOptions holds options related to the fleet agent service.
type FleetServiceOptions struct {
	ConfigLocation string
	config         *FleetServiceConfig
}

// AddFlags implements Flagger.AddFlags.
func (o *FleetServiceOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ConfigLocation, "config", "", "Path to fleet agent service configuration")
}

// Complete implements Completer.Complete.
func (o *FleetServiceOptions) Complete() error {
	if o.ConfigLocation == "" {
		return errors.New("config location is not set")
	}

	data, err := ioutil.ReadFile(o.ConfigLocation)
	if err != nil {
		return err
	}

	config := apisconfig.FleetAgentConfig{}
	_, _, err = decoder.Decode(data, nil, &config)
	if err != nil {
		return err
	}

	o.config = &FleetServiceConfig{
		config: config,
	}

	return nil
}

// Completed returns the decoded FleetServiceConfiguration instance. Only call this if `Complete` was successful.
func (o *FleetServiceOptions) Completed() *FleetServiceConfig {
	return o.config
}

// FleetServiceConfig contains configuration information about the fleet service.
type FleetServiceConfig struct {
	config apisconfig.FleetAgentConfig
}

// Apply applies the FleetServiceOptions to the passed ControllerOptions instance.
func (c *FleetServiceConfig) Apply(config *controllerconfig.Config) {
	config.FleetAgentConfig = c.config
}

// ControllerSwitches are the cmd.SwitchOptions for the provider controllers.
func ControllerSwitches() *cmd.SwitchOptions {
	return cmd.NewSwitchOptions(
		cmd.Switch(controller.ControllerName, controller.AddToManager),
		cmd.Switch(extensionshealthcheckcontroller.ControllerName, healthcheckcontroller.AddToManager),
	)
}

// ApplyHealthCheckConfig applies healthcheck config
func (c *FleetServiceConfig) ApplyHealthCheckConfig(config *healthcheckconfig.HealthCheckConfig) {
	if c.config.HealthCheckConfig != nil {
		*config = *c.config.HealthCheckConfig
	}
}
