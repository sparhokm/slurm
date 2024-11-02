package config

import (
	"github.com/spf13/viper"
)

type Option func(*options)

type options struct {
	configPath []string
}

func WithDefaults(defaults map[string]any) Option {
	return func(_ *options) {
		for k, v := range defaults {
			viper.SetDefault(k, v)
		}
	}
}

func WithFileDefaults(path string) Option {
	return func(o *options) {
		o.configPath = append(o.configPath, path)
	}
}

func WithFileOverwrite(path string) Option {
	return func(o *options) {
		if path != "" {
			o.configPath = append(o.configPath, path)
		}
	}
}

func Load[T any](cfg *T, opts ...Option) (*T, error) {
	viper.AutomaticEnv()

	options := &options{}
	for i := range opts {
		opts[i](options)
	}

	return loadConfigs(cfg, *options)
}

func loadConfigs[T any](cfg *T, opts options) (*T, error) {
	for i := range opts.configPath {
		if err := loadConfig(&cfg, opts.configPath[i]); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func loadConfig[T any](cfg *T, path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	return nil
}
