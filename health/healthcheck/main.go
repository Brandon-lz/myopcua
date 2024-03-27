package main

// go build -ldflags=-w -o healthcheck main.go
import (
	"net/http"
	"os"
)

func main() {
	_, err := http.Get("http://127.0.0.1:8080/health")
	if err != nil {
		os.Exit(1)
	}
}