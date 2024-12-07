package config

import (
	"fmt"
	"time"
)

type Config struct {
	Env         string     `mapstructure:"env"`
	GRPCServer  GRPCServer `mapstructure:"grpc_server"`
	Prom        Prom       `mapstructure:"prom"`
	Log         Log        `mapstructure:"log"`
	Database    Database   `mapstructure:"database"`
	Broker      Broker     `mapstructure:"broker"`
	Outbox      Outbox     `mapstructure:"outbox"`
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

type GRPCServer struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

type Prom struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

func (g Prom) ListenAddrAndPort() string {
	return fmt.Sprintf("%s:%d", g.Address, g.Port)
}

type Outbox struct {
	Duration  time.Duration `mapstructure:"duration_ms"`
	BatchSize int64         `mapstructure:"batch_size"`
}

func (g GRPCServer) ListenAddrAndPort() string {
	return fmt.Sprintf("%s:%d", g.Address, g.Port)
}

type Tracer struct {
	Version  string `mapstructure:"version"`
	Enable   bool   `mapstructure:"enable"`
	Endpoint string `mapstructure:"endpoint"`
}
