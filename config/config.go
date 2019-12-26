package config

import (
	"fmt"
	"os/user"

	"github.com/spf13/viper"
	"github.com/vitwit/faas-akash/types"
)

// Read the config from a yaml file <defaults to $HOME/.akash/config.yaml>
func Read(paths ...string) (*types.Config, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("ERROR_READING_USER_HOME: %v", err)
	}

	// default path for config
	viper.AddConfigPath(u.HomeDir + "/.akash")
	// add any custom paths for the config
	for _, p := range paths {
		viper.AddConfigPath(p)
	}
	viper.SetConfigName("config")
	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error read faas-akash config: %w", err)
	}

	var cfg types.Config
	if err = viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling faas-akash config: %w", err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error validating faas-akash config: %w", err)
	}

	return &cfg, nil
}
