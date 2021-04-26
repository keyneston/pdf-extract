package pdfextract

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/cheggaaa/go-poppler"
)

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

func FindSets(inputs []*PDFImage) [][]*PDFImage {
	matches := map[string][]*PDFImage{}

	for _, in := range inputs {
		key := in.matchKey()
		matches[key] = append(matches[key], in)
	}

	res := [][]*PDFImage{}

	for _, v := range matches {
		res = append(res, v)
	}

	return res
}

func (p *PDFImage) Save(loc string) error {
	f, err := os.OpenFile(loc, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return p.Write(f)
}

func (p *PDFImage) Write(w io.Writer) error {
	return jpeg.Encode(w, p.Image, &jpeg.Options{Quality: 100})
}
