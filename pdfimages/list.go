package pdfimages

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
)

const (
	Binary = "pdfimages"
)

var (
	splitRe = regexp.MustCompile(`[\t ]+`)
)

type List struct {
	Entries []*Entry
	Pages   int
}

func GetList(fileName string) (*List, error) {
	buf := &bytes.Buffer{}
	cmd := exec.Command(Binary, "-list", fileName)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	log.Printf("Got %d bytes", buf.Len())

	list := &List{}
	if err := list.parse(buf); err != nil {
		return nil, err
	}

	return list, nil
}

func (l *List) parse(r io.Reader) error {
	red := bufio.NewReader(r)

	count := 0
	var mappings map[string]int
	for {
		line, err := red.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if line == nil {
			break
		}
		line = bytes.TrimSpace(line)

		// parse column headers here
		if len(mappings) == 0 {
			mappings = parseHeaders(string(line))
			continue
		}
		if skipLine(line) {
			continue
		}

		entry, err := NewEntry(count, mappings, splitRe.Split(string(line), -1))
		if err != nil {
			return err
		}

		l.Entries = append(l.Entries, entry)

		if entry.Page > l.Pages {
			l.Pages = entry.Page
		}

		count += 1
	}

	return nil
}

func (l *List) Matches() [][]*Entry {
	matches := map[string][]int{}

	for i, e := range l.Entries {
		key := e.matchKey()
		matches[key] = append(matches[key], i)
	}

	res := make([][]*Entry, len(matches))

	i := 0
	for _, v := range matches {
		for _, id := range v {
			res[i] = append(res[i], l.Entries[id])
		}
		i++
	}
	return res
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
