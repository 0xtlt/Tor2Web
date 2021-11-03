package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Server struct{}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	transport := &http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			return url.Parse("socks5://localhost:9001")
		},
	}
	client := &http.Client{
		Transport: transport,
	}

	proxyEnv := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		proxyEnv[pair[0]] = pair[1]
	}

	r.URL.Scheme = os.Getenv("SCHEME")
	r.URL.Host = os.Getenv("HOST")

	res, err := client.Do(&http.Request{
		Proto:  r.URL.Scheme,
		Method: r.Method,
		Host:   r.URL.Host,
		URL:    r.URL,
		Header: r.Header,
		Body:   r.Body,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	for name, values := range res.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}

func main() {
	http.ListenAndServe(":8080", Server{})
}
