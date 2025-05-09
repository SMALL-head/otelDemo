package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var ApplicationConfig *AnalyzerConfig

type AnalyzerConfig struct {
	Tempo struct {
		Host string `mapstructure:"host"`
	} `mapstructure:"tempo"`
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	DataSource struct {
		Host string `mapstructure:"host"`
	} `mapstructure:"datasource"`
	Httpclient struct {
		Tempo struct {
			Dev struct {
				Host string `mapstructure:"host"`
			} `mapstructure:"dev"`
		} `mapstructure:"tempo"`
		FlareAdmin struct {
			Dev struct {
				Host string `mapstructure:"host"`
			} `mapstructure:"dev"`
		} `mapstructure:"flare-admin"`
	} `mapstructure:"httpclient"`
}

func AnalyzerConfigDeepCopy(cfg *AnalyzerConfig) *AnalyzerConfig {
	newCfg := &AnalyzerConfig{}
	marshal, err := json.Marshal(cfg)
	if err != nil {
		logrus.Errorf("[AnalyzerConfigDeepCopy] - json.Marshal(cfg) failed, err = %v", err)
		return nil
	}
	if err = json.Unmarshal(marshal, newCfg); err != nil {
		logrus.Errorf("[AnalyzerConfigDeepCopy] - json.Unmarshal(marshal, newCfg) failed, err = %v", err)
		return nil
	}
	return newCfg
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

func init() {
	// 读取配置文件
	cfgFile := "application.yaml"
	cfg, err := LoadConfig(cfgFile)
	if err != nil {
		logrus.Fatalf("[init] - 加载配置文件失败, err = %v", err)
	}
	ApplicationConfig = cfg
	logrus.Infof("[init] - 配置文件加载成功. configFile = %s", cfgFile)
}
