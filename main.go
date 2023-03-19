package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type GPU struct {
	Name        string
	Index       string
	FreeMemory  string
	UsedMemory  string
	TotalMemory string
}

type Server struct {
	Name    string
	Result  string
	Devices []GPU
}

const (
	BASE_COMMAND   = "nvidia-smi"
	DEFAULT_QUERY  = "name,index,memory.free,memory.used,memory.total"
	DEFAULT_FORMAT = "csv"
)

var COMMAND string = fmt.Sprintf("%s --query-gpu=%s --format=%s", BASE_COMMAND, DEFAULT_QUERY, DEFAULT_FORMAT)

func main() {
	// Get remote servers.
	servers, _ := GetRemoteServers()
	info := make(chan Server, len(servers))

	// Get info from remote servers.
	GetGPUInfoFromServers(servers, info)

	// Extract info from remote servers.
	serverInfo := ProcessGPUInfo(info)

	// Sort by server name.
	sort.Slice(serverInfo, func(i, j int) bool {
		return serverInfo[i].Name < serverInfo[j].Name
	})

	// Print info.
	tbl := CreateOutput(serverInfo)
	tbl.Print()
}

func CreateOutput(servers []Server) table.Table {
	headerFmt := color.New(color.FgWhite, color.Bold, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgHiBlue).SprintfFunc()

	tbl := table.New("Server", "Name", "Index", "Free Memory", "Used Memory", "Total Memory")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, server := range servers {
		serverName := server.Name
		for _, gpu := range server.Devices {
			tbl.AddRow(serverName, gpu.Name, gpu.Index, AlignRight(gpu.FreeMemory), AlignRight(gpu.UsedMemory), AlignRight(gpu.TotalMemory))
			serverName = ""
		}
	}

	return tbl
}

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

func AlignRight(s string) string {
	return fmt.Sprintf("%10s", s)
}
