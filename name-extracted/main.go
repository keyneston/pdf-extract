package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var tmpl = template.Must(template.New("rename").Parse(`
<html>
<body>
<h1>Page {{.page}}</h1>
<img src="/images/{{.path}}">
<ul>
{{range $i, $a := .text}}
<li>Title: {{$a}}</li>
{{end}}
</ul>
</body>
</html>
`))

type Server struct {
	files []string
	text  map[int][]string
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	var page, id int
	//q := r.URL.Query()

	f := s.files[50]
	if _, err := fmt.Sscanf(f, "page_%d_id_%d.png", &page, &id); err != nil {
		fmt.Fprintf(w, "error: %v", err)
	}

	tmpl.ExecuteTemplate(w, "rename", map[string]interface{}{
		"page": page,
		"id":   id,
		"path": f,
		"text": s.text[page],
	})
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

	mux := http.NewServeMux()
	mux.Handle(
		"/images/",
		http.StripPrefix("/images/",
			http.FileServer(http.Dir("/tmp/tokens"))))
	mux.HandleFunc("/", s.handler)
	log.Fatal(http.ListenAndServe(net.JoinHostPort("", strconv.Itoa(port)), logRequest(mux)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
