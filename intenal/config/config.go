package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type (
	Config struct {
		Env      string     `yaml:"env" env-default:"local"`
		DBPath   string     `yaml:"db_path" env-required:"true"`
		TokenTTL time.Time  `yaml:"token_ttl" env-required:"true"`
		GRPC     GRPCConfig `yaml:"grpc"`
	}
	GRPCConfig struct {
		Port    int       `yaml:"port" `
		Timeout time.Time `yaml:"timeout"`
	}
)

func MustLoad() *Config {
	envPath := fetchConfigEnvPath()
	if len(envPath) == 0 {
		panic("config file is empty")
	}
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		panic("config file does not exists: " + envPath)
	}
	var cnf *Config
	if err := cleanenv.ReadConfig(envPath, &cnf); err != nil {
		panic("failed to read config : " + envPath)
	}
	return cnf
}

// fetchConfigEnvPath for laoding config files from right path: ex: "go run main.go config=../local/config.yml" loading
// files from 'local' folder ../local/config.yaml
func fetchConfigEnvPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if len(res) == 0 {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
