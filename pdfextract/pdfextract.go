package pdfextract

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/cheggaaa/go-poppler"
)

type PDFExtractOptions struct {
	Input       string
	Destination string
	Format      string
}

func PDFExtract(options *PDFExtractOptions) error {
	doc, err := poppler.Open(options.Input)
	if err != nil {
		return err
	}
	defer doc.Close()

	pageCount := doc.GetNPages()

	for i := 0; i < pageCount; i++ {
		images, err := checkPage(doc, i)
		if err != nil {
			return fmt.Errorf("Error checking page %d: %w", i, err)
		}

		for _, img := range images {
			err := img.Save(filepath.Join(options.Destination, options.Format))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func checkPage(doc *poppler.Document, pageID int) ([]*PDFImage, error) {
	log.Printf("Checking page %d", pageID)

	page := doc.GetPage(pageID)
	defer page.Close()

	res := []*PDFImage{}
	log.Printf("Page %d has %d images", pageID, len(page.Images()))
	for _, img := range page.Images() {
		res = append(res, NewImage(pageID, img))
	}

	return res, nil
}
