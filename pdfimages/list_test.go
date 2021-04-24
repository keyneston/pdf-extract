package pdfimages

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestListParse(t *testing.T) {
	input := bytes.NewBufferString(
		`page   num  type   width height color comp bpc  enc interp  object ID x-ppi y-ppi size ratio
		--------------------------------------------------------------------------------------------
		   1     0 image     677  1990  gray    1   8  jpeg   no     11717  0   151   151 6553B 0.5%
		   1     1 image     677  1990  icc     3   8  jpeg   no     11723  0   151   151 47.7K 1.2%
		   1     2 smask     677  1990  gray    1   8  jpeg   no     11723  0   151   151 6553B 0.5%
		   1     3 image    1993   673  icc     3   8  jpeg   no     11727  0   445    51 94.9K 2.4%
		   1     4 image     663  1418  gray    1   8  image  no     11734  0    71   321 1024B 0.1%
		   1     5 image     663  1418  icc     3   8  jpeg   no     11740  0    71   321 25.4K 0.9%
		   1     6 smask     663  1418  gray    1   8  image  no     11740  0    71   321 1024B 0.1%
		   1     7 image    1410   664  icc     3   8  jpeg   no     11744  0   151   151 65.0K 2.4%`,
	)

	expected := []*Entry{
		{Page: 1, Num: 0, Type: "image", Width: 677, Height: 1990, Color: "gray",
			Comp: 1, BPC: 8, ENC: "jpeg", Interp: false, Object: 11717,
			ID: 0, XPPI: 151, YPPI: 151, Size: "6553B", Ratio: 0.5 / 100},
	}

	l := &List{}
	if err := l.parse(input); err != nil {
		t.Fatalf("l.parse(...) = %v; want nil", err)
	}
	if l.Entries != nil {
		l.Entries = l.Entries[0:1] // Only test the first entry for now.
	}
	if diff := deep.Equal(l.Entries, expected); diff != nil {
		t.Errorf("list.parse(...) =\n%v", strings.Join(diff, "\n"))
	}

}
