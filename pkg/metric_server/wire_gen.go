// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package metric_server

import (
	"github.com/coreos/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/pkg/client"
	"github.com/redhat-marketplace/redhat-marketplace-operator/pkg/controller"
	"github.com/redhat-marketplace/redhat-marketplace-operator/pkg/generated/clientset/versioned/typed/marketplace/v1alpha1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/pkg/managers"
	"github.com/redhat-marketplace/redhat-marketplace-operator/pkg/meter_definition"
	"github.com/redhat-marketplace/redhat-marketplace-operator/pkg/utils/reconcileutils"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Injectors from wire.go:

func NewServer(opts *Options) (*Service, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	restMapper, err := managers.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}
	opsSrcSchemeDefinition := controller.ProvideOpsSrcScheme()
	monitoringSchemeDefinition := controller.ProvideMonitoringScheme()
	olmV1SchemeDefinition := controller.ProvideOLMV1Scheme()
	olmV1Alpha1SchemeDefinition := controller.ProvideOLMV1Alpha1Scheme()
	openshiftConfigV1SchemeDefinition := controller.ProvideOpenshiftConfigV1Scheme()
	localSchemes := controller.ProvideLocalSchemes(opsSrcSchemeDefinition, monitoringSchemeDefinition, olmV1SchemeDefinition, olmV1Alpha1SchemeDefinition, openshiftConfigV1SchemeDefinition)
	scheme, err := managers.ProvideScheme(restConfig, localSchemes)
	if err != nil {
		return nil, err
	}
	clientOptions := getClientOptions()
	cache, err := managers.ProvideNewCache(restConfig, restMapper, scheme, clientOptions)
	if err != nil {
		return nil, err
	}
	clientClient, err := managers.ProvideClient(restConfig, restMapper, scheme, cache, clientOptions)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	options := ConvertOptions(opts)
	registry := provideRegistry()
	logger := _wireLoggerValue
	clientCommandRunner := reconcileutils.NewClientCommand(clientClient, scheme, logger)
	context := provideContext()
	dynamicInterface, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	findOwnerHelper := client.NewFindOwnerHelper(dynamicInterface, restMapper)
	monitoringV1Client, err := v1.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	marketplaceV1alpha1Client, err := v1alpha1.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	meterDefinitionStore := meter_definition.NewMeterDefinitionStore(context, logger, clientCommandRunner, clientset, findOwnerHelper, monitoringV1Client, marketplaceV1alpha1Client, scheme)
	statusProcessor := meter_definition.NewStatusProcessor(logger, clientCommandRunner, meterDefinitionStore)
	cacheIsIndexed, err := addIndex(context, cache)
	if err != nil {
		return nil, err
	}
	cacheIsStarted := managers.StartCache(context, cache, logger, cacheIsIndexed)
	service := &Service{
		k8sclient:       clientClient,
		k8sRestClient:   clientset,
		opts:            options,
		cache:           cache,
		metricsRegistry: registry,
		cc:              clientCommandRunner,
		meterDefStore:   meterDefinitionStore,
		statusProcessor: statusProcessor,
		isCacheStarted:  cacheIsStarted,
	}
	return service, nil
}

var (
	_wireLoggerValue = log
)
