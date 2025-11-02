package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB  *DBconfig
	Srv *ServicePorts
	App *App
}

type DBconfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
	MaxRetries int    `yaml:"max_retries"`
}

type ServicePorts struct {
	BusinessServicePort string `yaml:"business_service"`
	ConsumerServicePort string `yaml:"consumer_service"`
	AdminServicePort    string `yaml:"admin_service"`
	AuthServicePort     string `yaml:"auth_service"`
	WebServicePort      string `yaml:"web_service"`
}

type App struct {
	JwtSecret string `yaml:"jwt_secret"`
}

func New(path string) (*Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read YAML file: %w", err)
	}

	// Create a Config instance
	cnf := &Config{}

	// Unmarshal the YAML data into the Config struct
	if err := yaml.Unmarshal(data, cnf); err != nil {
		return nil, fmt.Errorf("unable to unmarshal YAML: %w", err)
	}

	// Print out the parsed values for debugging (optional)
	fmt.Printf("Parsed config: %+v\n", cnf)

	return cnf, nil
}
