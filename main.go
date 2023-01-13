package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

var (
	hosts    map[string]*Client
	targets  map[string]*Client
)

type Client struct {
	hostname string
	id       string
	conn     net.Conn
}

func main() {
	fmt.Println("Starting server...")

	// Create maps to hold the known hosts and targets
	hosts = make(map[string]*Client)
	targets = make(map[string]*Client)

	// Register the HTTP handlers
	http.HandleFunc("/add-host", handleAddHost)
	http.HandleFunc("/list-hosts", handleListHosts)
	http.HandleFunc("/pick-target", handlePickTarget)
	http.HandleFunc("/send-message", handleSendMessage)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleListHosts(w http.ResponseWriter, r *http.Request) {
	// Send a list of known hosts
	hostList := "Known hosts:\n"
	for _, host := range hosts {
		hostList += host.hostname + "\n"
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
		Id       string `json:"id"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Add the host to the list of known hosts
	host := &Client{
		hostname: request.Hostname,
		id:       request.Id,
		conn:     nil,
	}
	hosts[request.Id] = host
	fmt.Println("Added host", request.Hostname)
	w.Write([]byte("Host added!"))
}

func handlePickTarget(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Decode the request body as JSON
	var request struct {
		Id string `json:"id"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Add the target to the list of targets
	target, ok := hosts[request.Id]
	if !ok {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		return
	}
	targets[target.id] = target
	fmt.Println("Picked target", target.hostname)
	w.Write([]byte("Target picked!"))
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Decode the request body as JSON
	var request struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Send the message to the target
	target, ok := targets[request.Id]
	if !ok {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		return
	}
	_, err = target.conn.Write([]byte(request.Message))
	if err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}
	fmt.Println("Sent message to", target.hostname)
	w.Write([]byte("Message sent!"))
}
