package main

import (
	"bufio"
	"os"
	"strings"
)

type Server struct {
	Name    string
	Result  string
	Devices []GPU
}

func GetRemoteServers() ([]string, error) {
	filename := os.Getenv("HOME") + "/.ssh/config"

	// Open the SSH config file.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the hosts.
	var hosts []string
	scanner := bufio.NewScanner(file)
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
