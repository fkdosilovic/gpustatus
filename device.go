package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type GPU struct {
	Name        string
	Index       string
	FreeMemory  string
	UsedMemory  string
	TotalMemory string
}

const (
	BASE_COMMAND   = "nvidia-smi"
	DEFAULT_QUERY  = "name,index,memory.free,memory.used,memory.total"
	DEFAULT_FORMAT = "csv"
)

var COMMAND string = fmt.Sprintf("%s --query-gpu=%s --format=%s", BASE_COMMAND, DEFAULT_QUERY, DEFAULT_FORMAT)

// We consider a GPU to be free if it has at most 5% of its memory used.
const FREE_MEMORY_PERCENTAGE = 0.05

// Get info from remote servers.
func GetGPUInfoFromServers(servers []string, info chan Server) {
	var wg sync.WaitGroup

	wg.Add(len(servers))
	for _, server := range servers {
		go GetGPUInfo(server, COMMAND, info, &wg)
	}
	wg.Wait()
}

// Run command on remote machine and return the output.
func GetGPUInfo(server string, command string, info chan Server, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("ssh", server, command)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		info <- Server{Name: server, Result: ""}
		return
	}

	info <- Server{Name: server, Result: out.String()}
}

func ExtractGPUInfo(info string) []GPU {
	var devices []GPU

	lines := strings.Split(info, "\n")
	for _, ln := range lines[1:] {
		fields := strings.Split(ln, ",")
		devices = append(devices, GPU{
			Name:        fields[0],
			Index:       fields[1],
			FreeMemory:  fields[2],
			UsedMemory:  fields[3],
			TotalMemory: fields[4],
		})
	}

	return devices
}

// Process info from remote servers.
func ProcessGPUInfo(info <-chan Server) []Server {
	var serverInfo []Server

	for len(info) > 0 {
		server := <-info
		server.Devices = ExtractGPUInfo(strings.TrimSpace(server.Result))
		serverInfo = append(serverInfo, server)
	}

	return serverInfo
}

func CheckIsGPUFree(gpu GPU) bool {
	usedMemory, _ := GetMemoryInMB(gpu.UsedMemory)
	totalMemory, _ := GetMemoryInMB(gpu.TotalMemory)
	return float64(usedMemory)/float64(totalMemory) <= FREE_MEMORY_PERCENTAGE
}
