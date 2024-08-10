package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

func parseArgs() (string, *url.URL) {
	if len(os.Args) != 3 {
		log.Errorf("Usage: %v <bind address/port> <proxying URL>", os.Args[0])
		os.Exit(1)
	}
	webServerAddr := os.Args[1]
	if _, err := strconv.Atoi(webServerAddr); err == nil {
		webServerAddr = fmt.Sprintf(":%v", webServerAddr)
	}

	proxyingURL, err := url.Parse(os.Args[2])
	if err != nil {
		log.Errorf("Error parsing URL: %v", err)
		os.Exit(1)
	}

	log.Infof("Proxying requests to %v", proxyingURL)

	return webServerAddr, proxyingURL
}

func launchWebserver(webServerAddr string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			log.Errorf("Error handling request: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	log.Infof("Listening on %v", webServerAddr)
	log.Fatal(http.ListenAndServe(webServerAddr, nil))
}

var proxyingURL *url.URL

func main() {
	webServerAddr, localProxyingURL := parseArgs()
	proxyingURL = localProxyingURL
	launchWebserver(webServerAddr)
}

func handler(w http.ResponseWriter, r *http.Request) error {
	newURL := *proxyingURL // copy
	newURL.Path = path.Join(newURL.Path, r.URL.Path)
	r.URL = &newURL
	r.RequestURI = ""
	log.Infof("Proxying request to %v", r.URL)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("error make proxied request: %v", err)
	}

	// copy headers
	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	if resp.ContentLength == -1 {
		log.Info("Content-Length is not set, reading response body to set it")
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(respBody)))
		w.WriteHeader(resp.StatusCode)
		_, err = w.Write(respBody)
		if err != nil {
			return fmt.Errorf("error writing response body: %v", err)
		}
	} else {
		log.Infof("Content-Length already set (%v)", resp.ContentLength)
		w.WriteHeader(resp.StatusCode)
		_, err := io.Copy(w, resp.Body)
		if err != nil {
			return fmt.Errorf("error copying response body: %v", err)
		}
	}
	return nil

}
