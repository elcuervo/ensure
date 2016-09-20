package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
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
var replicas = flag.Int("replicas", 1, "Docker Service name")

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
	cli.ServiceInspectWithRaw(context.Background(), *name)
}
