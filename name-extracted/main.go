package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Server struct {
	files []string
	text  map[int][]string
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (s *Server) parseDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, ent := range entries {
		s.files = append(s.files, ent.Name())

		if ent.Name() == "text.json" {
			f, err := os.Open(filepath.Join(dir, ent.Name()))
			if err != nil {
				return err
			}
			defer f.Close()
			if err := json.NewDecoder(f).Decode(&s.text); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	var port int

	flag.IntVar(&port, "port", 8888, "Port to listen to")
	flag.Parse()

	s := Server{}
	s.parseDirectory("/tmp/tokens")
	log.Fatalf("%#v", s)

	http.HandleFunc("/", s.handler)
	log.Fatal(http.ListenAndServe(net.JoinHostPort("", strconv.Itoa(port)), nil))
}
