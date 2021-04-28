package pdfextract

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/cheggaaa/go-poppler"
	"github.com/slongfield/pyfmt"
	"github.com/ungerik/go-cairo"
)

type PDFImage struct {
	Page int
	ID   int

	Content   string
	Width     int
	Height    int
	Extension string

	Surface *cairo.Surface `json:"-"`

	X1 float64
	X2 float64
	Y1 float64
	Y2 float64
}

func NewImage(page int, img poppler.Image) *PDFImage {
	p := &PDFImage{
		Page:    page,
		ID:      img.Id,
		Surface: img.GetSurface(),
		X1:      img.Area.X1,
		Y1:      img.Area.Y1,
		X2:      img.Area.X2,
		Y2:      img.Area.Y2,
	}

	p.setHeight()
	p.setWidth()
	p.setContent()

	return p
}

func (p *PDFImage) GetImage() image.Image {
	return p.Surface.GetImage()
}

func (p *PDFImage) setHeight() {
	p.Height = p.Surface.GetHeight()
}

func (p *PDFImage) setWidth() {
	p.Width = p.Surface.GetWidth()
}

func (p *PDFImage) setContent() {
	switch p.Surface.GetContent() {
	case cairo.CONTENT_COLOR:
		p.Content = "color"
	case cairo.CONTENT_ALPHA:
		p.Content = "alpha"
	case cairo.CONTENT_COLOR_ALPHA:
		p.Content = "color_alpha"
	default:
		p.Content = ""
	}
}

func (p *PDFImage) String() string {
	return fmt.Sprintf("PDFImage{P%d-%d, %2.02f,%2.02f - %2.02f,%2.02f}", p.Page, p.ID, p.X1, p.Y1, p.X2, p.Y2)
}

func (p *PDFImage) Save(loc string) error {
	loc, err := p.evaluateTemplate(loc)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(loc, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return p.WritePNG(f)
}

func (p *PDFImage) WriteJPEG(w io.Writer) error {
	return jpeg.Encode(w, p.GetImage(), &jpeg.Options{Quality: 100})
}

func (p *PDFImage) WritePNG(w io.Writer) error {
	return png.Encode(w, p.GetImage())
}

func (p *PDFImage) evaluateTemplate(input string) (string, error) {
	return pyfmt.Fmt(input, p)
}
