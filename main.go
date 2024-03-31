package main

import "webserver/app"

func main() {
	app := app.NewApp()
	defer app.Close()
	app.Run()
}
