package sdk

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/snowflakedb/gosnowflake"
)

func DefaultConfig() *gosnowflake.Config {
	config, err := ProfileConfig("default")
	if err != nil || config == nil {
		log.Printf("[DEBUG] No Snowflake config file found, returning empty config: %v\n", err)
		config = &gosnowflake.Config{}
	}
	return config
}

func ProfileConfig(profile string) (*gosnowflake.Config, error) {
	configs, err := loadConfigFile()
	if err != nil {
		return nil, err
	}

	if profile == "" {
		profile = "default"
	}
	var config *gosnowflake.Config
	if cfg, ok := configs[profile]; ok {
		log.Printf("[DEBUG] loading config for profile: \"%s\"", profile)
		config = cfg
	}

	if config == nil {
		log.Printf("[DEBUG] no config found for profile: \"%s\"", profile)
		return nil, nil
	}

	// us-west-2 is Snowflake's default region, but if you actually specify that it won't trigger the default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if config.Region == "us-west-2" {
		config.Region = ""
	}

	return config, nil
}

func MergeConfig(baseConfig *gosnowflake.Config, mergeConfig *gosnowflake.Config) *gosnowflake.Config {
	if baseConfig == nil {
		return mergeConfig
	}
	if baseConfig.Account == "" {
		baseConfig.Account = mergeConfig.Account
	}
	if baseConfig.User == "" {
		baseConfig.User = mergeConfig.User
	}
	if baseConfig.Password == "" {
		baseConfig.Password = mergeConfig.Password
	}
	if baseConfig.Role == "" {
		baseConfig.Role = mergeConfig.Role
	}
	if baseConfig.Region == "" {
		baseConfig.Region = mergeConfig.Region
	}
	if baseConfig.Host == "" {
		baseConfig.Host = mergeConfig.Host
	}
	return baseConfig
}

func configFile() (string, error) {
	// has the user overwridden the default config path?
	if configPath, ok := os.LookupEnv("SNOWFLAKE_CONFIG_PATH"); ok {
		if configPath != "" {
			return configPath, nil
		}
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// default config path is ~/.snowflake/config.
	return filepath.Join(dir, ".snowflake", "config"), nil
}

func loadConfigFile() (map[string]*gosnowflake.Config, error) {
	path, err := configFile()
	if err != nil {
		return nil, err
	}
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s map[string]*gosnowflake.Config
	err = toml.Unmarshal(dat, &s)
	if err != nil {
		log.Printf("[DEBUG] error unmarshalling config file: %v\n", err)
		return nil, nil
	}
	return s, nil
}
