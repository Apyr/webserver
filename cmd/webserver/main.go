package main

import "github.com/apyr/webserver/app"

func main() {
	app := app.NewApp()
	defer app.Close()
	app.Run()
}
