package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type AnalyzerConfig struct {
	Tempo struct {
		Host string `mapstructure:"host"`
	} `mapstructure:"tempo"`
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func LoadConfig(cfgFile string) (*AnalyzerConfig, error) {
	var cfg AnalyzerConfig
	if err := loadViperConfig(cfgFile, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadViperConfig(cfgFile string, config *AnalyzerConfig) error {
	// check文件是否存在
	if _, err := os.Stat(cfgFile); err != nil {
		return err
	}

	viper.SetConfigFile(cfgFile)
	//viper.SetConfigName("application")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(config); err != nil {
		logrus.Errorf("[LoadViperConfig] - 解析配置文件失败, err = %v", err)
		return err
	}
	return nil
}
