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
	orderedmap "github.com/wk8/go-ordered-map"
)

type Server struct {
	root string

	count int
	total int
	files *orderedmap.OrderedMap
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

	if next := s.files.GetPair(f).Next(); next != nil {
		tmplVars["next"] = next.Key.(string)
	}
	if prev := s.files.GetPair(f).Prev(); prev != nil {
		tmplVars["prev"] = prev.Key.(string)
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
	newName = fmt.Sprintf("%s%s", cleanFileName(newName), filepath.Ext(oldName))

	if !contains(oldName, s.files) {
		httpError(w, fmt.Errorf("Can't find file %q", oldName))
		return
	}

	if err := safeFile(s.root, newName); err != nil {
		httpError(w, err)
		return
	}

	log.Printf("Renaming %q to %q", oldName, newName)
	if err := os.Rename(filepath.Join(s.root, oldName), filepath.Join(s.root, newName)); err != nil {
		httpError(w, err)
		return
	}

	pair := s.files.GetPair(oldName)

	redirect := "/"
	if next := pair.Next(); next != nil {
		redirect = filepath.Join("/rename", next.Key.(string))
	}

	s.files.Delete(oldName)

	// Redirect to next file
	// TODO: add some type of ordering to files
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (s *Server) parseDirectory() error {
	if s.files == nil {
		s.files = orderedmap.New()
	}

	entries, err := os.ReadDir(s.root)
	if err != nil {
		return err
	}

	for _, ent := range entries {
		s.files.Set(ent.Name(), true)

		if ent.Name() == "text.json" {
			f, err := os.Open(filepath.Join(s.root, ent.Name()))
			if err != nil {
				return err
			}
			defer f.Close()
			if err := json.NewDecoder(f).Decode(&s.text); err != nil {
				return err
			}
		}
	}
	s.total = s.files.Len()

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

func contains(needle string, haystack *orderedmap.OrderedMap) bool {
	_, ok := haystack.Get(needle)
	return ok
}

func cleanFileName(name string) string {
	return strings.Trim(cleanTextRE.ReplaceAllString(name, "-"), "-")
}

func safeFile(root, name string) error {
	path := filepath.Join(root, name)
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(root, abs)
	switch {
	case err != nil:
		return err
	case strings.HasPrefix(rel, "."):
		return fmt.Errorf("%q does not seem to be a subpath of %q", path, root)
	}

	return nil
}
