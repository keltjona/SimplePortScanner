package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	
	"github.com/user/portscanner/internal/scanner"
)

// scanRequest represents a request to scan ports
type scanRequest struct {
	Target string `json:"target"`
	Ports  string `json:"ports"`
}

// scanResponse represents a response for a port scan
type scanResponse struct {
	Target     string            `json:"target"`
	Results    []scanner.ScanResult `json:"results"`
	ScanTime   string            `json:"scanTime"`
	TotalPorts int               `json:"totalPorts"`
}

// StartServer starts the web server
func StartServer(port int) error {
	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("web/static")))
	
	// API endpoint for scanning
	http.HandleFunc("/api/scan", handleScan)
	
	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting web server on http://localhost%s\n", serverAddr)
	
	return http.ListenAndServe(serverAddr, nil)
}

// handleScan processes scan requests
func handleScan(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Decode request
	var req scanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	
	// Parse ports
	ports, err := scanner.ParsePorts(req.Ports)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid port specification: %v", err), http.StatusBadRequest)
		return
	}
	
	// Create scanner
	scannerInst := scanner.NewScanner(req.Target)
	
	// Start scan
	startTime := time.Now()
	results := scannerInst.ScanPorts(ports)
	duration := time.Since(startTime)
	
	// Create response
	resp := scanResponse{
		Target:     req.Target,
		Results:    results,
		ScanTime:   duration.String(),
		TotalPorts: len(ports),
	}
	
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	
	// Encode response
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
