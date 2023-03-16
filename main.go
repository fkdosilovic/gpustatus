package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type GPUInfo struct {
	Name        string
	Index       string
	TotalMemory string
	FreeMemory  string
	UsedMemory  string
}

type ServerInfo struct {
	Name   string
	Result string
	GPU    []GPUInfo
}

const (
	BASE_COMMAND   = "nvidia-smi"
	DEFAULT_QUERY  = "name,index,memory.total,memory.free,memory.used"
	DEFAULT_FORMAT = "csv"
)

var COMMAND string = fmt.Sprintf("%s --query-gpu=%s --format=%s", BASE_COMMAND, DEFAULT_QUERY, DEFAULT_FORMAT)

func main() {
	// Get remote servers.
	servers := GetRemoveServers()
	info := make(chan ServerInfo, len(servers))

	// Get info from remote servers.
	GetGPUInfoFromServers(servers, info)

	// Extract info from remote servers.
	serverInfo := ProcessGPUInfo(info)

	// Sort by server name.
	sort.Slice(serverInfo, func(i, j int) bool {
		return serverInfo[i].Name < serverInfo[j].Name
	})

	// Print info.
	FormatOutput(serverInfo)
}

func FormatOutput(servers []ServerInfo) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Server", "Name", "Index", "Total Memory", "Used Memory", "Free Memory")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, server := range servers {
		for i, gpu := range server.GPU {
			if i == 0 {
				tbl.AddRow(server.Name, gpu.Name, gpu.Index, gpu.TotalMemory, gpu.UsedMemory, gpu.FreeMemory)
			} else {
				tbl.AddRow("", gpu.Name, gpu.Index, gpu.TotalMemory, gpu.UsedMemory, gpu.FreeMemory)
			}
		}
	}

	tbl.Print()
}

// Get info from remote servers.
func GetGPUInfoFromServers(servers []string, info chan ServerInfo) {
	var wg sync.WaitGroup

	wg.Add(len(servers))
	for _, server := range servers {
		go GetGPUInfo(server, COMMAND, info, &wg)
	}
	wg.Wait()
}

// Run command on remote machine and return the output.
func GetGPUInfo(server string, command string, info chan ServerInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("ssh", server, command)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		info <- ServerInfo{Name: server, Result: ""}
		return
	}

	info <- ServerInfo{Name: server, Result: out.String()}
}

func GetRemoveServers() []string {
	return []string{"jadranka", "adriana", "tuga", "buga", "snjezana", "suncica"}
}

func ExtractGPUInfo(info string) []GPUInfo {
	var gpuInfo []GPUInfo

	lines := strings.Split(info, "\n")
	for _, ln := range lines[1:] {
		fields := strings.Split(ln, ",")
		gpuInfo = append(gpuInfo, GPUInfo{
			Name:        fields[0],
			Index:       fields[1],
			TotalMemory: fields[2],
			FreeMemory:  fields[3],
			UsedMemory:  fields[4],
		})
	}

	return gpuInfo
}

// Process info from remote servers.
func ProcessGPUInfo(info <-chan ServerInfo) []ServerInfo {
	var serverInfo []ServerInfo

	// Populate the serverInfo slice.
	for len(info) > 0 {
		serverInfo = append(serverInfo, <-info)
	}

	for i, server := range serverInfo {
		serverInfo[i].GPU = ExtractGPUInfo(strings.TrimSpace(server.Result))
	}

	return serverInfo
}
