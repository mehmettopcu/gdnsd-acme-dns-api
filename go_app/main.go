package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/mehmettopcu/gdnsd-acme-dns-api/client"
)

// GdnsdCtlHandler is the handler used to handle HTTP requests.
type GdnsdCtlHandler struct {
	Client *client.GdnsdCtlClient
}

// generateRandomToken generates a random token.
func generateRandomToken() string {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		// Handle error
		return ""
	}
	return hex.EncodeToString(token)
}

// ServeHTTP handles HTTP requests.
func (h *GdnsdCtlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Token checking
	authHeader := r.Header.Get("Authorization")
	token := os.Getenv("TOKEN")
	if token == "" {
		token = generateRandomToken()
		os.Setenv("TOKEN", token)
		log.Println("Random TOKEN: " + token)
	}
	expectedAuthHeader := "Bearer " + token
	if authHeader != expectedAuthHeader {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.URL.Path {
	case "/acme-dns-01":
		h.handleAcmeDns01(w, r)
	case "/acme-dns-01-flush":
		h.handleAcmeDns01Flush(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// handleAcmeDns01 handles the /acme-dns-01 endpoint.
func (h *GdnsdCtlHandler) handleAcmeDns01(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payloads map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payloads); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Execute the acme-dns-01 command
	output, err := h.Client.ExecuteCommand(append([]string{"acme-dns-01"}, h.payloadsToArgs(payloads)...)...)
	if err != nil {
		http.Error(w, "Error executing acme-dns-01 command", http.StatusInternalServerError)
		return
	}

	// Respond with the result
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "success", "output": %q}`, output)
}

// handleAcmeDns01Flush handles the /acme-dns-01-flush endpoint.
func (h *GdnsdCtlHandler) handleAcmeDns01Flush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payloads map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payloads); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	// Execute the acme-dns-01 command
	if submit, ok := payloads["submit"]; ok && submit == "true" {
		output, err := h.Client.ExecuteCommand([]string{"acme-dns-01-flush"}...)
		if err != nil {
			http.Error(w, "Error executing acme-dns-01-flush command", http.StatusInternalServerError)
			return
		}

		// Respond with the result
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "success", "output": %q}`, output)
	} else {
		http.Error(w, "\"submit\" : \"true\" not exist on the payloads", http.StatusBadRequest)
	}
}

// payloadsToArgs converts the sent payloads to an argument list for gdnsdctl.
func (h *GdnsdCtlHandler) payloadsToArgs(payloads map[string]string) []string {
	args := []string{}
	for name, payload := range payloads {
		args = append(args, name, payload)
	}
	return args
}

func main() {
	log.Println("Starting gdnsdctl HTTP server...")
	configDir := flag.String("configDir", "/etc/gdnsd/", "Configuration directory for gdnsd")
	tcpSocket := flag.String("tcpSocket", "", "TCP socket address for gdnsd")
	// Minimum and maximum values to restrict the port
	const (
		minPort = 0
		maxPort = 65000
	)

	// Define the port flag and its default value
	port := flag.Int("port", 8080, "HTTP server port")
	flag.Parse()

	// Check the user input for the port number
	if *port < minPort || *port > maxPort {
		fmt.Printf("Port number must be between %d and %d.\n", minPort, maxPort)
		return
	}

	// Start the HTTP server on the specified port

	// Parse the arguments
	flag.Parse()

	// Create GdnsdCtlClient
	client := &client.GdnsdCtlClient{
		ConfigDir: *configDir,
		TcpSocket: *tcpSocket,
	}
	if !client.IsGdnsdCtlInstalled() {
		return
	}
	log.Println("Created gdnsdctl client with config directory:", client.ConfigDir, "and TCP socket:", client.TcpSocket)

	// Start the HTTP server
	log.Println("Handler for / registered")
	handler := &GdnsdCtlHandler{Client: client}
	http.Handle("/", handler)

	fmt.Println("Server started on port ", strconv.Itoa(*port))
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
