package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
)

//go:embed static templates
var content embed.FS
var funcMap = template.FuncMap{}
var tmpls = template.Must(template.New("").Funcs(funcMap).ParseFS(content, "templates/*.tmpl"))

func main() {
	var port int

	flag.IntVar(&port, "port", 8888, "Port to listen to")
	flag.Parse()

	s := Server{}
	s.parseDirectory("/tmp/tokens")

	mux := http.NewServeMux()
	mux.Handle(
		"/images/",
		http.StripPrefix("/images/",
			http.FileServer(http.Dir("/tmp/tokens"))))
	mux.Handle("/static/", http.FileServer(http.FS(content)))
	mux.HandleFunc("/", s.handler)

	address := net.JoinHostPort("localhost", strconv.Itoa(port))
	log.Printf("Starting server %s%s", "http://", address)
	log.Fatal(http.ListenAndServe(address, logRequest(mux)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
