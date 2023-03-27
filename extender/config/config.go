package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"violin-home.cn/violin-extender/common"
)

type Config struct {
	viper                 *viper.Viper
	SC                    *ServerConfig
	LC                    *common.LogConfig
	PodSelectorLabelKey   string
	PodSelectorLabelValue string
	NodeSelectorLabel     string
	ZoneSelectorLabel     string
}

type ServerConfig struct {
	Name string
	Addr string
}

func InitConfig() *Config {
	v := viper.New()
	config := &Config{viper: v}
	applicationPath, err := os.Getwd()
	configPath := applicationPath + string(os.PathSeparator) + "config"
	_, err = os.Stat("/app/config/application.yaml")
	if err == nil {
		configPath = "/app/config/"
	}
	log.Println("load config file : [" + configPath + string(os.PathSeparator) + "application.yaml")
	config.viper.SetConfigName("application")
	config.viper.SetConfigType("yaml")
	config.viper.AddConfigPath(configPath)
	err = v.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	config.ReadServerConfig()
	config.ReadLogsConfig()
	common.InitConfig(config.LC)
	config.ReadPodSelectorLabel()

	return config
}

func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{
		Name: c.viper.GetString("server.name"),
		Addr: c.viper.GetString("server.addr"),
	}
	c.SC = sc
}

func (c *Config) ReadLogsConfig() {
	lc := &common.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("zap.maxSize"),
		MaxAge:        c.viper.GetInt("zap.maxAge"),
		MaxBackups:    c.viper.GetInt("zap.maxBackups"),
	}
	c.LC = lc
}

func (c *Config) ReadPodSelectorLabel() {
	c.PodSelectorLabelKey = os.Getenv("scheduler.pod.selectorLabel.key")
	c.PodSelectorLabelValue = os.Getenv("scheduler.pod.selectorLabel.value")
	c.NodeSelectorLabel = os.Getenv("scheduler.node.selectorLabel")
	c.ZoneSelectorLabel = os.Getenv("scheduler.zone.selectorLabel")
}
