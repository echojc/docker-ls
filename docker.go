package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
)

type Container struct {
	Id     string
	Names  []string
	Image  string
	Ports  []*Port
	Status string
}

type Port struct {
	IP          string
	PrivatePort int
	PublicPort  int
	Type        string
}

const (
	Socket   = "/var/run/docker.sock"
	Protocol = "unix"
)

func Containers() ([]*Container, error) {
	req, err := http.NewRequest("GET", "/containers/json", nil)
	if err != nil {
		return nil, err
	}

	dial, err := net.Dial(Protocol, Socket)
	if err != nil {
		return nil, err
	}

	conn := httputil.NewClientConn(dial, nil)
	defer conn.Close()

	res, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Got status %d"))
	}

	var containers []*Container
	err = json.NewDecoder(res.Body).Decode(&containers)
	return containers, err
}

func (c *Container) PublicPorts() []int {
	ports := []int{}
	for _, port := range c.Ports {
		if port.PublicPort > 0 {
			ports = append(ports, port.PublicPort)
		}
	}
	return ports
}
