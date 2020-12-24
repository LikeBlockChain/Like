package fakes

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net"

	"github.com/ethereum/hive/internal/hive"
)

// BackendHooks can be used to override the behavior of the fake backend.
type BackendHooks struct {
	StartClient   func(name string, env map[string]string) (*hive.ClientInfo, error)
	StopContainer func(string) error
	RunEnodeSh    func(string) (string, error)

	NetworkNameToID     func(string) (string, error)
	CreateNetwork       func(string) (string, error)
	RemoveNetwork       func(networkID string) error
	ContainerIP         func(containerID, networkID string) (net.IP, error)
	ConnectContainer    func(containerID, networkID string) error
	DisconnectContainer func(containerID, networkID string) error
}

// fakeBackend implements Backend without docker.
type fakeBackend struct {
	hooks         BackendHooks
	clientCounter uint64
	netCounter    uint64
}

// NewBackend creates a new fake container backend.
func NewBackend(hooks *BackendHooks) hive.Backend {
	b := &fakeBackend{}
	if hooks != nil {
		b.hooks = *hooks
	}
	return b
}

func (b *fakeBackend) StartClient(name string, env map[string]string, files map[string]*multipart.FileHeader, checklive bool) (*hive.ClientInfo, error) {
	var info hive.ClientInfo
	if b.hooks.StartClient != nil {
		info2, err := b.hooks.StartClient(name, env)
		if err != nil {
			return nil, err
		}
		info = *info2
	}

	b.clientCounter++
	info.ID = fmt.Sprintf("%0.8x", b.clientCounter)
	if info.IP == "" {
		ip := net.IP{192, 0, 2, byte(b.clientCounter)}
		info.IP = ip.String()
	}
	if info.MAC == "" {
		info.MAC = "00:80:41:ae:fd:7e"
	}
	return &info, nil
}

func (b *fakeBackend) StopContainer(containerID string) error {
	if b.hooks.StopContainer != nil {
		return b.hooks.StopContainer(containerID)
	}
	return nil
}

func (b *fakeBackend) RunEnodeSh(containerID string) (string, error) {
	if b.hooks.RunEnodeSh != nil {
		return b.hooks.RunEnodeSh(containerID)
	}
	return "enode://a61215641fb8714a373c80edbfa0ea8878243193f57c96eeb44d0bc019ef295abd4e044fd619bfc4c59731a73fb79afe84e9ab6da0c743ceb479cbb6d263fa91@192.0.2.1:30303", nil
}

func (b *fakeBackend) NetworkNameToID(name string) (string, error) {
	if b.hooks.NetworkNameToID != nil {
		return b.hooks.NetworkNameToID(name)
	}
	return "", errors.New("network not found")
}

func (b *fakeBackend) CreateNetwork(name string) (string, error) {
	if b.hooks.CreateNetwork != nil {
		return b.hooks.CreateNetwork(name)
	}
	b.netCounter++
	id := fmt.Sprintf("%0.8x", b.netCounter)
	return id, nil
}

func (b *fakeBackend) RemoveNetwork(networkID string) error {
	if b.hooks.RemoveNetwork != nil {
		return b.hooks.RemoveNetwork(networkID)
	}
	return nil
}

func (b *fakeBackend) ContainerIP(containerID, networkID string) (net.IP, error) {
	if b.hooks.ContainerIP != nil {
		return b.hooks.ContainerIP(containerID, networkID)
	}
	return net.IP{203, 0, 113, 2}, nil
}

func (b *fakeBackend) ConnectContainer(containerID, networkID string) error {
	if b.hooks.ConnectContainer != nil {
		return b.hooks.ConnectContainer(containerID, networkID)
	}
	return nil
}

func (b *fakeBackend) DisconnectContainer(containerID, networkID string) error {
	if b.hooks.DisconnectContainer != nil {
		return b.hooks.DisconnectContainer(containerID, networkID)
	}
	return nil
}
