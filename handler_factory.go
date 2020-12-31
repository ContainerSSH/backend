package backend

import (
	"github.com/containerssh/configuration"
	"github.com/containerssh/log"
	"github.com/containerssh/metrics"
	"github.com/containerssh/sshserver"
)

// New creates a new backend handler.
//goland:noinspection GoUnusedExportedFunction
func New(
	config configuration.AppConfig,
	logger log.Logger,
	loggerFactory log.LoggerFactory,
	metricsCollector metrics.Collector,
	defaultAuthResponse sshserver.AuthResponse,
) (sshserver.Handler, error) {
	loader, err := configuration.NewHTTPLoader(
		config.ConfigServer,
		logger,
		metricsCollector,
	)
	if err != nil {
		return nil, err
	}

	backendRequestsCounter := metricsCollector.MustCreateCounter(
		MetricNameBackendRequests,
		MetricUnitBackendRequests,
		MetricHelpBackendRequests,
	)
	backendErrorCounter := metricsCollector.MustCreateCounter(
		MetricNameBackendError,
		MetricUnitBackendError,
		MetricHelpBackendError,
	)

	return &handler{
		config:                 config,
		configLoader:           loader,
		loggerFactory:          loggerFactory,
		authResponse:           defaultAuthResponse,
		metricsCollector:       metricsCollector,
		logger:                 logger,
		backendRequestsCounter: backendRequestsCounter,
		backendErrorCounter:    backendErrorCounter,
	}, nil
}
