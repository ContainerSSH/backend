package backend_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/containerssh/configuration"
	"github.com/containerssh/log"
	"github.com/containerssh/service"
	"github.com/containerssh/sshserver"
	"github.com/containerssh/structutils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"

	"github.com/containerssh/backend"
)

func TestSimpleContainerLaunch(t *testing.T) {
	config := configuration.AppConfig{}
	structutils.Defaults(&config)
	err := config.SSH.GenerateHostKey()
	assert.NoError(t, err)

	loggerFactory := log.NewFactory(os.Stdout)

	backendLogger, err := loggerFactory.Make(config.Log, "backend")
	assert.NoError(t, err)
	b, err := backend.New(
		config,
		backendLogger,
		loggerFactory,
		sshserver.AuthResponseSuccess,
	)
	assert.NoError(t, err)

	sshServerLogger, err := loggerFactory.Make(config.Log, "ssh")
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	defer func() {
		if sshConnection != nil {
			_ = sshConnection.Close()
		}
	}()

	session, err := sshConnection.NewSession()
	assert.NoError(t, err)

	_, stdout, err := createPipe(session)
	assert.NoError(t, err)

	err = session.Start("echo 'Hello world!'")
	assert.NoError(t, err)
	if err != nil && !errors.Is(err, io.EOF) {
		assert.NoError(t, err)
	}
	output, err := ioutil.ReadAll(stdout)
	assert.NoError(t, err)

	assert.NoError(t, session.Wait())

	assert.NoError(t, sshConnection.Close())
	assert.EqualValues(t, []byte("Hello world!\n"), output)
}

func createPipe(session *ssh.Session) (io.WriteCloser, io.Reader, error) {
	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to request stdin (%w)", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to request stdout (%w)", err)
	}
	return stdin, stdout, nil
}
