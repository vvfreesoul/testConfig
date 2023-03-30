package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"regexp"
)

///Users/guoweizhang/GolandProjects/testConfig/main.go

type HostConfig struct {
	Hostname string
	IP      string
	CPU     string
	Mem     string
	Storage string
}

func FindConfigByHostname(hostname string) (*HostConfig, error) {
	// Resolve hostname to IP address
	ipAddr, err := net.LookupIP(hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve IP address: %v", err)
	}

	// Convert IP address to string
	ip := ipAddr[0].String()

	// Run command to retrieve CPU information
	cpuCmd := exec.Command("cat", "/proc/cpuinfo")
	cpuOut, err := cpuCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve CPU information: %v", err)
	}

	// Parse CPU information
	cpuRegex := regexp.MustCompile(`model name\s+:\s+(.*)\n`)
	cpuMatch := cpuRegex.FindStringSubmatch(string(cpuOut))
	cpu := ""
	if len(cpuMatch) > 1 {
		cpu = cpuMatch[1]
	}

	// Run command to retrieve memory information
	memCmd := exec.Command("cat", "/proc/meminfo")
	memOut, err := memCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve memory information: %v", err)
	}

	// Parse memory information
	memRegex := regexp.MustCompile(`MemTotal:\s+(\d+) kB\n`)
	memMatch := memRegex.FindStringSubmatch(string(memOut))
	mem := ""
	if len(memMatch) > 1 {
		memBytes := 1024 * 1024 * atoi(memMatch[1])
		mem = fmt.Sprintf("%d GB", memBytes/(1024*1024*1024))
	}

	// Run command to retrieve storage information
	storageCmd := exec.Command("df", "-h", "/")
	storageOut, err := storageCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve storage information: %v", err)
	}

	// Parse storage information
	storageRegex := regexp.MustCompile(`/dev/.*?\s+(\d+G)\s+(\d+G)\s+(\d+G)\s+(\d+)%\s+/`)
	storageMatch := storageRegex.FindStringSubmatch(string(storageOut))
	storage := ""
	if len(storageMatch) > 4 {
		storage = fmt.Sprintf("%s used (%s/%s total)", storageMatch[4], storageMatch[2], storageMatch[1])
	}

	// Construct HostConfig object
	config := &HostConfig{
		Hostname: hostname,
		IP:      ip,
		CPU:     cpu,
		Mem:     mem,
		Storage: storage,
	}

	return config, nil
}

func main() {
	configs := make([]*HostConfig, 0)

	config, err := FindConfigByHostname("example.com")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		configs = append(configs, config)
	}

	config, err = FindConfigByHostname("mdw")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		configs = append(configs, config)
	}

	for _, config := range configs {
		fmt.Println("Hostname", config.Hostname)
		fmt.Println("Config for", config.IP)
		fmt.Println("CPU model:", config.CPU)
		fmt.Println("Memory size:", config.Mem)
		fmt.Println("Storage size:", config.Storage)
	}
	jsonData, err := json.Marshal(configs)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = ioutil.WriteFile("configs.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func atoi(s string) int {
	i := 0
	for _, c := range s {
		i = i*10 + int(c-'0')
	}
	return i
}
