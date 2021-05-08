package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func httpError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Error: %v", err)
}

func getParsedBody(r *http.Request) (url.Values, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	parsed, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func (s *Server) renamePOST(w http.ResponseWriter, r *http.Request) {
	parsed, err := getParsedBody(r)
	if err != nil {
		httpError(w, err)
		return
	}

	oldName := mux.Vars(r)["image"]
	newName := parsed.Get("new_name")
	newName = fmt.Sprintf("%s%s", newName, filepath.Ext(oldName))

	// TODO: verify file is in the root directory
	// TODO: actually rename the files

	log.Printf("Renaming %q to %q", oldName, newName)

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
