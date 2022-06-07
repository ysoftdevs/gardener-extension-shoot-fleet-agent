<p>Packages:</p>
<ul>
<li>
<a href="#shoot-fleet-agent-service.extensions.config.gardener.cloud%2fv1alpha1">shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1</a>
</li>
</ul>
<h2 id="shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1">shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1</h2>
<p>
<p>Package v1alpha1 contains the Azure provider configuration API resources.</p>
</p>
Resource Types:
<ul><li>
<a href="#shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1.FleetAgentConfig">FleetAgentConfig</a>
</li></ul>
<h3 id="shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1.FleetAgentConfig">FleetAgentConfig
</h3>
<p>
<p>FleetAgentConfig configuration resource</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>FleetAgentConfig</code></td>
</tr>
<tr>
<td>
<code>defaultConfig</code></br>
<em>
<a href="#shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1.ProjectConfig">
ProjectConfig
</a>
</em>
</td>
<td>
<p>DefaultConfiguration holds default config applied if no project config found</p>
</td>
</tr>
<tr>
<td>
<code>projectConfig</code></br>
<em>
<a href="#shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1.ProjectConfig">
map[string]github.com/ysoftdevs/gardener-extension-shoot-fleet-agent/pkg/apis/config/v1alpha1.ProjectConfig
</a>
</em>
</td>
<td>
<p>ProjectConfiguration holds configuration overrides for each project</p>
</td>
</tr>
<tr>
<td>
<code>healthCheckConfig</code></br>
<em>
<a href="https://github.com/gardener/gardener/extensions/pkg/controller/healthcheck/config">
github.com/gardener/gardener/extensions/pkg/controller/healthcheck/config.HealthCheckConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1.ProjectConfig">ProjectConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#shoot-fleet-agent-service.extensions.config.gardener.cloud/v1alpha1.FleetAgentConfig">FleetAgentConfig</a>)
</p>
<p>
<p>ProjectConfig holds configuration for single project</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>kubeconfig</code></br>
<em>
string
</em>
</td>
<td>
<p>Kubeconfig contains base64 encoded kubeconfig</p>
</td>
</tr>
<tr>
<td>
<code>labels</code></br>
<em>
map[string]string
</em>
</td>
<td>
<p>labels to use in Fleet Cluster registration</p>
</td>
</tr>
<tr>
<td>
<code>namespace</code></br>
<em>
string
</em>
</td>
<td>
<p>namespace to store clusters registrations in Fleet managers cluster</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <a href="https://github.com/ahmetb/gen-crd-api-reference-docs">gen-crd-api-reference-docs</a>
</em></p>
