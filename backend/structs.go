package main

type Config struct {
	Servers map[string]ServerConfig `yaml:"servers"`
	Config  struct {
		MaxServers int `yaml:"maxServers"`
	} `yaml:"config"`
	Users map[string]UserConfig `yaml:"users"`
}

type ServerConfig struct {
	Directory           string `yaml:"directory"`
	InitialTTL          string `yaml:"initialTTL"`
	ExtendedTTL         string `yaml:"extendedTTL"`
	MaxTimeBeforeExtend string `yaml:"maxTimeBeforeExtend"`
	MaxExtensions       int    `yaml:"maxExtensions"`
}

type UserConfig struct {
	Password       string   `yaml:"password"`
	CanStart       bool     `yaml:"canStart"`
	CanExtend      bool     `yaml:"canExtend"`
	CanStop        bool     `yaml:"canStop"`
	AllowedServers []string `yaml:"allowedServers"`
}

type State struct {
	Servers map[string]ServerState `json:"servers"`
}

type ServerState struct {
	StartedAt  int64   `json:"startedAt"`
	Extensions []int64 `json:"extensions"`
	EndsAt     int64   `json:"endsAt"`
}

type PasswordRequest struct {
	Password string `json:"password"`
}

type UserInfo struct {
	Name           string   `json:"name"`
	CanStart       bool     `json:"canStart"`
	CanExtend      bool     `json:"canExtend"`
	CanStop        bool     `json:"canStop"`
	AllowedServers []string `json:"allowedServers"`
}
