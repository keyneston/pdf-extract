package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	count int
	total int
	files []string
	text  map[int][]string
}

func (s *Server) renameGET(w http.ResponseWriter, r *http.Request) {
	f := mux.Vars(r)["image"]
	tmplVars := map[string]interface{}{
		"count": s.count,
		"total": s.total,
		"path":  f,
		"text":  []string{},
	}

	if !contains(f, s.files) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"status": "not found", "code": 404}`)
		return
	}

	var page, id int
	_, err := fmt.Sscanf(f, "page_%d_id_%d.png", &page, &id)
	if err == nil {
		tmplVars["page"] = page
		tmplVars["id"] = id
		tmplVars["text"] = filterUnique(cleanText(s.text[page]))
	}

	if err := tmpls.ExecuteTemplate(w, "rename.html.tmpl", tmplVars); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
	}

	s.count += 1
}

func (s *Server) indexGET(w http.ResponseWriter, r *http.Request) {
	log.Printf("In indexGET")
	if err := tmpls.ExecuteTemplate(w, "index.html.tmpl", map[string]interface{}{
		"files": s.files,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
	}
}

func (s *Server) renamePOST(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	s.total = len(s.files)

	return nil
}

var cleanTextRE = regexp.MustCompile(`[^a-zA-Z]+`)

func cleanText(in []string) []string {
	out := make([]string, 0, len(in))
	for _, i := range in {
		scrubbed := cleanTextRE.ReplaceAllString(i, "-")
		scrubbed = strings.ToLower(strings.Trim(scrubbed, "-"))

		if scrubbed == "" {
			continue
		}

		out = append(out, scrubbed)
	}

	return out
}

func filterUnique(in []string) []string {
	tmp := map[string]bool{}

	for _, i := range in {
		tmp[i] = true
	}

	out := make([]string, 0, len(tmp))
	for k := range tmp {
		out = append(out, k)
	}

	return out
}

func contains(needle string, haystack []string) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}

	return false
}
