package pdfimages

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func Extract(totalPages int, fileName, outputdir string) error {
	for i := 1; i <= totalPages; i++ {
		log.Printf("Extracting Page: %d", i)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		cmd := exec.CommandContext(ctx, Binary,
			"-f", strconv.Itoa(i), "-l", strconv.Itoa(i), "-p",
			"-all", fileName, fmt.Sprintf("%s/file", outputdir))
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
