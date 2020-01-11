package cfg

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Load loads configuration from file.
func Load(file string) (cfg Config, err error) {
	cfg.Server.Addr = ":8080"
	cfg.Server.ReadTimeout = 5 * time.Second
	cfg.Server.WriteTimeout = 5 * time.Second
	cfg.Server.IdleTimeout = 5 * time.Second

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return
	}

	return
}

// Config ...
type Config struct {
	Server struct {
		Name         string
		Addr         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}
	DB struct {
		Driver     string
		DataSource string
	}
}
