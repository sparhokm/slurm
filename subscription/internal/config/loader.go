package config

import (
	"errors"
	"flag"
	"fmt"
	"log"

	pkgCfg "github.com/sparhokm/slurm/file-storage/pkg/config"
)

const (
	defaultConfig = "config/default.yaml"
)

func MustLoad() *Config {
	var overwriteConfig string
	flag.StringVar(&overwriteConfig, "config", "", "Add overwrite config")
	flag.Parse()

	cfg, err := pkgCfg.Load(
		new(Config),
		pkgCfg.WithFileDefaults(defaultConfig),
		pkgCfg.WithDefaults(map[string]any{
			"env":           "local",
			"log.min_level": "debug",
		}),
		pkgCfg.WithFileOverwrite(overwriteConfig),
	)

	if err != nil {
		log.Fatal(fmt.Errorf("can't load config: %w", err))
	}

	if err = validateConfig(cfg); err != nil {
		log.Fatal(fmt.Errorf("error validate config: %w", err))
	}

	return cfg
}

func validateConfig(config *Config) error {
	if config.Env == "" {
		return errors.New("env must be set")
	}

	if config.GRPCServer.Port <= 0 {
		return errors.New("grpc server port must be positive")
	}

	if config.GRPCServer.Address == "" {
		return errors.New("grpc server address must be set")
	}

	return nil
}
