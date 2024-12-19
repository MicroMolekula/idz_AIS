package main

import (
	"idz_ais/internal/gui"
	"log"
	"os"

	"gioui.org/app"
)

type PC struct {
	HostName string
	UserName string
	Password string
}

func main() {
	go func() {
		w := new(app.Window)
		if err := gui.Loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
