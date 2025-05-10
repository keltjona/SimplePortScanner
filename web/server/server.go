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
	Target     string               `json:"target"`
	Results    []scanner.ScanResult `json:"results"`
	ScanTime   string               `json:"scanTime"`
	TotalPorts int                  `json:"totalPorts"`
}

// StartServer starts the web server
func StartServer(port int) error {
	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("web/static")))

	// API endpoint for scanning
	http.HandleFunc("/api/scan", enableCORS(handleScan))

	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting web server on http://localhost%s\n", serverAddr)

	return http.ListenAndServe(serverAddr, nil)
}

// enableCORS is a middleware function to add CORS headers to responses
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

// handleScan processes scan requests
func handleScan(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request
	var req scanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request format: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request format: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Received scan request for target: %s, ports: %s", req.Target, req.Ports)

	// Parse ports
	ports, err := scanner.ParsePorts(req.Ports)
	if err != nil {
		log.Printf("Invalid port specification: %v", err)
		http.Error(w, fmt.Sprintf("Invalid port specification: %v", err), http.StatusBadRequest)
		return
	}

	// Create scanner
	scannerInst := scanner.NewScanner(req.Target)

	// Start scan
	log.Printf("Starting scan of %d ports on %s", len(ports), req.Target)
	startTime := time.Now()
	results := scannerInst.ScanPorts(ports)
	duration := time.Since(startTime)
	log.Printf("Scan completed in %s, found %d results", duration, len(results))

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
