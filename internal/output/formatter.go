package output

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/user/portscanner/internal/scanner"
	"strings"
	"time"
)

// PrintHeader prints a formatted header for the scan
func PrintHeader(target string, portCount int) {
	color.Cyan("SimplePortScanner v1.0")
	color.Cyan("----------------------------")
	color.White("Target: %s", target)
	color.White("Scanning %d ports...", portCount)
	color.White("Started at: %s", time.Now().Format(time.RFC1123))
	fmt.Println()
}

// PrintResult prints a single scan result
func PrintResult(result scanner.ScanResult) {
	if result.State == "open" {
		if result.Service != "Unknown" && result.Service != "" {
			// Truncate long service banners
			banner := result.Service
			if len(banner) > 40 {
				banner = banner[:37] + "..."
			}
			banner = strings.ReplaceAll(banner, "\n", " ")
			banner = strings.ReplaceAll(banner, "\r", "")
			
			color.Green("Port %5d: %-6s  %s", result.Port, result.State, banner)
		} else {
			color.Green("Port %5d: %s", result.Port, result.State)
		}
	}
}

// PrintSummary prints a summary of the scan results
func PrintSummary(results []scanner.ScanResult, duration time.Duration) {
	var openCount int
	
	for _, result := range results {
		if result.State == "open" {
			openCount++
		}
	}
	
	fmt.Println()
	color.Cyan("Scan completed in %s", duration)
	color.White("Results: %d ports scanned, %d open ports", len(results), openCount)
	
	if openCount == 0 {
		color.Yellow("No open ports found.")
	}
}
