package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Server struct {
	count int
	total int
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

	if err := tmpls.ExecuteTemplate(w, "rename.html.tmpl", map[string]interface{}{
		"count": s.count,
		"total": s.total,
		"page":  page,
		"id":    id,
		"path":  f,
		"text":  filterUnique(cleanText(s.text[page])),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
	}

	s.count += 1
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
