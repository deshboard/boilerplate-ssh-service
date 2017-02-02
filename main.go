package main

import "log"

func main() {
	app, err := NewApp()
	if err != nil {
		log.Panic(err)
	}
	defer app.Shutdown()

	err = app.Listen()
	if err != nil {
		log.Panic(err)
	}
}
