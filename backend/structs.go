package main

type ServerConfig struct {
	Servers struct {
		Additional map[string]interface{} `yaml:",inline"`
	} `yaml:"servers"`
	Config struct {
		MaxServers int `yaml:"maxServers"`
	} `yaml:"config"`
}
