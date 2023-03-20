package main

import (
	"sort"
)

func main() {
	args := ParseArguments()

	// Get remote servers.
	servers, _ := GetRemoteServers()
	info := make(chan Server, len(servers))

	// Get info from remote servers.
	GetGPUInfoFromServers(servers, info)

	// Extract info from remote servers.
	serverInfo := ProcessGPUInfo(info)

	// Filter GPUs.
	if args.ShowFree {
		serverInfo = FilterFreeGPUs(serverInfo)
	} else if args.ShowUsed {
		serverInfo = FilterUsedGPUs(serverInfo)
	}

	// Sort by server name.
	sort.Slice(serverInfo, func(i, j int) bool {
		return serverInfo[i].Name < serverInfo[j].Name
	})

	// Print info.
	tbl := CreateOutput(serverInfo)
	tbl.Print()
}
