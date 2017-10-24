package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
)

type EmptyStruct struct{}

type Buildstruct struct {
	// remember to use caps so that they can be exported
	Context    string `yaml:"context,omitempty"`
	Dockerfile string `yaml:"dockerfile,omitempty"`
}

type ComposeService struct {
	Image       string      `yaml:"image,omitempty"`
	Ports       []string    `yaml:"ports,omitempty"`
	Labels      []string    `yaml:"labels,omitempty"`
	Environment []string    `yaml:"environment,omitempty"`
	Command     string      `yaml:"command,flow,omitempty"`
	Restart     string      `yaml:"restart,omitempty"`
	Build       Buildstruct `yaml:"build,omitempty"`
	Volumes     []string    `yaml:"volumes,omitempty"`
	Links       []string    `yaml:"links,omitempty"`
}

type DockerComposeConfig struct {
	Version  string                    `yaml:"version,omitempty"`
	Services map[string]ComposeService `yaml:"services"`
	Volumes  map[string]string         `yaml:"volumes,omitempty"`
}

type DockerContainer struct {
	ServiceName       string
	ComposeService    ComposeService
	NetworkID         string
	NetworkName       string
	FollowLogs        bool
	DockerComposeFile string
	ContainerID       string
	// this assumes that there can only be one container per docker-compose service
	LogMedium io.Writer
	Color     string
}

func (dc *DockerContainer) UpdateContainerID(containerID string) {
	dc.ContainerID = containerID
}

var Colors = []string{
	"\x1b[30;1m", // black
	"\x1b[31;1m", // red
	"\x1b[32;1m", // green
	"\x1b[33;1m", // yellow
	"\x1b[34;1m", // blue
	"\x1b[35;1m", // magenta
	"\x1b[36;1m", // cyan
	"\x1b[37;1m"} // white

type MeliAPiClient interface {
	// we implement this interface so that we can be able to mock it in tests
	// https://medium.com/@zach_4342/dependency-injection-in-golang-e587c69478a8
	ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error)
	NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error)
	NetworkConnect(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error
	VolumeCreate(ctx context.Context, options volumetypes.VolumesCreateBody) (types.Volume, error)
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
}

type MockDockerClient struct{}

func (m *MockDockerClient) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer([]byte("Pulling from library/testImage"))), nil
}
func (m *MockDockerClient) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{Body: ioutil.NopCloser(bytes.NewBuffer([]byte("BUILT library/testImage"))), OSType: "linux baby!"}, nil
}
func (m *MockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	return container.ContainerCreateCreatedBody{ID: "myContainerId001"}, nil
}

func (m *MockDockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return nil
}

func (m *MockDockerClient) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer([]byte("SHOWING LOGS for library/testImage"))), nil
}

func (m *MockDockerClient) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
	return []types.NetworkResource{}, nil
}
func (m *MockDockerClient) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	return types.NetworkCreateResponse{ID: "myNetworkId002"}, nil
}

func (m *MockDockerClient) NetworkConnect(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error {
	return nil
}

func (m *MockDockerClient) VolumeCreate(ctx context.Context, options volumetypes.VolumesCreateBody) (types.Volume, error) {
	return types.Volume{Name: "MyVolume007"}, nil
}

func (m *MockDockerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return []types.Container{types.Container{ID: "myExistingContainerId00912"}}, nil
}

func CopyBufferWithColor(dst io.Writer, src io.Reader, buf []byte, serviceName, color string) (written int64, err error) {
	// undo the set TERM color
	defer fmt.Println("\x1b[0m")

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			// also set TERM color
			fmt.Fprintf(dst, "%sSERVICE=%s:: ", color, serviceName)
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}

	return written, err
}
