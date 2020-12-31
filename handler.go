package backend

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/containerssh/configuration"
	"github.com/containerssh/docker"
	"github.com/containerssh/kubernetes"
	"github.com/containerssh/log"
	"github.com/containerssh/metrics"
	"github.com/containerssh/security"
	"github.com/containerssh/sshserver"
	"github.com/containerssh/structutils"
)

type handler struct {
	config                 configuration.AppConfig
	configLoader           configuration.Loader
	loggerFactory          log.LoggerFactory
	authResponse           sshserver.AuthResponse
	metricsCollector       metrics.Collector
	logger                 log.Logger
	backendRequestsCounter metrics.Counter
	backendErrorCounter    metrics.Counter
}

func (h *handler) OnReady() error {
	return nil
}

func (h *handler) OnShutdown(_ context.Context) {
	//TODO send SIGTERM to containers?
}

func (h *handler) OnNetworkConnection(
	remoteAddr net.TCPAddr,
	connectionID string,
) (sshserver.NetworkConnectionHandler, error) {
	//TODO add early loading for some backends?
	return &networkHandler{
		rootHandler:  h,
		remoteAddr:   remoteAddr,
		connectionID: connectionID,
	}, nil
}

type networkHandler struct {
	rootHandler  *handler
	remoteAddr   net.TCPAddr
	connectionID string
	backend      sshserver.NetworkConnectionHandler
}

func (n *networkHandler) OnAuthPassword(_ string, _ []byte) (response sshserver.AuthResponse, reason error) {
	return n.authResponse()
}

func (n *networkHandler) authResponse() (sshserver.AuthResponse, error) {
	switch n.rootHandler.authResponse {
	case sshserver.AuthResponseUnavailable:
		return sshserver.AuthResponseUnavailable, fmt.Errorf("the backend handler does not support authentication")
	default:
		return n.rootHandler.authResponse, nil
	}
}

func (n *networkHandler) OnAuthPubKey(_ string, _ string) (response sshserver.AuthResponse, reason error) {
	return n.authResponse()
}

func (n *networkHandler) OnHandshakeFailed(_ error) {
}

func (n *networkHandler) OnHandshakeSuccess(username string) (
	connection sshserver.SSHConnectionHandler,
	failureReason error,
) {
	appConfig, err := n.loadConnectionSpecificConfig(username)
	if err != nil {
		return nil, err
	}

	backendLogger, err := n.rootHandler.loggerFactory.Make(appConfig.Log, appConfig.Backend)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger for backend (%w)", err)
	}

	return n.initBackend(username, appConfig, backendLogger)
}

func (n *networkHandler) initBackend(
	username string,
	appConfig configuration.AppConfig,
	backendLogger log.Logger,
) (sshserver.SSHConnectionHandler, error) {
	backend, failureReason := n.getConfiguredBackend(
		appConfig,
		backendLogger,
		n.rootHandler.backendRequestsCounter.WithLabels(metrics.Label(MetricLabelBackend, appConfig.Backend)),
		n.rootHandler.backendErrorCounter.WithLabels(metrics.Label(MetricLabelBackend, appConfig.Backend)),
	)
	if failureReason != nil {
		return nil, failureReason
	}

	// Inject security overlay
	n.backend, failureReason = security.New(appConfig.Security, backend)
	if failureReason != nil {
		return nil, failureReason
	}

	return n.backend.OnHandshakeSuccess(username)
}

func (n *networkHandler) getConfiguredBackend(
	appConfig configuration.AppConfig,
	backendLogger log.Logger,
	backendRequestsCounter metrics.Counter,
	backendErrorCounter metrics.Counter,
) (backend sshserver.NetworkConnectionHandler, failureReason error) {
	switch appConfig.Backend {
	case "docker":
		backend, failureReason = docker.New(
			n.remoteAddr,
			n.connectionID,
			appConfig.Docker,
			backendLogger,
			backendRequestsCounter,
			backendErrorCounter,
		)
	case "dockerrun":
		//goland:noinspection GoDeprecation
		backend, failureReason = docker.NewDockerRun(
			n.remoteAddr,
			n.connectionID,
			appConfig.DockerRun,
			backendLogger,
			backendRequestsCounter,
			backendErrorCounter,
		)
	case "kubernetes":
		backend, failureReason = kubernetes.New(
			n.remoteAddr,
			n.connectionID,
			appConfig.Kubernetes,
			backendLogger,
			backendRequestsCounter,
			backendErrorCounter,
		)
	case "kuberun":
		//goland:noinspection GoDeprecation
		backend, failureReason = kubernetes.NewKubeRun(
			n.remoteAddr,
			n.connectionID,
			appConfig.KubeRun,
			backendLogger,
			backendRequestsCounter,
			backendErrorCounter,
		)
	default:
		failureReason = fmt.Errorf("invalid backend: %s", appConfig.Backend)
	}
	return backend, failureReason
}

func (n *networkHandler) loadConnectionSpecificConfig(
	username string,
) (
	configuration.AppConfig,
	error,
) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFunc()

	appConfig := configuration.AppConfig{}
	if err := structutils.Copy(&appConfig, &n.rootHandler.config); err != nil {
		return appConfig, fmt.Errorf("failed to copy application configuration (%w)", err)
	}

	if err := n.rootHandler.configLoader.LoadConnection(
		ctx,
		username,
		n.remoteAddr,
		n.connectionID,
		&appConfig,
	); err != nil {
		return appConfig, fmt.Errorf("failed to load connections-specific configuration (%w)", err)
	}

	if err := appConfig.Validate(true); err != nil {
		newErr := fmt.Errorf("configuration server returned invalid configuration (%w)", err)
		n.rootHandler.logger.Warninge(newErr)
		return appConfig, newErr
	}

	return appConfig, nil
}

func (n *networkHandler) OnDisconnect() {
	if n.backend != nil {
		n.backend.OnDisconnect()
		n.backend = nil
	}
}
