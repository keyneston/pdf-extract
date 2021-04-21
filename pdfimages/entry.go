package pdfimages

import (
	"fmt"
	"strconv"
	"strings"

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
			entry.Width, err = strconv.Atoi(cur)
		case "height":
			entry.Height, err = strconv.Atoi(cur)
		case "color":
			entry.Color = cur
		case "comp":
			entry.Comp, err = strconv.Atoi(cur)
		case "bpc":
			entry.BPC, err = strconv.Atoi(cur)
		case "enc":
			entry.ENC = cur
		case "interp":
			entry.Interp, err = parseBool(cur)
		case "object":
			entry.Object, err = strconv.Atoi(cur)
		case "ID":
			entry.ID, err = strconv.Atoi(cur)
		case "x-ppi":
			entry.XPPI, err = strconv.Atoi(cur)
		case "y-ppi":
			entry.YPPI, err = strconv.Atoi(cur)
		case "size":
			entry.Size = cur
		case "ratio":
			entry.Ratio, err = parsePercent(cur)
		}
		if err != nil {
			return nil, errors.Wrap(err, "parsing failed")
		}

	}

	return &Entry{}, nil
}

func parsePercent(input string) (float32, error) {
	output, err := strconv.ParseFloat(strings.Trim("%", input), 32)
	return float32(output / 100), err
}

func parseBool(input string) (bool, error) {
	switch input {
	case "no":
		return false, nil
	case "yes":
		return true, nil
	default:
		return false, fmt.Errorf("unable to parse bool %q", input)

	}
}
