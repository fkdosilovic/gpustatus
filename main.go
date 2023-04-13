package main

import (
	"log"
	"sort"
)

func main() {
	args := ParseArguments()

	// Read SSH config.
	sshConfigReader, err := ReadDefaultSSHConfig()
	if err != nil {
		log.Fatalf("Failed to read SSH config: %v", err)
	}

	// Get remote hosts.
	hosts, _ := GetHosts(sshConfigReader)

	// Get info from remote hosts.
	info := GetGPUInfoFromHosts(hosts)

	// Extract info from the query result.
	servers := GetServers(info)

	// Filter GPUs.
	if args.ShowFree {
		servers = FilterFreeGPUs(servers)
	} else if args.ShowUsed {
		servers = FilterUsedGPUs(servers)
	}

	// Sort by server name.
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Name < servers[j].Name
	})

	// Print info.
	tbl := CreateOutput(servers)
	tbl.Print()
}
