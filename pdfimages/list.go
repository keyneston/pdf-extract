package pdfimages

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
)

var (
	splitRe = regexp.MustCompile(`\S+`)
)

type List struct {
	Entries []*Entry
}

func GetList(fileName string) (*List, error) {
	buf := &bytes.Buffer{}
	cmd := exec.Command("pdfimages", "-list", fileName)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	log.Printf("Got %d bytes", buf.Len())

	list := &List{}

	var mappings map[string]int
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if line == nil {
			break
		}

		// parse column headers here
		if mappings == nil {
			mappings = parseHeaders(string(line))
		}
		if skipLine(line) {
			continue
		}

		entry, err := NewEntry(mappings, splitRe.Split(string(line), -1))
		if err != nil {
			return nil, err
		}

		list.Entries = append(list.Entries, entry)

	}

	return list, nil
}

func parseHeaders(line string) map[string]int {
	res := map[string]int{}

	splitLine := splitRe.Split(line, -1)
	for i, name := range splitLine {
		res[name] = i
	}

	return res
}

func skipLine(line []byte) bool {
	return bytes.HasPrefix(line, []byte("----"))
}
