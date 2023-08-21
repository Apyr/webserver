package server

import (
	"log"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
)

func Watch(files []string) bool {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Printf("Watcher error: %s\n", err)
		return false
	}
	defer watcher.Close()

	for _, file := range files {
		if err := watcher.Add(file); err != nil {
			log.Printf("Watcher error: %s\n", err)
			return false
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		return false
	case <-watcher.Events:
		log.Println("Config changed. Reloading...")
		return true
	case err := <-watcher.Errors:
		log.Printf("Watcher error: %s\n", err)
		return false
	}
}
