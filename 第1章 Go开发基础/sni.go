package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tls"))
}

func main() {
	t := log.Logger{}
	var err error
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = make([]tls.Certificate, 3)
	// go http server treats the 0'th key as a default fallback key
	tlsConfig.Certificates[0], err = tls.LoadX509KeyPair("test0.pem", "key.pem")
	if err != nil {
		t.Fatal(err)
	}
	tlsConfig.Certificates[1], err = tls.LoadX509KeyPair("test1.pem", "key.pem")
	if err != nil {
		t.Fatal(err)
	}
	tlsConfig.Certificates[2], err = tls.LoadX509KeyPair("test2.pem", "key.pem")
	if err != nil {
		t.Fatal(err)
	}
	tlsConfig.BuildNameToCertificate()

	http.HandleFunc("/", myHandler)
	server := &http.Server{
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConfig,
	}

	listener, err := tls.Listen("tcp", ":8443", tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
	log.Fatal(server.Serve(listener))
}

/*
// request sent to '127.0.0.1:443'
req, _ := http.NewRequest("GET", "https://127.0.0.1/example", nil)

// virtual host set to 'example.com'
req.Host =  "example.com"

// SNI set to 'example.com'
client := http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName:         req.Host, // here
		},
	},
}

client.Do(req)
*/
