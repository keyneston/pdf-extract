package main

import (
	"context"
	"flag"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/keyneston/pdf-extract/pdfimages"
	"github.com/keyneston/pdf-extract/unit"
)

func main() {
	var fileName, outputDir string
	var skipClean bool
	var count int

	flag.StringVar(&fileName, "f", "", "File to extract images from")
	flag.StringVar(&outputDir, "d", "", "Directory to output images to")
	flag.BoolVar(&skipClean, "skip-clean", false, "Skip cleaning up work in progress")
	flag.IntVar(&count, "n", -1, "Number of files to convert; for testing purposes")
	flag.Parse()

	list, err := pdfimages.GetList(fileName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := pdfimages.Extract(list.Pages, fileName, outputDir); err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Extracted to %s", outputDir)
	if !skipClean {
		defer func() {
			if err := filepath.Walk(outputDir, cleanup); err != nil {
				log.Printf("Error cleaning: %v", err)
			}
		}()
	}

	units, err := unit.NewUnits(outputDir, list.Matches())
	if err != nil {
		log.Fatal(err)
	}

	for i, u := range units {
		names := []string{}

		if len(names) == 0 {
			continue
		}

		if count > 0 && i >= count {
			continue
		}

		combinedName := filepath.Join(outputDir, u.CombinedName())
		if len(names) == 1 {
			log.Printf("Moving: %q", names[0])
			if err := os.Rename(names[0], combinedName); err != nil {
				log.Fatalf("Error: %v", err)
			}
		} else {
			if err := combine(combinedName, names); err != nil {
				log.Fatalf("Error: %v", err)
			}
		}

	}
}

func combine(output string, inputs []string) error {
	log.Printf("convert %v => %s", inputs, output)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	args := []string{}
	args = append(args, inputs...)
	args = append(args, "-alpha", "off", "-compose", "copy-opacity", "-background", "none", "-composite", output)

	log.Printf(`Running "convert %s"`, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, "convert", args...)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func cleanup(path string, info fs.FileInfo, err error) error {
	if !info.Mode().IsRegular() || strings.Contains("-comb", info.Name()) {
		return nil
	}

	if !strings.Contains("-set-", info.Name()) {
		return nil
	}

	return os.Remove(filepath.Join(path, info.Name()))
}
