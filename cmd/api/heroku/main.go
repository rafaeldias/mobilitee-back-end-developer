package main

import (
	"os"
	"net/http"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("Hello"))
	})

	if err := http.ListenAndServe(post, nil); err != nil {
		panic(err)
	}
}
