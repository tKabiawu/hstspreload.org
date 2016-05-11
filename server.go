package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/chromium/hstspreload.appspot.com/api"
	"github.com/chromium/hstspreload.appspot.com/database/gcd"
)

func main() {
	staticHandler := http.FileServer(http.Dir("files"))
	http.Handle("/", staticHandler)
	http.Handle("/favicon.ico", staticHandler)
	http.Handle("/static/", staticHandler)

	http.HandleFunc("/robots.txt", http.NotFound)

	db, shutdown, err := gcd.NewLocalBackend()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	defer shutdown()

	api := api.API{
		Backend: db,
	}
	if err := api.TestConnection(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	http.HandleFunc("/preloadable", api.Preloadable)
	http.HandleFunc("/removable", api.Removable)
	http.HandleFunc("/status", api.Status)
	http.HandleFunc("/submit", api.Submit)

	http.HandleFunc("/pending", api.Pending)
	http.HandleFunc("/update", api.Update)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
}
