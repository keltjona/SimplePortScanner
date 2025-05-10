package scanner

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePorts parses port specifications into a slice of port numbers
// Supports formats like:
// - Single port: "80"
// - Port range: "20-25"
// - Comma-separated list: "22,80,443"
// - Combinations: "22,80,100-200"
func ParsePorts(portsSpec string) ([]int, error) {
	var ports []int
	
	// Split by comma
	parts := strings.Split(portsSpec, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		
		// Check if it's a range (contains "-")
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}
			
			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start port in range %s: %v", part, err)
			}
			
			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end port in range %s: %v", part, err)
			}
			
			if start > end {
				return nil, fmt.Errorf("start port cannot be greater than end port in range: %s", part)
			}
			
			for port := start; port <= end; port++ {
				if port > 0 && port < 65536 {
					ports = append(ports, port)
				}
			}
		} else {
			// It's a single port
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", part)
			}
			
			if port <= 0 || port >= 65536 {
				return nil, fmt.Errorf("port number out of range (1-65535): %d", port)
			}
			
			ports = append(ports, port)
		}
	}
	
	return ports, nil
}
