package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

func main() {
	f, err := os.Open("input.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	m := make(map[string]string)
	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		url := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])
		if url == "" {
			continue
		}
		if name == "" {
			name = uuid.New().String()
		}
		m[name] = url
	}

	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	gold := color.RGBA{R: 255, G: 215, B: 0, A: 255}

	os.MkdirAll("output", os.ModePerm)

	for name, url := range m {
		path := filepath.Join("output", name+".png")
		err := qrcode.WriteColorFile(url, qrcode.Medium, 1000, gold, blue, path)
		if err != nil {
			fmt.Printf("failed to generate QR for %s: %v\n", name, err)
		}
	}

	outF, err := os.Create("end.csv")
	if err != nil {
		panic(err)
	}
	defer outF.Close()

	w := csv.NewWriter(outF)
	for name, url := range m {
		err := w.Write([]string{url, name})
		if err != nil {
			panic(err)
		}
	}
	w.Flush()
}
