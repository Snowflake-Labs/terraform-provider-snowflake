package sdk

import (
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/snowflakedb/gosnowflake"
)

func DefaultConfig() (*gosnowflake.Config, error) {
	return ProfileConfig("default")
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

	if envConfig := envConfig(); envConfig != nil {
		// envConfig takes precedence
		config = MergeConfig(config, envConfig)
	}

	return config, nil
}

func MergeConfig(baseConfig *gosnowflake.Config, mergeConfig *gosnowflake.Config) *gosnowflake.Config {
	if baseConfig == nil {
		return mergeConfig
	}
	if mergeConfig.Account != "" {
		baseConfig.Account = mergeConfig.Account
	}
	if mergeConfig.User != "" {
		baseConfig.User = mergeConfig.User
	}
	if mergeConfig.Password != "" {
		baseConfig.Password = mergeConfig.Password
	}
	if mergeConfig.Role != "" {
		baseConfig.Role = mergeConfig.Role
	}
	if mergeConfig.Region != "" {
		baseConfig.Region = mergeConfig.Region
	}
	if mergeConfig.Host != "" {
		baseConfig.Host = mergeConfig.Host
	}
	return baseConfig
}

func configFile() (string, error) {
	// has the user overwridden the default config path?
	if configPath, ok := os.LookupEnv("SNOWFLAKE_CONFIG_PATH"); ok {
		return configPath, nil
	}
	dir, err := homeDir()
	if err != nil {
		return "", err
	}
	// default config path is ~/.snowflake/config.
	return filepath.Join(dir, ".snowflake", "config"), nil
}

func envConfig() *gosnowflake.Config {
	config := &gosnowflake.Config{}

	if account, ok := os.LookupEnv("SNOWFLAKE_ACCOUNT"); ok {
		config.Account = account
	}
	if user, ok := os.LookupEnv("SNOWFLAKE_USER"); ok {
		config.User = user
	}
	if password, ok := os.LookupEnv("SNOWFLAKE_PASSWORD"); ok {
		config.Password = password
	}
	if role, ok := os.LookupEnv("SNOWFLAKE_ROLE"); ok {
		config.Role = role
	}
	if region, ok := os.LookupEnv("SNOWFLAKE_REGION"); ok {
		config.Region = region
	}
	if host, ok := os.LookupEnv("SNOWFLAKE_HOST"); ok {
		config.Host = host
	}

	return config
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

func configDir() (string, error) {
	dir, err := homeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, ".snowflake"), nil
}

func homeDir() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try build-in module
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	if user.HomeDir == "" {
		return "", errors.New("blank output")
	}

	return user.HomeDir, nil
}
