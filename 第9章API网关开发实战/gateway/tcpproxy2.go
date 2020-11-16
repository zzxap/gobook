package main

import (
	"log"

	"github.com/inetaf/tcpproxy"
)

func main() {

	var p tcpproxy.Proxy
	p.AddHTTPHostRoute(":80", "foo.com", tcpproxy.To("10.0.0.1:8081"))
	p.AddHTTPHostRoute(":80", "bar.com", tcpproxy.To("10.0.0.2:8082"))
	p.AddRoute(":80", tcpproxy.To("10.0.0.1:8081")) // fallback
	p.AddSNIRoute(":443", "foo.com", tcpproxy.To("10.0.0.1:4431"))
	p.AddSNIRoute(":443", "bar.com", tcpproxy.To("10.0.0.2:4432"))
	p.AddRoute(":443", tcpproxy.To("10.0.0.1:4431")) // fallback
	log.Fatal(p.Run())
}
