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

func main() {
	if !fileExists("config.yml") {
		config.SaveDefault()
		log.Println("Default config saved")
	}
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("%v", err)
	}
	//log.Print(cfg.AsYaml())

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

	server.WaitInterrupt()

	for _, c := range close {
		c()
	}
}
