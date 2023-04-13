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
func GetGPUInfoFromHosts(hosts []string) map[string]string {
	var wg sync.WaitGroup
	var m sync.Map

	wg.Add(len(hosts))
	for _, host := range hosts {
		go GetGPUInfo(host, COMMAND, &m, &wg)
	}
	wg.Wait()

	var info = make(map[string]string)
	m.Range(func(key, value interface{}) bool {
		info[key.(string)] = value.(string)
		return true
	})

	return info
}

// Run command on remote machine and return the output.
func GetGPUInfo(host string, command string, m *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("ssh", host, command)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		m.Store(host, "")
		return
	}

	m.Store(host, out.String())
}

func ExtractGPUInfo(gpuInfo string) []GPU {
	var devices []GPU

	lines := strings.Split(gpuInfo, "\n")
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

func CheckIsGPUFree(gpu GPU) bool {
	usedMemory, _ := GetMemoryInMB(gpu.UsedMemory)
	totalMemory, _ := GetMemoryInMB(gpu.TotalMemory)
	return float64(usedMemory)/float64(totalMemory) <= FREE_MEMORY_PERCENTAGE
}
