package main

import (
	"errors"
	"io/ioutil"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Blacklist map[string]struct{}
	Protocols map[string]map[int]string
	Host      string
	Port      int
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	if path != "" {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return nil, err
		}
	}

	if config.Host == "" {
		host, err := GetLocalName()
		if err != nil {
			return nil, err
		}
		config.Host = host
	}

	if config.Port == 0 {
		config.Port = 8000
	}

	return config, nil
}

func (c *Config) String() string {
	data, err := yaml.Marshal(&c)
	if err != nil {
		return "error"
	}
	return string(data)
}

func GetLocalName() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		host, err = GetFirstIP()
		if err != nil {
			return "", errors.New("Could not find a name or IP for this machine")
		}
	}
	return host, nil
}

func GetFirstIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("No IPv4 IPs found")
}
