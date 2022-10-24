package config

import "github.com/spf13/viper"

const DefaultConfigFilename = "config"
const FileExtension = "json"

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Pvtbc    PvtbcConfig    `mapstructure:"pvtbc"`
	Server   ServerConfig   `mapstructure:"server"`
	Bus      BusConfig      `mapstructure:"bus"`
}

func Load(path string, filename string) (*Config, error) {
	config := new(Config)

	if len(filename) == 0 {
		filename = DefaultConfigFilename
	}

	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType(FileExtension)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
