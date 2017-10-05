package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

/* DOCS:
1. https://godoc.org/github.com/moby/moby/client
2. https://docs.docker.com/engine/api/v1.31/
*/

type emptyStruct struct{}

type buildstruct struct {
	// remember to use caps so that they can be exported
	Context    string `yaml:"context,omitempty"`
	Dockerfile string `yaml:"dockerfile,omitempty"`
}

type serviceConfig struct {
	Image       string      `yaml:"image,omitempty"`
	Ports       []string    `yaml:"ports,omitempty"`
	Labels      []string    `yaml:"labels,omitempty"`
	Environment []string    `yaml:"environment,omitempty"`
	Command     string      `yaml:"command,flow,omitempty"`
	Restart     string      `yaml:"restart,omitempty"`
	Build       buildstruct `yaml:"build,omitempty"`
	Volumes     []string    `yaml:"volumes,omitempty"`
}

type dockerComposeConfig struct {
	Version  string                   `yaml:"version,omitempty"`
	Services map[string]serviceConfig `yaml:"services"`
	Volumes  map[string]string        `yaml:"volumes,omitempty"`
}

func main() {
	data, err := ioutil.ReadFile("docker-compose.yml")
	if err != nil {
		log.Fatal(err, "unable to read docker-compose file")
	}
	curentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err, "unable to get the current working directory")
	}
	networkName := "meli_network_" + getCwdName(curentDir)
	networkID, err := getNetwork(networkName)
	if err != nil {
		log.Fatal(err, "unable to create/get network")
	}

	var dockerCyaml dockerComposeConfig
	err = yaml.Unmarshal([]byte(data), &dockerCyaml)
	if err != nil {
		log.Fatal(err, "unable to parse docker-compose file contents")
	}

	ctx := context.Background()

	// Create top level volumes, if any
	if len(dockerCyaml.Volumes) > 0 {
		fmt.Println("len", len(dockerCyaml.Volumes))
		for k := range dockerCyaml.Volumes {
			// TODO we need to synchronise here else we'll get a race
			// but I think we can get away for now because:
			// 1. there are on average a lot more containers in a compose file
			// than volumes, so the sync in the for loop for containers is enough
			// 2. since we intend to stream logs as containers run(see; issues/24);
			// then meli will be up long enough for the volume creation goroutines to have finished.
			go CreateDockerVolume(ctx, "meli_"+k, "local")
		}
	}

	var wg sync.WaitGroup
	for _, v := range dockerCyaml.Services {
		wg.Add(1)
		fmt.Println("docker service", v)
		//go fakepullImage(ctx, v, networkID, networkName, &wg)
		go pullImage(ctx, v, networkID, networkName, &wg)
	}
	wg.Wait()
}

func fakepullImage(ctx context.Context, s serviceConfig, networkName, networkID string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println()
}

func pullImage(ctx context.Context, s serviceConfig, networkID, networkName string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 1. Pull Image
	formattedImageName := fomatImageName("containerFromBuild")
	if len(s.Image) > 0 {
		formattedImageName = fomatImageName(s.Image)
		PullDockerImage(ctx, s.Image)
	}

	// 2. Create a container
	containerCreateResp := CreateContainer(ctx, s, networkName, formattedImageName)

	// 3. Connect container to network
	networkConnect(ctx, networkID, containerCreateResp.ID)

	// 4. Start container
	ContainerStart(ctx, containerCreateResp.ID)

	// 5. Stream container logs to stdOut
	ContainerLogs(ctx, containerCreateResp.ID)
}
