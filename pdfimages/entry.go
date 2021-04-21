package pdfimages

import (
	"strconv"

	"github.com/pkg/errors"
)

type Entry struct {
	Page   int
	Num    int
	Type   string
	Width  int
	Height int
	Color  string
	Comp   int
	BPC    int
	ENC    string
	Interp bool
	Object int
	ID     int
	XPPI   int
	YPPI   int
	Size   string
	Ratio  float32
}

func NewEntry(mapping map[string]int, input []string) (*Entry, error) {
	entry := &Entry{}

	var err error
	for name, index := range mapping {
		cur := input[index]

		switch name {
		case "page":
			entry.Page, err = strconv.Atoi(cur)
		case "num":
			entry.Num, err = strconv.Atoi(cur)
		case "type":
			entry.Type = cur
		case "width":
		case "height":
		case "color":
			entry.Color = cur
		case "comp":
		case "bpc":
		case "enc":
		case "interp":
		case "object":
		case "ID":
		case "x-ppi":
		case "y-ppi":
		case "size":
		case "ratio":
		}
		if err != nil {
			return nil, errors.Wrap(err, "parsing failed")
		}

	}

	return &Entry{}, nil
}
