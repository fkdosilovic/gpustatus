package main

// Returns a list of remote servers containing free GPUs.
func FilterFreeGPUs(servers []Server) []Server {
	return filterGPUs(servers, CheckIsGPUFree)
}

func FilterUsedGPUs(servers []Server) []Server {
	return filterGPUs(servers, func(gpu GPU) bool {
		return !CheckIsGPUFree(gpu)
	})
}

func filterGPUs(servers []Server, f func(GPU) bool) []Server {
	var filtered []Server

	for _, server := range servers {
		var devices []GPU

		for _, gpu := range server.Devices {
			if f(gpu) {
				devices = append(devices, gpu)
			}
		}

		filtered = append(filtered, Server{Name: server.Name, Devices: devices})
	}

	return filtered
}
