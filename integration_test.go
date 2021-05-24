package backend_test

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/containerssh/configuration/v2"
	"github.com/containerssh/geoip"
	"github.com/containerssh/log"
	"github.com/containerssh/metrics"
	"github.com/containerssh/service"
	"github.com/containerssh/sshserver"
	"github.com/containerssh/structutils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"

	"github.com/containerssh/backend"
)

func TestSimpleContainerLaunch(t *testing.T) {
	t.Parallel()

	lock := &sync.Mutex{}
	for _, backendName := range []string{"docker", "dockerrun"} {
		t.Run("backend="+backendName, func(t *testing.T) {
			lock.Lock()
			defer lock.Unlock()
			config := configuration.AppConfig{}
			structutils.Defaults(&config)
			config.Backend = backendName
			config.Auth.URL = "http://localhost:8080"
			err := config.SSH.GenerateHostKey()
			assert.NoError(t, err)

			backendLogger := log.NewTestLogger(t)
			geoIPLookupProvider, err := geoip.New(
				geoip.Config{
					Provider: geoip.DummyProvider,
				},
			)
			assert.NoError(t, err)
			metricsCollector := metrics.New(
				geoIPLookupProvider,
			)
			b, err := backend.New(
				config,
				backendLogger,
				metricsCollector,
				sshserver.AuthResponseSuccess,
			)
			assert.NoError(t, err)

			sshServerLogger := log.NewTestLogger(t)
			sshServer, err := sshserver.New(config.SSH, b, sshServerLogger)
			assert.NoError(t, err)

			lifecycle := service.NewLifecycle(sshServer)
			running := make(chan struct{})
			lifecycle.OnRunning(
				func(s service.Service, l service.Lifecycle) {
					running <- struct{}{}
				})
			go func() {
				_ = lifecycle.Run()
			}()
			<-running

			processClientInteraction(t, config)

			lifecycle.Stop(context.Background())
			err = lifecycle.Wait()
			assert.NoError(t, err)
		})
	}
}

func processClientInteraction(t *testing.T, config configuration.AppConfig) {
	clientConfig := &ssh.ClientConfig{
		User: "foo",
		Auth: []ssh.AuthMethod{ssh.Password("bar")},
	}
	clientConfig.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	sshConnection, err := ssh.Dial("tcp", config.SSH.Listen, clientConfig)
	if !assert.NoError(t, err) {
		return
	}
	defer func() {
		if sshConnection != nil {
			_ = sshConnection.Close()
		}
	}()

	session, err := sshConnection.NewSession()
	assert.NoError(t, err)

	output, err := session.CombinedOutput("echo 'Hello world!'")
	assert.NoError(t, err)

	assert.NoError(t, sshConnection.Close())
	assert.EqualValues(t, []byte("Hello world!\n"), output)
}
