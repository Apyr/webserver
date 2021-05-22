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

func main() {
	if !fileExists("config.yml") {
		config.SaveDefault()
		log.Println("Default config saved")
	}
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Print(cfg.AsYaml())

	go func() {
		needStart := true
		for {
			var close func()
			if needStart {
				close = startServer(cfg)
			} else {
				close = func() {}
			}

			exit := !cfg.Watch()

			close()
			if exit {
				break
			}

			newCfg, err := config.LoadConfig("config.yml")
			if err != nil {
				needStart = false
				log.Printf("Config loading error: %s\n", err)
			} else {
				log.Print(cfg.AsYaml())
				needStart = true
				cfg = newCfg
			}
		}
	}()

	server.WaitInterrupt()
}
