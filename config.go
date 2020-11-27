package backend

import (
	"errors"
)

// Name is an enum for supported backend names.
// swagger:enum BackendName
type Name string

const (
	// NameDockerRun is the Docker backend.
	NameDockerRun Name = "dockerrun"
	// NameKubeRun is the Kubernetes backend.
	NameKubeRun Name = "kuberun"
)

// ErrInvalidBackend is returned if the backend returned is invalid.
var ErrInvalidBackend = errors.New("invalid backend")

// Validate returns an error if the provided name is invalid.
func (n Name) Validate() error {
	switch n {
	case NameDockerRun:
		return nil
	case NameKubeRun:
		return nil
	default:
		return ErrInvalidBackend
	}
}
