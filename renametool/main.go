package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	orderedmap "github.com/wk8/go-ordered-map"
)

//go:embed static templates
var content embed.FS
var funcMap = template.FuncMap{
	"fmt":      fmt.Sprintf,
	"joinpath": filepath.Join,
	"sansext": func(in string) string {
		return strings.TrimSuffix(in, filepath.Ext(in))
	},
	"iter": func(in *orderedmap.OrderedMap) []string {
		ret := make([]string, 0, in.Len())
		for p := in.Oldest(); p != nil; p = p.Next() {
			ret = append(ret, p.Key.(string))
		}
		return ret
	},
}
var tmpls = template.Must(template.New("").Funcs(funcMap).ParseFS(content, "templates/*.tmpl"))

func main() {
	var port int
	var root string

	flag.IntVar(&port, "port", 8888, "Port to listen to")
	flag.StringVar(&root, "root", "/tmp/tokens", "Root directory to work from")
	flag.Parse()

	s := Server{root: root}
	s.parseDirectory()

	r := mux.NewRouter()
	r.PathPrefix("/images/").Handler(
		http.StripPrefix("/images/", http.FileServer(http.Dir(root)))).Methods("GET")
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(content)))
	r.HandleFunc("/rename/{image}", s.renameGET).Methods("GET")
	r.HandleFunc("/rename/{image}", s.renamePOST).Methods("POST")
	r.HandleFunc("/", s.indexGET).Methods("GET")
	r.Use()

	address := net.JoinHostPort("localhost", strconv.Itoa(port))
	log.Printf("Starting server %s%s", "http://", address)
	log.Fatal(http.ListenAndServe(address, handlers.LoggingHandler(os.Stderr, r)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
