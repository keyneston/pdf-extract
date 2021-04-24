package unit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/keyneston/pdf-extract/pdfimages"
	"github.com/keyneston/tabslib"
	"github.com/pkg/errors"
)

// Unit represents a unit of work that is in progress of being done
type Unit struct {
	entries []*pdfimages.Entry
	inputs  []string
	id      int
	ext     string
	dir     string
}

func NewUnits(dir string, idSets [][]*pdfimages.Entry) ([]*Unit, error) {
	// fileMap, err := getDirectoryMap(dir)
	// if err != nil {
	// 	return nil, fmt.Errorf("error creating filemap: %w", err)
	// }
	// fmt.Println(tabslib.PrettyString(fileMap))

	pageStarts, err := getPageStartCounts(dir)
	if err != nil {
		return nil, fmt.Errorf("error creating pageStarts: %w", err)
	}
	fmt.Println(tabslib.PrettyString(pageStarts))

	res := make([]*Unit, 0, len(idSets))
	for _, set := range idSets {
		u, err := newUnit(pageStarts, dir, set)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}

	log.Fatalf("testing")

	return res, nil
}

func newUnit(pageStarts map[int]int, dir string, set []*pdfimages.Entry) (*Unit, error) {
	u := &Unit{
		entries: set,
		inputs:  make([]string, len(set)),
		id:      set[0].Object,
		ext:     "jpg", // TODO: discover this rather than assuming
		dir:     dir,
	}

	u.WriteDebugJSON()

	for i, entry := range set {
		if entry.ENC != "jpeg" {
			continue
		}

		log.Printf("entry.Num(%d), entry.Page(%d), pageStart(%d)", entry.Num, entry.Page, pageStarts[entry.Page])
		fileName := fmt.Sprintf("file-%03d-%03d.%s", entry.Page, (entry.Num - pageStarts[entry.Page]), enc2ext(entry.ENC))

		if _, err := os.Stat(filepath.Join(dir, fileName)); err != nil {
			log.Printf("Can't find %q: %v", fileName, err)
			continue
		}

		ext := filepath.Ext(fileName)
		newName := filepath.Join(dir, fmt.Sprintf("%04d-set-%d%s", u.id, i, ext))

		log.Printf("Renaming: %q => %q", fileName, newName)
		if err := os.Rename(filepath.Join(dir, fileName), newName); err != nil {
			return nil, err
		}

		u.inputs[i] = newName
	}

	return u, nil
}

func (u *Unit) WriteDebugJSON() error {
	f, err := os.OpenFile(
		filepath.Join(u.dir, fmt.Sprintf("%04d-debug.json", u.id)),
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error creating debug json: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(u.entries)
}

func (u *Unit) CombinedName() string {
	return fmt.Sprintf("%04d-comb.%s", u.id, u.ext)
}

func getDirectoryMap(dir string) (map[string]string, error) {
	fsInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "getting directory map")
	}

	mappings := map[string]string{}
	for _, fs := range fsInfos {
		name := fs.Name()
		mappings[name] = name
	}

	return mappings, nil
}

func getPageStartCounts(dir string) (map[int]int, error) {
	perPage := []int{}

	fsInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "getting directory map")
	}

	for _, fs := range fsInfos {
		name := fs.Name()
		var page int

		if !strings.HasPrefix(name, "file-") {
			log.Printf("Skipping: %q", name)
			continue
		}

		_, err := fmt.Sscanf(name, "file-%03d", &page)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		if len(perPage) < (page + 1) {
			newPerPage := make([]int, page+1)
			copy(newPerPage, perPage)
			perPage = newPerPage
		}

		perPage[page] += 1
	}

	pageStart := map[int]int{}
	sum := 0
	for page, count := range perPage {
		sum += count
		pageStart[page+1] = sum
	}

	return pageStart, nil
}

func enc2ext(enc string) string {
	switch enc {
	case "jpeg":
		return "jpg"
	default:
		return enc
	}
}
