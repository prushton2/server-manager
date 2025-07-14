package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"server-manager/types"
	"strconv"

	"gopkg.in/yaml.v3"
)

func DecodeTime(formattedTime string) (int64, error) { //decodes "12h" or "5w" to 12 hours or 5 weeks in seconds
	unit := string(formattedTime[len(formattedTime)-1])
	duration, err := strconv.ParseInt(formattedTime[:len(formattedTime)-1], 10, 64)
	var time int64 = 0

	if err != nil {
		return 0, fmt.Errorf("Invalid time declaration %s. This should be in the format <time><duration>, ie 12h or 2w", formattedTime)
	}

	switch unit {
	case "h":
		time = 86400 / 24
	case "d":
		time = 86400
	case "w":
		time = 86400 * 7
	case "m":
		time = 86400 * 30
	default:
		return 0, fmt.Errorf("Invalid unit unit %s in time declaration %s", unit, formattedTime)
	}

	return duration * time, nil
}

// it uses config to create the state if not exists
func LoadState(config types.Config) {
	file, err := os.OpenFile("./state.json", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening state.json. Please make sure the file exists: ", err)
		return
	}

	body, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error loading state")
		return
	}

	var newState types.State = types.State{Servers: make(map[string]types.ServerState)}
	var oldState types.State

	err = json.Unmarshal(body, &oldState)
	if err != nil {
		fmt.Println("Error umarshalling state json, resetting: ", err)
		oldState = types.State{
			Servers: make(map[string]types.ServerState),
		}
	}

	for name := range config.Servers {
		old, exists := oldState.Servers[name]
		if exists {
			newState.Servers[name] = old
		} else {
			newState.Servers[name] = types.ServerState{
				StartedAt:  0,
				Extensions: make([]int64, 0),
				EndsAt:     0,
			}
		}
	}

	state = newState

	fmt.Println("Successfully loaded state")
}

func LoadConfig() error {
	file, err := os.OpenFile("./config.yaml", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, &config)
	if err != nil {
		return err
	}

	return nil
}

func SaveState() {
	file, err := os.OpenFile("./state.json", os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening state.json. Please make sure the file exists: ", err)
		return
	}

	str, err := json.Marshal(state)
	if err != nil {
		fmt.Println("Error marshaling state", err)
		return
	}

	_, err = file.Write(str)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
		return
	}
}

func ValidateConfig(config types.Config) error {
	// validate that time strings are correct
	for name, server := range config.Servers {
		_, err := DecodeTime(server.InitialTTL)
		if err != nil {
			return fmt.Errorf("Invalid InitialTTL for server %s: %v\n", name, err)
		}

		_, err = DecodeTime(server.ExtendedTTL)
		if err != nil {
			return fmt.Errorf("Invalid ExtendedTTL for server %s: %v\n", name, err)
		}

		_, err = DecodeTime(server.MaxTimeBeforeExtend)
		if err != nil {
			return fmt.Errorf("Invalid MaxTimeBeforeExtend for server %s: %v\n", name, err)
		}

		if server.MaxExtensions <= -2 {
			return fmt.Errorf("Warning: %s MaxExtensions is %d, consider setting it to a value between -1 and infinity", name, server.MaxExtensions)
		}
	}

	if config.Config.MaxServers <= -2 {
		return fmt.Errorf("Warning: config.MaxServers is %d, consider setting it to a value between -1 and infinity", config.Config.MaxServers)
	}

	return nil
}
