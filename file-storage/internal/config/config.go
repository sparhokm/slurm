package config

import (
	"fmt"
	"time"
)

type Config struct {
	Env         string     `mapstructure:"env"`
	HTTPServer  HTTPServer `mapstructure:"http_server"`
	Log         Log        `mapstructure:"log"`
	Minio       Minio      `mapstructure:"minio"`
	Register    Register   `mapstructure:"register"`
	Tracer      Tracer     `mapstructure:"tracer"`
	ServiceName string     `mapstructure:"service_name"`
}

type Log struct {
	MinLevel string `mapstructure:"min_level"`
}

type HTTPServer struct {
	Address     string        `mapstructure:"address"`
	Port        int           `mapstructure:"port"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
	Debug       bool          `mapstructure:"debug"`
}

func (h HTTPServer) ListenAddrAndPort() string {
	return fmt.Sprintf("%s:%d", h.Address, h.Port)
}

type Register struct {
	URL         string        `mapstructure:"url"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type Tracer struct {
	Version  string `mapstructure:"version"`
	Enable   bool   `mapstructure:"enable"`
	Endpoint string `mapstructure:"endpoint"`
}
