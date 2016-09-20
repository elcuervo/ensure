package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

type multipleValues []string

func (m *multipleValues) String() string {
	return fmt.Sprint(*m)
}

func (m *multipleValues) Set(value string) error {
	*m = append(*m, value)

	return nil
}

var socket = flag.String("socket", "/var/run/docker.sock", "Docker socket")
var name = flag.String("name", "", "Docker Service name")
var image = flag.String("image", "", "Docker Service image")
var replicas = flag.Uint64("replicas", 1, "Docker Service name")

var portFlag multipleValues
var envFlag multipleValues

func init() {
	log.SetLevel(log.DebugLevel)

	flag.Var(&portFlag, "port", "Ports to bind the service")
	flag.Var(&envFlag, "env", "Env to bind the service")
}

func main() {
	flag.Parse()

	if _, err := os.Stat(*socket); os.IsNotExist(err) {
		log.Fatal("I need access to the Docker socket to be able to work.")
	}

	if *name == "" {
		log.Fatal("--name is a mandatory flag")
	}

	log.Debug("Name: ", *name)
	log.Debug("image: ", *image)
	log.Debug("Replicas: ", *replicas)

	for _, p := range portFlag {
		log.Debug("Port: ", p)
	}

	for _, p := range envFlag {
		log.Debug("Env: ", p)
	}

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix://"+*socket, "v1.22", nil, defaultHeaders)

	if err != nil {
		panic(err)
	}

	service, _, err := cli.ServiceInspectWithRaw(context.Background(), *name)

	if err != nil {
		log.Info("Service " + *name + " does not exist, creating.")
		containerSpec := swarm.ContainerSpec{
			Image: *image,
			Env:   envFlag,
		}

		portsConfig := []swarm.PortConfig{}

		for _, port := range portFlag {
			parts := strings.FieldsFunc(port, func(c rune) bool {
				// 58 == : || 47 == /
				return c == 58 || c == 47
			})

			if len(parts) == 0 {
				return
			}

			pub, _ := strconv.ParseUint(parts[0], 10, 32)
			target, _ := strconv.ParseUint(parts[1], 10, 32)
			proto := swarm.PortConfigProtocolTCP

			if len(parts) == 3 {
				proto = swarm.PortConfigProtocolUDP
			}

			portConfig := swarm.PortConfig{
				Protocol:      proto,
				TargetPort:    uint32(target),
				PublishedPort: uint32(pub),
			}
			portsConfig = append(portsConfig, portConfig)

		}

		endpointSpec := &swarm.EndpointSpec{Ports: portsConfig}
		taskSpec := swarm.TaskSpec{ContainerSpec: containerSpec}
		replicatedService := &swarm.ReplicatedService{Replicas: replicas}
		serviceMode := swarm.ServiceMode{Replicated: replicatedService}

		spec := swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: *name,
			},
			Mode:         serviceMode,
			TaskTemplate: taskSpec,
			EndpointSpec: endpointSpec,
		}

		cli.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{})
	} else {
		log.Info("Service " + *name + " exists, checking for differences")
		log.Printf("%+v\n", service)

		spec := service.Spec
		serviceImage := spec.TaskTemplate.ContainerSpec.Image
		serviceReplicas := spec.Mode.Replicated.Replicas

		if serviceImage != *image {
			log.Debug("Service Image does not match")
		} else {
			log.Debug("Service Image matches, not updating")
		}

		if int(*serviceReplicas) != int(*replicas) {
			log.Debug("Service Replicas does not match.")
		} else {
			log.Debug("Service Replicas matches, not updating")
		}

		for _, port := range portFlag {
			log.Info(port)
		}

	}

}
