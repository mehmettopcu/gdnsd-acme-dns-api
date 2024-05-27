package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

// GdnsdCtlClient struct is used to execute gdnsdctl commands.
type GdnsdCtlClient struct {
	ConfigDir        string
	TcpSocket        string
	Debug            bool
	Syslog           bool
	Timeout          int
	OneShot          bool
	IgnoreDeadDaemon bool
}

// GdnsdCtlHandler is the handler used to handle HTTP requests.
type GdnsdCtlHandler struct {
	Client *GdnsdCtlClient
}

var gdnsdctl = "gdnsdctl"

// ServeHTTP handles HTTP requests.
func (h *GdnsdCtlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

// payloadsToArgs converts the sent payloads to an argument list for gdnsdctl.
func (h *GdnsdCtlHandler) payloadsToArgs(payloads map[string]string) []string {
	args := []string{}
	for name, payload := range payloads {
		args = append(args, name, payload)
	}
	return args
}

// ExecuteCommand function executes gdnsdctl commands and returns the output.
func (c *GdnsdCtlClient) ExecuteCommand(args ...string) (string, error) {
	cmdArgs := []string{gdnsdctl}
	if c.ConfigDir != "" {
		cmdArgs = append(cmdArgs, "-c", c.ConfigDir)
	}
	if c.TcpSocket != "" {
		cmdArgs = append(cmdArgs, "-s", c.TcpSocket)
	}
	// Other parameters will be added here
	cmdArgs = append(cmdArgs, args...)

	log.Println(cmdArgs)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (c *GdnsdCtlClient) IsGdnsdCtlInstalled() bool {
	cmd := exec.Command("which", gdnsdctl)
	if err := cmd.Run(); err != nil {
		log.Fatal(gdnsdctl, " not found")

		return false
	}
	return true
}
