package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {

	// define origin server URL
	originServerURL, err := url.Parse("http://127.0.0.1:2999")
	if err != nil {
		log.Fatal("invalid origin server URL")
	}

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("request received: %s %s\n", req.Method, req.URL.Path)

		// set req Host, URL and Request URI to forward a request to the origin server
		req.Host = originServerURL.Host
		req.URL.Host = originServerURL.Host
		req.URL.Scheme = originServerURL.Scheme
		req.RequestURI = ""

		// send a request to the origin server
		originServerResponse, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("failed to send a request to the origin server: %s\n", err)
			log.Printf("probably the league client is not running or not in game")
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}

		if originServerResponse.StatusCode != http.StatusOK {
			log.Printf("origin server returned %d\n", originServerResponse.StatusCode)
		}

		// return response to the client
		rw.WriteHeader(originServerResponse.StatusCode)
		_, err = io.Copy(rw, originServerResponse.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}
	})

	log.Println("starting reverse proxy server on port 2998")
	log.Fatal(http.ListenAndServe(":2998", reverseProxy))

}
