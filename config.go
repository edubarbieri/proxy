package main

import (
	"encoding/json"
	"httpproxy/middleware"
	"httpproxy/route"
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

func ReadConfigJson() Config {
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

func CreateRoutes(routesConfig []RouteConfig) map[string]*route.Route {
	routesMap := map[string]*route.Route{}
	for _, config := range routesConfig {
		route, error := CreateRoute(config)
		if error != nil {
			log.Printf("error creating route %s - %v", config.Pattern, error)
			continue
		}
		routesMap[config.Pattern] = route
	}

	return routesMap
}

func CreateRoute(r RouteConfig) (*route.Route, error) {
	return route.NewRoute(r.Pattern, r.Backends, []middleware.Middleware{})
}

func CreateRateLimiteRules(limitsConfig []LimitConfig) []middleware.RateLimiteRule {
	rules := []middleware.RateLimiteRule{}
	for _, config := range limitsConfig {
		rule := CreateRateLimiteRule(config)
		rules = append(rules, rule)
	}
	return rules
}
func CreateRateLimiteRule(limitsConfig LimitConfig) middleware.RateLimiteRule {
	return middleware.RateLimiteRule{
		ID:          limitsConfig.ID,
		Limit:       uint64(limitsConfig.RequestMin),
		TargetPath:  limitsConfig.TargetPath,
		SourceIP:    limitsConfig.SourceIp,
		HeaderValue: limitsConfig.HeaderValue,
	}
}
