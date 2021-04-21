package main

import (
	"flag"
	"log"

	"github.com/keyneston/pdf-extract/pdfimages"
)

func main() {
	var fileName, dest string

	flag.StringVar(&fileName, "f", "", "File to extract images from")
	flag.StringVar(&dest, "d", "", "Directory to output images to")
	flag.Parse()

	list, err := pdfimages.GetList(fileName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("%v", len(list.Entries))
}
