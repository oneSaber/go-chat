package main

import (
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static", files))

	mux.HandleFunc("/register", registrer)
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/login", login)
	server := &http.Server{
		Addr:           "127.0.0.1:5000",
		Handler:        mux,
		ReadTimeout:    time.Duration(10 * int64(time.Second)),
		WriteTimeout:   time.Duration(600 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
