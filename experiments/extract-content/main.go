package main

import (
	"fmt"
	"log"
	"os"

	"github.com/heussd/pdftotext-go"
)

func main() {
	// Read PDF file
	pdfPath := "../content/Wahlprogramm_2024_Zukunft_Erleben_29.01.2024.pdf"
	pdf, err := os.ReadFile(pdfPath)
	if err != nil {
		log.Fatalf("Failed to open PDF: %v\n", err)
	}

	pages, err := pdftotext.Extract(pdf)
	if err != nil {
		log.Fatalf("Failed to extract text from PDF: %v", err)
	}

	// Create output file
	outputFile, err := os.Create("../chat/content/Wahlprogramm_2024_Zukunft_Erleben_29.01.2024.txt")
	if err != nil {
		log.Fatalf("Failed to create output file: %v\n", err)
	}
	defer outputFile.Close()

	// Write the text content of each page to file
	for _, page := range pages {
		_, err := fmt.Fprintf(outputFile, "Page %d:\n%s\n\n", page.Number, page.Content)
		if err != nil {
			log.Fatalf("Failed to write to file: %v\n", err)
		}
	}

}
