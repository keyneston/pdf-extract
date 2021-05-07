package main

import (
	"flag"
	"log"

	"github.com/keyneston/pdf-extract/pdfextract"
)

func main() {
	var fileName, outputDir, format string
	var writeText bool

	flag.StringVar(&fileName, "f", "", "File to extract images from")
	flag.StringVar(&outputDir, "d", "", "Directory to output images to")
	flag.StringVar(&format, "t", "page_{Page:03d}_id_{ID:03d}.png", "Template to use for generating file names")
	flag.BoolVar(&writeText, "text", false, "Write a copy of the text to a file")
	flag.Parse()

	if err := pdfextract.PDFExtract(&pdfextract.PDFExtractOptions{
		Input:       fileName,
		Destination: outputDir,
		Format:      format,
		WriteText:   writeText,
	}); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
