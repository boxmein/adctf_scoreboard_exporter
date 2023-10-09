package httpclient

import (
	"crypto/tls"
	"net/http"
)

var HttpClient http.Client

func init() {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	HttpClient = http.Client{
		Transport: transport,
	}
}
