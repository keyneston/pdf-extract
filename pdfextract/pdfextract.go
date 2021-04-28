package pdfextract

import (
	"fmt"
	"log"

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
	images := []*PDFImage{}

	for i := 0; i < pageCount; i++ {
		newImages, err := checkPage(doc, i)
		if err != nil {
			return fmt.Errorf("Error checking page %d: %w", i, err)
		}
		images = append(images, newImages...)

	}

	for i, set := range FindSets(images) {
		for _, img := range set {
			err := img.Save(fmt.Sprintf(`/tmp/tokens/page_{Page:02d}_set_%03d_id_{ID:03d}.png`, i))
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
