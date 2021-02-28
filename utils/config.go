package utils

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config is structure defining configuration options parsed into the
// application at runtime.
type Config struct {
	Server struct {
		// Port the server listens on
		GrpcPort int    `yaml:"grpc_api_port", envconfig:"GRPC_API_PORT"`
		RestPort int    `yaml:"rest_api_port", envconfig:"REST_API_PORT"`
		Host     string `yaml:"host", envconfig:"SERVER_HOST"`
	} `yaml:"server"`
	Auth0 struct {
		APIClientID     string `yaml:"api_client_id", envconfig:"AUTH0_CLIENT_ID"`
		APIClientSecret string `yaml:"api_client_secret", envconfig:"AUTH0_CLIENT_SECRET"`
	} `yaml:"auth0"`
	Redis struct {
		Host string `yaml:"redis_host", envconfig:"REDIS_HOST"`
		Port string `yaml:"redis_port", envconfig:"REDIS_PORT"`
	} `yaml:"redis"`
}

// ProcessError is the default error handler.
func ProcessError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

// ReadFile reads in the config.yaml file for user defined variables.
func ReadFile(cfg *Config) {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		ProcessError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		ProcessError(err)
	}
}

// ReadEnv reads in environment variables
func ReadEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		ProcessError(err)
	}
}

// GetFullHost generates the full host name given the hostname and port
func GetFullHost(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
