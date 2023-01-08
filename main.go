package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	//"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	//"strings"
)

var (
	hosts map[string]net.Conn
)

func main() {
	fmt.Println("Starting server...")

	// Create a map to hold the known hosts
	hosts = make(map[string]net.Conn)

	// Register the HTTP handler
	http.HandleFunc("/add-host", handleAddHost)
	http.HandleFunc("/list-hosts", handleListHosts)

		// Send a list of known hosts
		hostList := "Known hosts:\n"
		for h := range hosts {
			hostList += h + "\n"
		}

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleListHosts(w http.ResponseWriter, r *http.Request) {
	// Send a list of known hosts
	hostList := "Known hosts:\n"
	for h := range hosts {
		hostList += h + "\n"
	}
	w.Write([]byte(hostList))
}

func handleAddHost(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Decode the request body as JSON
	var request struct {
		Hostname string `json:"hostname"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Add the host to the list of known hosts
	hosts[request.Hostname] = nil
	fmt.Println("Added host", request.Hostname)
}

