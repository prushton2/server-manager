package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/SherClockHolmes/webpush-go"
)

var stateMutex sync.RWMutex
var state State = State{
	Servers: make(map[string]ServerState),
}

var config Config

// this prevents docker compose up/down commands from being run at the same time
var containerManagerMutex sync.RWMutex

func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	file, err := os.Open("state.json")
	if err != nil {
		http.Error(w, "Could not open state.json", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, file)

}

func notify(w http.ResponseWriter, r *http.Request) {
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
		io.WriteString(w, "")
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		io.WriteString(w, "")
		return
	}

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
		serverState = ServerState{
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

		}
	}

	stateMutex.Unlock()
	containerManagerMutex.Unlock()
}

func main() {
	vapidPublicKey := os.Getenv("PUBLIC_VAPID_KEY")
	vapidPrivateKey := os.Getenv("PRIVATE_VAPID_KEY")
	email := os.Getenv("VAPID_EMAIL")

	s := webpush.Subscription{}
	json.Unmarshal([]byte("<YOUR_SUBSCRIPTION>"), &s)

	// Send Notification

	resp, err := webpush.SendNotification([]byte("Test"), &s, &webpush.Options{
		Subscriber:      email,
		VAPIDPublicKey:  vapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	err = LoadConfig()
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
	http.HandleFunc("/server/", server)

	http.ListenAndServe(":3000", nil)
}
