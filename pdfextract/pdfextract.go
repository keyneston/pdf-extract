package pdfextract

import (
	"fmt"
	"log"

	"github.com/cheggaaa/go-poppler"
)

func DoThing(name string) error {
	doc, err := poppler.Open(name)
	if err != nil {
		return err
	}
	defer doc.Close()

	pageCount := doc.GetNPages()
	for i := 0; i <= pageCount; i++ {
		if err := checkPage(doc, i); err != nil {
			return fmt.Errorf("Error checking page %d: %w", i, err)
		}
	}

	return nil
}

func checkPage(doc *poppler.Document, pageID int) error {
	log.Printf("Checking page %d", pageID)

	page := doc.GetPage(pageID)
	defer page.Close()

	log.Printf("Page %d has %d images", pageID, len(page.Images()))
	for _, img := range page.Images() {
		fmt.Printf("ID %d\tHeight: %.02f\tWidth: %.02f\t\n", img.Id, img.Area.X2-img.Area.X1, img.Area.Y2-img.Area.Y1)
	}

	return nil
}
