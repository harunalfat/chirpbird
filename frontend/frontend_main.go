package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func requestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		log.Printf("Body: %s\nAddr: %s\nTarget Host: %s\nUser-Agent: %s", body, r.RemoteAddr, r.Host, r.UserAgent())

		next.ServeHTTP(rw, r)
	})
}

func main() {
	fs := http.FileServer(http.Dir("./dist/client"))
	http.Handle("/", requestLogMiddleware(fs))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Serving on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("Cannot start frontend http server\n%s", err)
	}
}
