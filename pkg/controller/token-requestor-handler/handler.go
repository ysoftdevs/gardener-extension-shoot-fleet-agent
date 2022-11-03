package token_requestor_handler

import (
	"context"
	"fmt"
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	constsv1alpha1 "github.com/gardener/gardener/pkg/apis/resources/v1alpha1"
	"github.com/gardener/gardener/pkg/extensions"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	"github.com/gardener/gardener/pkg/utils/retry"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"
	"time"
)

const (
	TokenRequestorSecretName      = "shoot-access-extension-shoot-fleet-agent"
	TokenRequestorSecretKey       = constsv1alpha1.DataKeyKubeconfig
	TargetServiceAccountName      = "extension-shoot-fleet-agent"
	TargetServiceAccountNamespace = "kube-system"
)

type TokenRequestorHandler struct {
	ctx     context.Context
	client  client.Client
	cluster *extensions.Cluster
	logger  logr.Logger
}

func NewTokenRequestorHandler(ctx context.Context, client client.Client, cluster *extensions.Cluster) *TokenRequestorHandler {
	return &TokenRequestorHandler{
		ctx:     ctx,
		client:  client,
		cluster: cluster,
		logger:  log.Log.WithValues("logger", "token-requestor-handler", "cluster", cluster.ObjectMeta.Name),
	}
}

func (t *TokenRequestorHandler) getGenericKubeconfigName() string {
	return extensionscontroller.GenericTokenKubeconfigSecretNameFromCluster(t.cluster)
}

func (t *TokenRequestorHandler) getGenericKubeconfig() (*clientapiv1.Config, error) {
	kubeconfigSecret := &v1.Secret{}
	kubeconfigSecretName := t.getGenericKubeconfigName()
	if err := t.client.Get(t.ctx, kutil.Key(t.cluster.ObjectMeta.Name, kubeconfigSecretName), kubeconfigSecret); err != nil {
		return nil, err
	}

	kubeconfigData, ok := kubeconfigSecret.Data[TokenRequestorSecretKey]
	if !ok {
		return nil, fmt.Errorf("secret %s doesn't have data key %s", kubeconfigSecretName, TokenRequestorSecretKey)
	}

	kubeconfig := &clientapiv1.Config{}
	if err := yaml.Unmarshal(kubeconfigData, kubeconfig); err != nil {
		return nil, err
	}
	return kubeconfig, nil
}

func (t *TokenRequestorHandler) isTokenInserted() bool {
	kubeconfigSecret := &v1.Secret{}
	if err := t.client.Get(t.ctx, kutil.Key(t.cluster.ObjectMeta.Name, TokenRequestorSecretName), kubeconfigSecret); err != nil {
		t.logger.Info(fmt.Sprintf("Couldn't fet the kubeconfig secret %s, %+v", TokenRequestorSecretName, err))
		return false
	}

	kubeconfigData, ok := kubeconfigSecret.Data[TokenRequestorSecretKey]
	if !ok {
		t.logger.Info(fmt.Sprintf("Kubeconfig secret %s doesn't contain data item %s", TokenRequestorSecretName, TokenRequestorSecretKey))
		return false
	}

	kubeconfig := &clientapiv1.Config{}
	if err := yaml.Unmarshal(kubeconfigData, kubeconfig); err != nil {
		t.logger.Info(fmt.Sprintf("Data item %s in kubeconfig secret %s couldn't be parsed: %+v", TokenRequestorSecretKey, TokenRequestorSecretName, err))
		return false
	}

	if len(kubeconfig.AuthInfos) == 0 || kubeconfig.AuthInfos[0].AuthInfo.Token == "" {
		t.logger.Info(fmt.Sprintf("Kubeconfig in secret %s doesn't contain token (yet?)", TokenRequestorSecretName))
		return false
	}

	t.logger.Info(fmt.Sprintf("Kubeconfig in secret %s contains a proper token", TokenRequestorSecretName))
	return true
}

func (t *TokenRequestorHandler) createKubeconfigSecretFromGeneric() error {
	kubeconfig, err := t.getGenericKubeconfig()
	s := false
	for _, value := range t.cluster.Shoot.Status.AdvertisedAddresses {
		// check for external URL server
		if value.Name == "external" {
			kubeconfig.Clusters[0].Cluster.Server = value.URL
			s = true
			break
		}
	}
	if s != true {
		t.logger.Info(fmt.Sprintf("Shoot status doesn't contain external URL address."))
		return err
	}
	if err != nil {
		return err
	}

	kubeconfigYaml, err := yaml.Marshal(kubeconfig)
	if err != nil {
		return err
	}

	existingSecret := &v1.Secret{}
	err = t.client.Get(t.ctx, kutil.Key(t.cluster.ObjectMeta.Name, TokenRequestorSecretName), existingSecret)
	if err == nil || !apierrors.IsNotFound(err) {
		return err
	}

	secretObject := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: t.cluster.ObjectMeta.Name,
			Name:      TokenRequestorSecretName,
			Labels: map[string]string{
				constsv1alpha1.ResourceManagerPurpose: constsv1alpha1.LabelPurposeTokenRequest,
			},
			Annotations: map[string]string{
				constsv1alpha1.ServiceAccountName:      TargetServiceAccountName,
				constsv1alpha1.ServiceAccountNamespace: TargetServiceAccountNamespace,
			},
		},
		StringData: map[string]string{
			TokenRequestorSecretKey: string(kubeconfigYaml),
		},
	}
	return t.client.Create(t.ctx, secretObject)
}

// EnsureKubeconfig creates the secret with a token-requestor label and waits until gardener fills in
// the token into the kubeconfig
func (t *TokenRequestorHandler) EnsureKubeconfig() error {
	t.logger.Info(fmt.Sprintf("Trying to create the token requestor secret"))
	if err := t.createKubeconfigSecretFromGeneric(); err != nil {
		t.logger.Error(err, "Generic kubeconfig could not be copied")
		return err
	}

	t.logger.Info("Waiting until the token is propagated into the token requestor secret")
	if err := retry.UntilTimeout(t.ctx, 5*time.Second, 60*time.Second, func(ctx context.Context) (bool, error) {
		return t.isTokenInserted(), nil
	}); err != nil {
		t.logger.Error(err, "Kubeconfig could not be created")
		return err
	}

	return nil
}
