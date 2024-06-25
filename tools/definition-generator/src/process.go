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
	"sort"
	"strings"
	"text/template"
	"unicode"
)

type TemplateData struct {
	Url      string
	PptxFile string
}

func removeControlCharacter(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, s)
}

func processSlides(files map[string][]byte, slidesPath, outputPath, url, pptxFile string) {
	imageMappings := make(map[string]string)

	for name, data := range files {
		if strings.HasPrefix(name, slidesPath) && strings.HasSuffix(name, ".xml") {
			fmt.Printf("Found slide %s\n", name)
			processSlide(filepath.Base(name), data, files, imageMappings)
		}
	}

	generateOutput(outputPath, imageMappings, url, pptxFile)
}

func processSlide(slideFile string, data []byte, files map[string][]byte, imageMappings map[string]string) {
	slide := &Slide{}
	if err := xml.Unmarshal(data, slide); err != nil {
		log.Fatal(err)
	}

	names := []string{}
	for _, sp := range slide.CSld.SpTree.Sps {
		if sp.TxBody == nil {
			continue
		}
		if sp.TxBody.P == nil {
			continue
		}
		if sp.TxBody.P.Rs == nil {
			continue
		}
		name := ""
		for _, r := range sp.TxBody.P.Rs {
			if r.T != "" {
				name = name + " " + r.T
			}
		}
		fmt.Println(name)
		name = removeControlCharacter(strings.TrimSpace(name))
		names = append(names, name)

	}
	sort.Slice(names, func(i, j int) bool { return len(names[i]) > len(names[j]) })
	for _, pic := range slide.CSld.SpTree.Pics {
		if pic.NvPicPr != nil && pic.NvPicPr.CNvPr != nil {
			rId := pic.BlipFill.Blip.Embed
			desc := removeControlCharacter(pic.NvPicPr.CNvPr.Descr)

			name := ""
			for _, x := range names {
				if strings.Contains(desc, x) {
					name = x
				}
			}

			fmt.Println("Name:", name)
			fmt.Println("rId:", rId)
			image := findImage(slideFile, rId, files)
			if image == "" {
				continue
			}
			fmt.Println("Image:", image)

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
}

func findImage(slide, rId string, files map[string][]byte) string {
	relPath := "ppt/slides/_rels/" + slide + ".rels"
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
