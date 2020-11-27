package backend

import (
	"context"
	"net"

	"github.com/containerssh/sshserver"
)

type handler struct {
}

func (h *handler) OnReady() error {
	return nil
}

func (h *handler) OnShutdown(shutdownContext context.Context) {

}

func (h *handler) OnNetworkConnection(client net.TCPAddr, connectionID []byte) (sshserver.NetworkConnectionHandler, error) {
	return &networkHandler{}, nil
}

type networkHandler struct {
}

func (n *networkHandler) OnAuthPassword(username string, password []byte) (response sshserver.AuthResponse, reason error) {
	panic("implement me")
}

func (n *networkHandler) OnAuthPubKey(username string, pubKey []byte) (response sshserver.AuthResponse, reason error) {
	panic("implement me")
}

func (n *networkHandler) OnHandshakeFailed(reason error) {
	panic("implement me")
}

func (n *networkHandler) OnHandshakeSuccess() (connection sshserver.SSHConnectionHandler, failureReason error) {
	panic("implement me")
}

func (n *networkHandler) OnDisconnect() {
	panic("implement me")
}
