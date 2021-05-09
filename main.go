package main

import (
	"flag"
	"log"
	"regexp"

	"github.com/keyneston/pdf-extract/pdfextract"
)

func main() {
	var fileName, outputDir, format string
	var writeText bool
	var seperator string
	var startPage int // TODO: parse pages as: 1-5,6-10

	flag.StringVar(&fileName, "f", "", "File to extract images from")
	flag.StringVar(&outputDir, "d", "", "Directory to output images to")
	flag.StringVar(&format, "t", "page_{Page:03d}_id_{ID:03d}.png", "Template to use for generating file names")
	flag.StringVar(&seperator, "s", "\n+", "Seperator to break text on")
	flag.IntVar(&startPage, "p", -1, "Page to start on")
	flag.BoolVar(&writeText, "text", false, "Write a copy of the text to a file")
	flag.Parse()

	var pages []int = nil
	if startPage >= 0 {
		pages = []int{startPage}
	}
	if err := pdfextract.PDFExtract(&pdfextract.PDFExtractOptions{
		Input:       fileName,
		Destination: outputDir,
		Format:      format,
		WriteText:   writeText,
		Pages:       pages,
		Seperator:   regexp.MustCompile(seperator),
	}); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

type PageRanges struct {
	pages []int
}
