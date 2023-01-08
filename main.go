package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
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

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
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

func handleConnection(conn net.Conn) {
	// Read the hostname from the client
	hostname, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	hostname = strings.TrimSpace(hostname)

	// Add the host to the list of known hosts
	hosts[hostname] = conn
	fmt.Println("Added host", hostname)

	// Create a buffer to hold incoming data
	buf := make([]byte, 1024)

	for {
		// Read incoming data
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			delete(hosts, hostname)
			conn.Close()
			return
		}

		// Check for special commands
		data := strings.TrimSpace(string(buf[:n]))
		if data == "list" {
			// Send a list of known hosts
			hostList := "Known hosts:\n"
			for h := range hosts {
				hostList += h + "\n"
			}
			conn.Write([]byte(hostList))
			} else if strings.HasPrefix(data, "connect ") {
				// Connect to another host
				otherHost := strings.TrimPrefix(data, "connect ")
				otherConn, ok := hosts[otherHost]
				if !ok {
					conn.Write([]byte("Error: Unknown host.\n"))
					continue
				}
	
				fmt.Println("Connecting", hostname, "to", otherHost)
	
				// Copy data between the two connections
				go func() {
					_, err := io.Copy(conn, otherConn)
					if err != nil {
						log.Println(err)
					}
					conn.Close()
					otherConn.Close()
				}()
				go func() {
					_, err := io.Copy(otherConn, conn)
					if err != nil {
						log.Println(err)
					}
					conn.Close()
					otherConn.Close()
				}()
			} else {
				// Send an error message
				conn.Write([]byte("Error: Unknown command.\n"))
			}
		}
	}
	