package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	target, err := url.Parse("http://192.168.2.8:8000")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/users/", http.StripPrefix("/users/", httputil.NewSingleHostReverseProxy(target)))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./Documents"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
