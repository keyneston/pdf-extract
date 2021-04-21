package pdfimages

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Extract(fileName string) (string, error) {
	dir, err := os.MkdirTemp("", "pdf-extract-*")
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := exec.CommandContext(ctx, Binary, "-all", fileName, fmt.Sprintf("%s/file", dir))
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return dir, err
}
