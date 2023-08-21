package internal

import (
	"fmt"
	"path/filepath"
)

func (cfg Config) loadServices(dir string) ([]ServiceEx, error) {
	result := []ServiceEx{}
	for _, item := range cfg.Services {
		key, ok := item.Key.(string)
		if !ok {
			return nil, fmt.Errorf("Service key must be string")
		}
		value, ok := item.Value.(string)
		if !ok {
			return nil, fmt.Errorf("Service value of key \"%s\" must be string", key)
		}
		var err error
		value, err = absPath(value, dir)
		if err != nil {
			return nil, errServiceName(err, key)
		}

		service, err := loadService(key, value)
		if err != nil {
			return nil, errServiceName(err, key)
		}
		result = append(result, service)
	}
	return result, nil
}

func loadService(name, fileName string) (ServiceEx, error) {
	dir := filepath.Dir(fileName)
	service := Service{}
	if err := loadYamlFile(fileName, &service); err != nil {
		return ServiceEx{}, err
	}
	newEndpoints := []Endpoint{}
	for index, endpoint := range service.Endpoints {
		if err := endpoint.check(dir, index); err != nil {
			return ServiceEx{}, err
		}
		newEndpoints = append(newEndpoints, endpoint.splitByHosts()...)
	}
	service.Endpoints = newEndpoints
	return ServiceEx{service, name, fileName}, nil
}

func Load(fileName string) (ConfigEx, error) {
	cfg := Config{}
	if err := loadYamlFile(fileName, &cfg); err != nil {
		return ConfigEx{}, err
	}
	dir := filepath.Dir(fileName)
	if err := cfg.check(dir); err != nil {
		return ConfigEx{}, err
	}
	services, err := cfg.loadServices(dir)
	if err != nil {
		return ConfigEx{}, err
	}
	return ConfigEx{cfg, services}, nil
}
