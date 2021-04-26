package pdfextract

import (
	"fmt"
	"image"
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
		i := img.GetSurface().GetImage()
		log.Printf("Image %d-%d: %#v", pageID, img.Id, i.ColorModel())

		res = append(res, NewImage(pageID, img))

		//if err := SaveImage(pageID, img.Id, i); err != nil {
		//	return nil, err
		//}
	}

	return res, nil
}

type PDFImage struct {
	Page  int
	ID    int
	Image image.Image `json:"-"`

	X1 float64
	X2 float64
	Y1 float64
	Y2 float64
}

func NewImage(page int, img poppler.Image) *PDFImage {
	return &PDFImage{
		Page:  page,
		ID:    img.Id,
		Image: img.GetSurface().GetImage(),
		X1:    img.Area.X1,
		Y1:    img.Area.Y1,
		X2:    img.Area.X2,
		Y2:    img.Area.Y2,
	}
}

func (p *PDFImage) Height() float64 {
	return p.Y2 - p.Y1
}

func (p *PDFImage) Width() float64 {
	return p.X2 - p.X1
}

func (p *PDFImage) matchKey() string {
	return fmt.Sprintf("%d-%2.02f-%2.02f-%2.02f-%2.02f", p.Page, p.X1, p.X2, p.Y1, p.Y2)
}

func (p *PDFImage) String() string {
	return fmt.Sprintf("PDFImage{P%d-%d, %2.02f,%2.02f - %2.02f,%2.02f}", p.Page, p.ID, p.X1, p.Y1, p.X2, p.Y2)
}

func FindSets(inputs []*PDFImage) map[string][]*PDFImage {
	matches := map[string][]*PDFImage{}

	for _, in := range inputs {
		key := in.matchKey()
		matches[key] = append(matches[key], in)
	}

	return matches
}

/*
func SaveImage(page, id int, img image.Image) error {

	f, err := os.OpenFile(fmt.Sprintf("/tmp/tokens/test_%03d_%03d.jpeg", page, id), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
}
*/
