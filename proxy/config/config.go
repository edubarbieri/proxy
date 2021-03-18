package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type RouteConfig struct {
	Pattern  string   `json:"pattern"`
	Backends []string `json:"backends"`
}

type LimitConfig struct {
	ID          string `json:"id"`
	RequestMin  int    `json:"requestMin"`
	TargetPath  string `json:"targetPath"`
	SourceIp    bool   `json:"sourceIp"`
	HeaderValue string `json:"headerValue"`
}

type Config struct {
	Routes []RouteConfig `json:"routes"`
	Limits []LimitConfig `json:"limits"`
}

type ConfigObserver interface {
	UpdateConfig(Config)
}

type ConfigManager struct {
	currentConfig      Config
	configObserverList []ConfigObserver
}

func (c *ConfigManager) AddUpdateObserver(observer ConfigObserver) {
	c.configObserverList = append(c.configObserverList, observer)
}

func (c *ConfigManager) Init() {
	c.currentConfig = c.readConfigJson()
	c.notifyUpdateConfig()
}

func (c *ConfigManager) UpdateConfig(config Config) error {
	c.currentConfig = config
	c.notifyUpdateConfig()

	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		log.Printf("error parsing config to json %v\n", err)
		return err
	}
	err = ioutil.WriteFile(os.Getenv("CONFIG_PATH"), file, 0644)
	if err != nil {
		log.Printf("error writing current json config %v\n", err)
		return err
	}
	return nil
}

func (c *ConfigManager) GetCurrentConfig() Config {
	return c.currentConfig
}

func (c *ConfigManager) readConfigJson() Config {
	log.Printf("reading %s config file", os.Getenv("CONFIG_PATH"))
	jsonFile, err := os.Open(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Printf("error open default config file %v", err)
		return Config{}
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Printf("error parsing json config %v", err)
	}
	return config
}

func (c *ConfigManager) notifyUpdateConfig() {
	for _, o := range c.configObserverList {
		o.UpdateConfig(c.currentConfig)
	}
}
