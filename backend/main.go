package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var stateMutex sync.RWMutex
var state State = State{
	Servers: make(map[string]ServerState),
}

var configMutex sync.RWMutex
var config Config

func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

}

func server(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var command = strings.Split(r.URL.String(), "/")
	if command[0] == "" {
		command = command[1:]
	}

	if len(command) < 3 {
		http.Error(w, "Provide a server name and command (start, extend)", http.StatusBadRequest)
		io.WriteString(w, "")
		return
	}

	var err error

	switch command[2] {
	case "start":
		err = startServer(command[1])
	case "extend":
		err = extendServer(command[1])
	default:
		http.Error(w, "Provide a valid command (start, extend)", http.StatusBadRequest)
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error", http.StatusBadRequest)
	}

	io.WriteString(w, "200")
}

func startServer(name string) error {
	now := time.Now().Unix()
	serverConfig, exists := config.Servers[name]

	if !exists {
		return fmt.Errorf("Invalid server name")
	}

	serverTTL, _ := DecodeTime(serverConfig.InitialTTL)

	// Get the number of started servers and check if its at or above cap

	stateMutex.Lock()
	serverState, exists := state.Servers[name]

	if !exists {
		serverState = ServerState{
			StartedAt:  0,
			Extensions: make([]int, 0),
			EndsAt:     0,
		}
	}

	serverState.StartedAt = now
	serverState.EndsAt = now + serverTTL

	state.Servers[name] = serverState
	stateMutex.Unlock()

	SaveState()
	return nil
}

func extendServer(name string) error {

	SaveState()
	return nil
}

func main() {
	LoadState()
	err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config: ", err)
		return
	}

	err = ValidateConfig(config)
	if err != nil {
		fmt.Println("Error validating config file: ", err)
		return
	}

	http.HandleFunc("/status/", status)
	http.HandleFunc("/server/", server)

	http.ListenAndServe(":3000", nil)
}
