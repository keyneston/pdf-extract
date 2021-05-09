package pdfextract

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/cheggaaa/go-poppler"
	"golang.org/x/net/html/charset"
)

type PDFExtractOptions struct {
	Input       string
	Destination string
	Format      string
	Pages       []int
	Seperator   *regexp.Regexp

	WriteText bool
}

func PDFExtract(options *PDFExtractOptions) error {
	fileHashes := map[string]struct{}{}

	doc, err := poppler.Open(options.Input)
	if err != nil {
		return err
	}
	defer doc.Close()

	pageCount := doc.GetNPages()
	content := map[int][]string{}

	pages := options.Pages
	if pages == nil {
		for i := 0; i < pageCount; i++ {
			pages = append(pages, i)
		}
	}
	for _, i := range pages {
		text, images, err := options.checkPage(doc, i)
		if err != nil {
			return fmt.Errorf("Error checking page %d: %w", i, err)
		}

		content[i] = text
		for _, img := range images {
			filename, err := img.FormatString(filepath.Join(options.Destination, options.Format))
			if err != nil {
				return err
			}

			hash, err := img.Hash()
			if err != nil {
				return err
			}
			if _, ok := fileHashes[hash]; ok {
				continue
			}
			fileHashes[hash] = struct{}{}

			if err := img.Save(filename); err != nil {
				return err
			}
		}
	}

	if options.WriteText {
		f, err := os.OpenFile(filepath.Join(options.Destination, "text.json"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(f)
		enc.SetIndent("", "    ")
		if err := enc.Encode(content); err != nil {
			return err
		}
	}

	return nil
}

func (options PDFExtractOptions) checkPage(doc *poppler.Document, pageID int) ([]string, []*PDFImage, error) {
	log.Printf("Checking page %d", pageID)

	page := doc.GetPage(pageID)
	defer page.Close()

	res := []*PDFImage{}
	log.Printf("Page %d has %d images", pageID, len(page.Images()))
	for _, img := range page.Images() {
		res = append(res, NewImage(pageID, img))
	}

	if !options.WriteText {
		return nil, res, nil
	}

	txt := page.Text()
	e, _, _ := charset.DetermineEncoding([]byte(txt), "")

	var err error
	txt, err = e.NewEncoder().String(txt)
	if err != nil {
		return nil, nil, fmt.Errorf("Error encoding string: %w", err)
	}

	return options.Seperator.Split(txt, -1), res, nil
}
