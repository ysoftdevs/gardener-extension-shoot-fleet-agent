package managed_resource_handler

import (
	"context"
	"fmt"
	"github.com/gardener/gardener/extensions/pkg/util"
	"github.com/gardener/gardener/pkg/apis/resources/v1alpha1"
	"github.com/gardener/gardener/pkg/extensions"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	"github.com/gardener/gardener/pkg/utils/kubernetes/health"
	"github.com/gardener/gardener/pkg/utils/managedresources"
	"github.com/gardener/gardener/pkg/utils/retry"
	"github.com/go-logr/logr"
	token_requestor_handler "github.com/javamachr/gardener-extension-shoot-fleet-agent/pkg/controller/token-requestor-handler"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

var (
	ChartsPath                  = filepath.Join("charts", "internal")
	ManagedResourcesChartValues = map[string]interface{}{
		"shootAccessServiceAccountName":      token_requestor_handler.TargetServiceAccountName,
		"shootAccessServiceAccountNamespace": token_requestor_handler.TargetServiceAccountNamespace,
	}
)

const (
	ManagedResourceChartName = "shoot-fleet-agent-shoot"
	ManagedResourceName      = "extension-shoot-fleet-agent-shoot"
)

type ManagedResourceHandler struct {
	ctx     context.Context
	client  client.Client
	cluster *extensions.Cluster
	logger  logr.Logger
}

func NewManagedResourceHandler(ctx context.Context, client client.Client, cluster *extensions.Cluster) *ManagedResourceHandler {
	return &ManagedResourceHandler{
		ctx:     ctx,
		client:  client,
		cluster: cluster,
		logger:  log.Log.WithValues("logger", "managed-resource-handler", "cluster", cluster.ObjectMeta.Name),
	}
}

// renderChart takes the helm chart at charts/internal/<ManagedResourceChartName> and renders it with the values
func (h *ManagedResourceHandler) renderChart() ([]byte, error) {
	chartPath := filepath.Join(ChartsPath, ManagedResourceChartName)
	renderer, err := util.NewChartRendererForShoot(h.cluster.Shoot.Spec.Kubernetes.Version)
	if err != nil {
		return nil, err
	}
	chart, err := renderer.Render(chartPath, ManagedResourceChartName, h.cluster.ObjectMeta.Name, ManagedResourcesChartValues)
	if err != nil {
		return nil, err
	}

	return chart.Manifest(), nil
}

func (h *ManagedResourceHandler) EnsureManagedResoruces() error {
	h.logger.Info("Rendering chart with the managed resources")
	renderedChart, err := h.renderChart()
	if err != nil {
		return err
	}
	data := map[string][]byte{ManagedResourceChartName: renderedChart}
	keepObjects := false
	forceOverwriteAnnotations := false

	h.logger.Info("Creating the ManagedResource")
	if err := managedresources.Create(h.ctx, h.client, h.cluster.ObjectMeta.Name, ManagedResourceName, false, "", data, &keepObjects, map[string]string{}, &forceOverwriteAnnotations); err != nil {
		h.logger.Error(err, "ManagedResource could not be created")
		return err
	}

	h.logger.Info("Waiting until the ManagedResource is healthy")
	if err := h.waitUntilManagedResourceHealthy(); err != nil {
		h.logger.Error(err, "ManagedResource has been created, but hasn't been applied")
		return err
	}

	h.logger.Info("ManagedResource created")
	return nil
}

func (h *ManagedResourceHandler) waitUntilManagedResourceHealthy() error {
	name := ManagedResourceName
	namespace := h.cluster.ObjectMeta.Name

	obj := &v1alpha1.ManagedResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	return retry.UntilTimeout(h.ctx, 5*time.Second, 60*time.Second, func(ctx context.Context) (done bool, err error) {
		if err := h.client.Get(ctx, kutil.Key(namespace, name), obj); err != nil {
			h.logger.Info(fmt.Sprintf("Could not wait for the managed resource to be ready: %+v", err))
			return retry.MinorError(err)
		}

		// Check whether ManagedResource has proper status
		if err := health.CheckManagedResource(obj); err != nil {
			h.logger.Info(fmt.Sprintf("Managed resource %s not ready (yet): %+v", name, err))
			return retry.MinorError(err)
		}

		return retry.Ok()
	})
}
