package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

}

func server(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Println(strings.Split(r.URL.String(), "/"))

	io.WriteString(w, "200")
}

func startServer(name string) {}

func extendServer(name string) {}

func main() {

	http.HandleFunc("/status/", status)
	http.HandleFunc("/server/", server)

	http.ListenAndServe(":3000", nil)
}
