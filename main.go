package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/keyneston/pdf-extract/pdfimages"
	"github.com/pkg/errors"
)

func main() {
	var fileName, dest string
	var skipClean bool

	flag.StringVar(&fileName, "f", "", "File to extract images from")
	flag.StringVar(&dest, "d", "", "Directory to output images to")
	flag.BoolVar(&skipClean, "skip-clean", false, "Skip cleaning up work in progress")
	flag.Parse()

	list, err := pdfimages.GetList(fileName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	matchset := list.Matches()

	tmpDir, err := pdfimages.Extract(fileName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Extracted to %s", tmpDir)
	if !skipClean {
		defer os.RemoveAll(tmpDir)
	}

	fileMap, err := getDirectoryMap(tmpDir)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for i, set := range matchset {
		names := []string{}

		for j, entry := range set {
			fileName := fileMap[fmt.Sprintf("file-%03d", entry.LineNum)]
			if fileName == "" {
				log.Printf("Missing file for %s", entry)
				continue
			}
			ext := filepath.Ext(fileName)
			newName := filepath.Join(tmpDir, fmt.Sprintf("set-%04d-%d%s", i, j, ext))

			if err := os.Rename(filepath.Join(tmpDir, fileName), newName); err != nil {
				log.Fatalf("Error: %v", err)
			}

			names = append(names, newName)
		}

		if len(names) == 0 {
			continue
		}

		// combinedName := filepath.Join(dest, fmt.Sprintf("comb-%03d.jpg", set[0].Object))
		// if len(names) == 1 {
		// 	log.Printf("Moving: %q", names[0])
		// 	if err := os.Rename(names[0], combinedName); err != nil {
		// 		log.Fatalf("Error: %v", err)
		// 	}
		// } else {
		// 	if err := combine(combinedName, names); err != nil {
		// 		log.Fatalf("Error: %v", err)
		// 	}
		// }

		// if i == 10 {
		// 	log.Fatalf("Only doing 10")
		// }
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

func getDirectoryMap(dir string) (map[string]string, error) {
	fsInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "getting directory map")
	}

	mappings := map[string]string{}
	for _, fs := range fsInfos {
		name := fs.Name()
		mappings[strings.TrimRight(name, filepath.Ext(name))] = name
	}

	return mappings, nil
}
