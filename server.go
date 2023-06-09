package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

type Server struct {
	Name    string
	Devices []GPU
}

func ReadDefaultSSHConfig() (io.Reader, error) {
	read, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/config")
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(read), nil
}

func GetHosts(sshConfigReader io.Reader) ([]string, error) {
	// Read the hosts.
	var hosts []string
	scanner := bufio.NewScanner(sshConfigReader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Host ") {
			hosts = append(hosts, strings.TrimSpace(strings.TrimPrefix(line, "Host ")))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func GetServers(info map[string]string) []Server {
	var servers []Server

	for host, gpuInfo := range info {
		server := Server{
			Name:    host,
			Devices: ExtractGPUInfo(strings.TrimSpace(gpuInfo)),
		}
		servers = append(servers, server)
	}

	return servers
}
