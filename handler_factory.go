package backend

import (
	"github.com/containerssh/sshserver"
)

// New creates a new backend handler.
//goland:noinspection GoUnusedExportedFunction
func New() (sshserver.Handler, error) {
	return &handler{}, nil
}
