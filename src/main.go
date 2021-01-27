package main

import (
	"os"
	"path/filepath"
	"webserver/config"
	"webserver/server"
)

func main() {
	cfg := config.DefaultConfig()

	certsDir, _ := filepath.Abs("./certs")
	os.Mkdir(certsDir, os.ModePerm|os.ModeDir)
	hosts := cfg.GetHosts()
	handler := cfg.GetHandler()

	closeHttp := server.StartHttpRedirectServer(80)
	closeHttps := server.StartHttpsServer(certsDir, hosts, 443, handler)

	closeChannel := make(chan struct{})

	doClose := func(close func()) {
		close()
		closeChannel <- struct{}{}
	}

	server.WaitInterrupt()

	go doClose(closeHttp)
	go doClose(closeHttps)

	<-closeChannel
	<-closeChannel
}
