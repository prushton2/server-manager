package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"server-manager/types"
)

// Given a config and password, get the associated user
func GetAuth(password string, config types.Config) (types.UserInfo, error) {
	for name, info := range config.Users {
		if info.Password == password {
			user := types.UserInfo{
				Name:           name,
				CanStart:       info.CanStart,
				CanExtend:      info.CanExtend,
				CanStop:        info.CanStop,
				AllowedServers: info.AllowedServers,
			}
			return user, nil
		}
	}

	return types.UserInfo{}, fmt.Errorf("Password not found")
}

// Read password from body and get the associated user
func ValidateUser(r *http.Request, config types.Config) (types.UserInfo, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return types.UserInfo{}, fmt.Errorf("Error reading request body (is this a post request?)")
	}

	// parse the body
	var parsedBody types.PasswordRequest
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return types.UserInfo{}, fmt.Errorf("Body is not valid JSON")
	}

	return GetAuth(parsedBody.Password, config)
}

// Does this user have the correct auth to do this action on this server
func HasAuth(user types.UserInfo, server string, action string) bool {

	// is the server the user wants to modify in their allowed servers?
	for _, allowed := range user.AllowedServers {
		if allowed == server {
			// can they perform the action?
			switch action {
			case "start":
				return user.CanStart
			case "extend":
				return user.CanExtend
			case "stop":
				return user.CanStop
			default:
				return false
			}
		}
	}

	return false
}
