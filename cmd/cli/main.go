package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	
	"github.com/user/portscanner/internal/output"
	"github.com/user/portscanner/internal/scanner"
)

func main() {
	// Parse command line flags
	target := flag.String("target", "", "Target IP address or domain name")
	portsSpec := flag.String("ports", "1-1000", "Ports to scan (e.g., 80, 22-100, 80,443,8080)")
	threads := flag.Int("threads", 100, "Number of concurrent scanning threads")
	timeout := flag.Duration("timeout", 2*time.Second, "Timeout for connection attempts")
	
	flag.Parse()
	
	// Validate target
	if *target == "" {
		fmt.Println("Error: --target flag is required")
		flag.Usage()
		os.Exit(1)
	}
	
	// Check if target is valid
	if _, err := net.ResolveIPAddr("ip", *target); err != nil {
		fmt.Printf("Error: Invalid target '%s': %v\n", *target, err)
		os.Exit(1)
	}
	
	// Parse ports
	ports, err := scanner.ParsePorts(*portsSpec)
	if err != nil {
		fmt.Printf("Error parsing ports: %v\n", err)
		os.Exit(1)
	}
	
	if len(ports) == 0 {
		fmt.Println("Error: No valid ports to scan")
		os.Exit(1)
	}
	
	// Create and configure scanner
	scannerInst := scanner.NewScanner(*target)
	scannerInst.Timeout = *timeout
	scannerInst.Concurrency = *threads
	
	// Print header
	output.PrintHeader(*target, len(ports))
	
	// Start the scan
	startTime := time.Now()
	
	results := scannerInst.ScanPorts(ports)
	
	// Print results
	for _, result := range results {
		if result.State == "open" {
			output.PrintResult(result)
		}
	}
	
	// Print summary
	duration := time.Since(startTime)
	output.PrintSummary(results, duration)
}
