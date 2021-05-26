package main

import (
	"log"
	"os"
	"webserver/config"
	"webserver/server"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func startServer(cfg config.Config) func() {
	close := make([]func(), 0)

	handler := server.BuildHandler(cfg)

	if cfg.HttpsEnabled() {
		os.Mkdir(cfg.CertsDir, os.ModePerm|os.ModeDir)
		c := server.StartHttpsServer(cfg.CertsDir, cfg.GetHosts(), cfg.HttpsPort, handler)
		close = append(close, c)
	}

	if cfg.HttpEnabled() {
		if cfg.HttpsEnabled() && cfg.RedirectToHttps {
			c := server.StartHttpRedirectServer(cfg.HttpPort)
			close = append(close, c)
		} else {
			c := server.StartHttpServer(cfg.HttpPort, handler)
			close = append(close, c)
		}
	}

	log.Println("Server started")

	return func() {
		for _, c := range close {
			c()
		}
	}
}

const configFileName = "config.yml"

type loader struct {
	files []string
}

func (loader loader) watch() bool {
	return server.Watch(loader.files)
}

func (loader *loader) load() (config.Config, bool) {
	for {
		cfg, err := config.LoadConfig(configFileName)
		if err != nil {
			log.Printf("Config loading error: %s\n", err.Error())
			ok := loader.watch()
			if !ok {
				return config.Config{}, ok
			}
			continue
		}
		loader.files = cfg.ConfigFiles
		return cfg, true
	}
}

func main() {
	if !fileExists(configFileName) {
		config.SaveDefault()
		log.Println("Default config saved")
	}

	loader := loader{[]string{configFileName}}

	for {
		cfg, ok := loader.load()
		if !ok {
			log.Println("Interrupted")
			break
		}
		close := startServer(cfg)
		ok = loader.watch()
		close()
		if !ok {
			log.Println("Interrupted")
			break
		}
	}
}
