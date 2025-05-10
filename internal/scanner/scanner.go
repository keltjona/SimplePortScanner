package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ScanResult contains information about a scanned port
type ScanResult struct {
	Port    int
	State   string
	Service string
}

// Scanner represents a port scanner
type Scanner struct {
	Target      string
	Timeout     time.Duration
	Concurrency int
}

// NewScanner creates a new scanner instance
func NewScanner(target string) *Scanner {
	return &Scanner{
		Target:      target,
		Timeout:     2 * time.Second,
		Concurrency: 100,
	}
}

// ScanPort scans a single port and returns the result
func (s *Scanner) ScanPort(port int) ScanResult {
	result := ScanResult{Port: port, State: "closed"}
	address := fmt.Sprintf("%s:%d", s.Target, port)
	
	conn, err := net.DialTimeout("tcp", address, s.Timeout)
	if err != nil {
		return result
	}
	defer conn.Close()
	
	result.State = "open"
	result.Service = s.grabBanner(conn)
	
	return result
}

// grabBanner attempts to grab the service banner from an open connection
func (s *Scanner) grabBanner(conn net.Conn) string {
	// Set a deadline for reading
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	
	// Try to read the first 1024 bytes
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "Unknown"
	}
	
	return string(buffer[:n])
}

// ScanPorts scans multiple ports concurrently
func (s *Scanner) ScanPorts(ports []int) []ScanResult {
	var results []ScanResult
	resultChan := make(chan ScanResult)
	var wg sync.WaitGroup
	
	// Use a semaphore to limit concurrency
	sem := make(chan struct{}, s.Concurrency)
	
	// Launch goroutines for each port
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			
			// Acquire semaphore slot
			sem <- struct{}{}
			defer func() { <-sem }() // Release slot
			
			resultChan <- s.ScanPort(p)
		}(port)
	}
	
	// Close the result channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Collect results
	for result := range resultChan {
		results = append(results, result)
	}
	
	return results
}
