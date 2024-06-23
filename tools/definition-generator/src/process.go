package main

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type TemplateData struct {
	Url      string
	PptxFile string
}

type Relationship struct {
	ID     string `xml:"Id,attr"`
	Target string `xml:"Target,attr"`
}

type Relationships struct {
	XMLName      xml.Name       `xml:"Relationships"`
	Relationship []Relationship `xml:"Relationship"`
}

func processSlides(files map[string][]byte, slidesPath, outputPath, url, pptxFile string) {
	imageMappings := make(map[string]string)

	for name, data := range files {
		if strings.HasPrefix(name, slidesPath) && strings.HasSuffix(name, ".xml") {
			fmt.Printf("Found slide %s\n", name)
			processSlide(name, data, files, imageMappings)
		}
	}

	generateOutput(outputPath, imageMappings, url, pptxFile)
}

func processSlide(name string, data []byte, files map[string][]byte, imageMappings map[string]string) {
	type Pic struct {
		CNvPr struct {
			Descr string `xml:"descr,attr"`
		} `xml:"nvPicPr>cNvPr"`
		Blip struct {
			Embed string `xml:"embed,attr"`
		} `xml:"blipFill>blip"`
	}

	var pics struct {
		Pic []Pic `xml:"pic"`
	}

	err := xml.Unmarshal(data, &pics)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", &pics)

	for _, pic := range pics.Pic {
		name := cleanName(pic.CNvPr.Descr)
		if name == "" {
			continue
		}
		fmt.Println(name)

		rId := pic.Blip.Embed
		image := findImage(name, rId, files)
		if image == "" {
			continue
		}

		if existing, ok := imageMappings[name]; ok && existing != image {
			for i := 2; i <= 9; i++ {
				newName := fmt.Sprintf("%s(%d)", name, i)
				if _, ok := imageMappings[newName]; !ok {
					imageMappings[newName] = image
					break
				}
			}
		} else {
			imageMappings[name] = image
		}
	}
}

func cleanName(name string) string {
	replacements := []struct {
		old string
		new string
	}{
		{" group.", ""},
		{" Service icon.", ""},
		{" service icon.", ""},
		{" group icon.", ""},
		{" instance icon for the Database category.", ""},
		{" resource icon for", ""},
		{" instance icon for", ""},
		{" storage class icon for", ""},
		{" standard category icon.", ""},
		{"A representation of a", ""},
		{".", ""},
		{"&#10;", ""},
		{"&amp;", "&"},
		{"&#x2013;", "–"},
	}

	for _, r := range replacements {
		name = strings.ReplaceAll(name, r.old, r.new)
	}

	name = strings.TrimSpace(name)
	return name
}

func findImage(slide, rId string, files map[string][]byte) string {
	relPath := slide + ".rels"
	relData, ok := files[relPath]
	if !ok {
		log.Printf("Relationship file not found for slide: %s", slide)
		return ""
	}

	var rels Relationships
	err := xml.Unmarshal(relData, &rels)
	if err != nil {
		log.Fatal(err)
	}

	for _, rel := range rels.Relationship {
		if rel.ID == rId {
			return filepath.Base(rel.Target)
		}
	}
	return ""
}

func generateOutput(outputPath string, imageMappings map[string]string, url, pptxFile string) {
	// Read header
	tmpl, err := template.ParseFiles("header.yaml")
	if err != nil {
		log.Fatal(err)
	}

	data := TemplateData{
		Url:      url,
		PptxFile: pptxFile,
	}

	// Generate main part of the output
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Fatal(err)
	}

	// Append image mappings to the output
	buf.WriteString("\n  # Resource Types\n")

	mappingsFile, err := os.Open("mappings")
	if err != nil {
		log.Fatal(err)
	}
	defer mappingsFile.Close()

	reader := csv.NewReader(mappingsFile)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		typeField := record[0]
		name := record[1]
		fmt.Println(name)
		if name == "" {
			continue
		}

		label := regexp.MustCompile(`\([0-9]\)`).ReplaceAllString(name, "")
		image, ok := imageMappings[name]
		if !ok {
			log.Printf("Not found: %s", name)
			continue
		}

		hasChildren := "false"
		if typeField == "AWS::ECS::Cluster" || typeField == "AWS::EKS::Cluster" || typeField == "AWS::CodePipeline::Pipeline" {
			hasChildren = "true"
		}

		buf.WriteString(fmt.Sprintf(`  %s:
    Type: Resource
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "%s"
    Label:
      Title: "%s"
      Color: "rgba(0, 0, 0, 255)"
    CFn:
      HasChildren: %s

`, typeField, image, label, hasChildren))
	}

	buf.WriteString("\n  # Presets\n")

	for key, image := range imageMappings {
		name := regexp.MustCompile(`\([0-9]\)`).ReplaceAllString(key, "")
		buf.WriteString(fmt.Sprintf(`  "%s":
    Type: Preset
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "%s"
    Label:
      Title: "%s"
      Color: "rgba(0, 0, 0, 255)"

`, key, image, name))
	}

	err = ioutil.WriteFile(outputPath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
