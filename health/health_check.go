package health

import (
	"io"
	"log"
	"net/http"
)

func Runhealthcheck() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) { 
		io.WriteString(rw, "I'm healthy") 
	})
	log.Println("Starting healthcheck on port 6060")
	err := http.ListenAndServe("127.0.0.1:6060", mux)
	// err := http.ListenAndServe(":6060", mux)
	log.Fatal(err)
}
