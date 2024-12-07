package config

import (
	"fmt"
)

type Config struct {
	Env         string     `mapstructure:"env"`
	GRPCServer  GRPCServer `mapstructure:"grpc_server"`
	Log         Log        `mapstructure:"log"`
	Database    Database   `mapstructure:"database"`
	Broker      Broker     `mapstructure:"broker"`
	Tracer      Tracer     `mapstructure:"tracer"`
	ServiceName string     `mapstructure:"service_name"`
}

type Log struct {
	MinLevel string `mapstructure:"min_level"`
}

type Database struct {
	URL string `mapstructure:"url"`
}

type Broker struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
}

type Tracer struct {
	Version  string `mapstructure:"version"`
	Enable   bool   `mapstructure:"enable"`
	Endpoint string `mapstructure:"endpoint"`
}

type GRPCServer struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

func (g GRPCServer) ListenAddrAndPort() string {
	return fmt.Sprintf("%s:%d", g.Address, g.Port)
}
