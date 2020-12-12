package backend

import (
	"github.com/containerssh/configuration"
	"github.com/containerssh/log"
	"github.com/containerssh/sshserver"
)

// New creates a new backend handler.
//goland:noinspection GoUnusedExportedFunction
func New(
	config configuration.AppConfig,
	logger log.Logger,
	loggerFactory log.LoggerFactory,
	defaultAuthResponse sshserver.AuthResponse,
) (sshserver.Handler, error) {
	loader, err := configuration.NewHTTPLoader(
		config.ConfigServer,
		logger,
	)
	if err != nil {
		return nil, err
	}
	return &handler{
		configLoader:  loader,
		config:        config,
		loggerFactory: loggerFactory,
		authResponse:  defaultAuthResponse,
	}, nil
}
