package pdfextract

import (
	"fmt"
	"log"

	"github.com/cheggaaa/go-poppler"
	"github.com/keyneston/tabslib"
)

func DoThing(name string) error {
	doc, err := poppler.Open(name)
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

		break
	}

	log.Printf("Matches: %v", tabslib.PrettyString(FindSets(images)))

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
