package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"server-manager/auth"
	"server-manager/types"
	"strings"
	"sync"
	"time"
)

var stateMutex sync.RWMutex
var state types.State = types.State{
	Servers: make(map[string]types.ServerState),
}

var config types.Config

// this prevents docker compose up/down commands from being run at the same time since the manager is invoked by both the goroutine and the start/stop endpoints
var containerManagerMutex sync.RWMutex

func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	userInfo, err := auth.ValidateUser(r, config)
	if err != nil {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	filteredState := types.State{Servers: make(map[string]types.ServerState)}

	for _, visibleServer := range userInfo.AllowedServers {
		server, exists := state.Servers[visibleServer]
		if exists {
			filteredState.Servers[visibleServer] = server
		}
	}

	data, err := json.Marshal(filteredState)

	w.Header().Set("Content-Type", "application/json")
	io.Writer.Write(w, data)

}

func authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// this does the body reading too
	UserInfo, err := auth.ValidateUser(r, config)
	if err != nil {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	str, err := json.Marshal(UserInfo)

	io.Writer.Write(w, str)
}

func server(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	UserInfo, err := auth.ValidateUser(r, config)
	if err != nil {
		http.Error(w, "Invalid Authentication", http.StatusForbidden)
		return
	}

	var url = strings.Split(r.URL.String(), "/")
	if url[0] == "" {
		url = url[1:]
	}

	if len(url) < 3 {
		http.Error(w, "Provide a server name and command (start, extend)", http.StatusBadRequest)
		io.WriteString(w, "")
		return
	}

	//command (start, stop, extend); server (satisfactory, minecraft)
	var command = url[2]
	var server = url[1]

	if !auth.HasAuth(UserInfo, server, command) {
		http.Error(w, "Missing permission", http.StatusForbidden)
		io.WriteString(w, "")
		return
	}

	switch command {
	case "start":
		err = startServer(server)
	case "extend":
		err = extendServer(server)
	case "stop":
		err = stopServer(server)
	default:
		http.Error(w, "Provide a valid command (start, extend)", http.StatusBadRequest)
		io.WriteString(w, "")
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		io.WriteString(w, "")
		return
	}

	fmt.Printf("%s: command: %s: %s\n", UserInfo.Name, command, server)

	io.WriteString(w, "")
}

func startServer(name string) error {
	now := time.Now().Unix()
	serverConfig, exists := config.Servers[name]

	if !exists {
		return fmt.Errorf("Invalid server name")
	}

	// this cant error since the input was validated already
	serverTTL, _ := DecodeTime(serverConfig.InitialTTL)

	// Get the number of started servers and check if its at or above cap

	stateMutex.Lock()
	serverState, exists := state.Servers[name]

	if !exists {
		serverState = types.ServerState{
			StartedAt:  0,
			Extensions: make([]int64, 0),
			EndsAt:     0,
		}
	}

	serverState.StartedAt = now
	serverState.Extensions = make([]int64, 0)
	serverState.EndsAt = now + serverTTL

	state.Servers[name] = serverState
	stateMutex.Unlock()

	// do the starting of the server
	manageDockerContainers()

	SaveState()
	return nil
}

func extendServer(name string) error {
	serverConfig, exists := config.Servers[name]

	if !exists {
		return fmt.Errorf("Invalid server name")
	}

	maxTimeLeftBeforeExtended, _ := DecodeTime(serverConfig.MaxTimeBeforeExtend)
	timeToExtendBy, _ := DecodeTime(serverConfig.ExtendedTTL)

	stateMutex.Lock()
	serverState, exists := state.Servers[name]
	stateMutex.Unlock()

	if !exists {
		return fmt.Errorf("Server state not declared, try starting the server.")
	}

	if serverState.EndsAt-maxTimeLeftBeforeExtended > time.Now().Unix() {
		return fmt.Errorf("Server has too much time remaining to extend")
	}

	if len(serverState.Extensions) >= serverConfig.MaxExtensions && serverConfig.MaxExtensions != -1 {
		return fmt.Errorf("Server has been extended the maximum number of times for this reboot")
	}

	serverState.EndsAt += timeToExtendBy
	serverState.Extensions = append(serverState.Extensions, timeToExtendBy)

	stateMutex.Lock()
	state.Servers[name] = serverState
	stateMutex.Unlock()

	SaveState()
	return nil
}

func stopServer(name string) error {
	_, exists := config.Servers[name]

	if !exists {
		return fmt.Errorf("Invalid server name")
	}

	stateMutex.Lock()
	serverState, exists := state.Servers[name]
	stateMutex.Unlock()

	if !exists {
		return fmt.Errorf("Server state not declared, try starting the server.")
	}

	if serverState.EndsAt <= time.Now().Unix() {
		return fmt.Errorf("Server is off")
	}

	serverState.EndsAt = 0
	serverState.StartedAt = 0
	serverState.Extensions = make([]int64, 0)

	stateMutex.Lock()
	state.Servers[name] = serverState
	stateMutex.Unlock()

	// do the stopping of the server
	manageDockerContainers()

	SaveState()
	return nil
}

func manageDockerContainersThread() {
	for {
		time.Sleep(1 * time.Second)
		now := time.Now().Unix()
		if now%60 == 0 {
			manageDockerContainers()
		}
	}
}

func manageDockerContainers() {
	containerManagerMutex.Lock()
	stateMutex.Lock()

	for name, serverState := range state.Servers {
		serverConfig, exists := config.Servers[name]
		if !exists {
			continue
		}

		// fmt.Println("Name: ", name)

		cmd := exec.Command("docker", "compose", "ps", "--format", "{{json .}}")
		cmd.Dir = serverConfig.Directory

		out, err := cmd.Output()

		if err != nil {
			fmt.Println("Error running `docker ps --format '{{json .}} in ", serverConfig.Directory, ": ", err)
			continue
		}

		var started bool = len(out) >= 5
		var shouldBeStarted bool = serverState.EndsAt > time.Now().Unix()

		// fmt.Println("Started: ", started, "\nShould be started: ", shouldBeStarted);
		// fmt.Println("Out: ", string(out))

		if started && !shouldBeStarted {
			// fmt.Println("Turning off")
			cmd := exec.Command("docker", "compose", "down")
			cmd.Dir = serverConfig.Directory
			_, err := cmd.Output()

			if err != nil {
				fmt.Println("Error shutting down ", serverConfig.Directory, ": ", err)
				continue
			}

			fmt.Printf("server: execute: stop: %s\n", name)

		}

		if !started && shouldBeStarted {
			// fmt.Println("Turning on")
			cmd := exec.Command("docker", "compose", "up", "-d")
			cmd.Dir = serverConfig.Directory
			_, err := cmd.Output()

			if err != nil {
				fmt.Println("Error starting up ", serverConfig.Directory, ": ", err)
				continue
			}

			fmt.Printf("server: execute: start: %s\n", name)
		}
	}

	stateMutex.Unlock()
	containerManagerMutex.Unlock()
}

func main() {
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

	LoadState(config)

	SaveState()

	go manageDockerContainersThread()

	http.HandleFunc("/status", status)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/server/", server)

	http.ListenAndServe(":3000", nil)
}
