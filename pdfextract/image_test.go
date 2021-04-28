package pdfextract

import "testing"

func TestEvaluateTemplate(t *testing.T) {
	p := &PDFImage{
		X1:        50,
		X2:        63,
		Y1:        1,
		Y2:        55,
		Height:    60,
		Width:     32,
		Page:      5,
		ID:        6,
		Extension: "jpg",
	}

	res, err := p.evaluateTemplate(`/tmp/tokens/page_{Page:02d}_id_{ID:02d}.{Extension}`)
	if err != nil {
		t.Fatalf("p.evaluateTemplate(...) = _, %q; want _, nil", err.Error())
	}

	expected := "/tmp/tokens/page_05_id_06.jpg"
	if res != expected {
		t.Errorf("p.evaluateTemplate(...) = %q; want %q", res, expected)
	}
}
