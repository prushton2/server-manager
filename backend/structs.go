package main

type Config struct {
	Servers map[string]ServerConfig `yaml:"servers"`
	Config  struct {
		MaxServers int `yaml:"maxServers"`
	} `yaml:"config"`
	Defaults ServerConfig `yaml:"defaults"`
}

type ServerConfig struct {
	Directory           string `yaml:"directory"`
	InitialTTL          string `yaml:"initialTTL"`
	ExtendedTTL         string `yaml:"extendedTTL"`
	MaxTimeBeforeExtend string `yaml:"maxTimeBeforeExtend"`
	MaxExtensions       int    `yaml:"maxExtensions"`
}

type State struct {
	Servers map[string]ServerState `json:"servers"`
}

type ServerState struct {
	StartedAt  int64 `json:"startedAt"`
	Extensions []int `json:"extensions"`
	EndsAt     int64 `json:"endsAt"`
}
