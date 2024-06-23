package main

import (
	"log"
)

func main() {
	// Step 1: Fetch URL
	url := fetchURL()
	if url == "" {
		log.Fatal("Failed to fetch URL")
	}

	// Step 2: Download the file
	zipData := downloadFile(url)

	// Step 3: Unzip the file
	unzippedData := unzip(zipData)

	// Step 4: Find the pptx file
	pptxFile := findPPTXFile(unzippedData)
	if pptxFile == "" {
		log.Fatal("Failed to find pptx file")
	}

	// Step 5: Unzip the pptx file
	pptxContent := unzip(unzippedData[pptxFile])

	// Step 6: Process slides and generate output
	outputPath := "definition-for-aws-icons-light.yaml"
	slidesPath := "ppt/slides"
	processSlides(pptxContent, slidesPath, outputPath, url, pptxFile)
}
