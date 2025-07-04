package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string
}

// env-default:"production"
type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}


// It is must, so we are should not return error if anything goes wrong
func MustLoad() *Config {
	var configPath string

	// Read configpath from env var
	configPath = os.Getenv("CONFIG_PATH")

	// Read configpath from command line args
	if configPath == "" {
		flags := flag.String("config", "", "Path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config Path is not set")
		}
	}

	// check if configfile exists or not
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exists %s", configPath)
	}

	// parse config using package
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read configfile %s", err.Error())
	}

	return &cfg
}

