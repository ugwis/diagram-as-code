package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	xmlFile := "sample.xml"
	slideXML, err := ioutil.ReadFile(xmlFile)
	if err != nil {
		log.Fatal(err)
	}

	slide := &Slide{}
	if err := xml.Unmarshal(slideXML, slide); err != nil {
		log.Fatal(err)
	}

	name := extractDescription(slide)
	name = cleanName(name)
	fmt.Println("Name:", name)

	rId := extractEmbedID(slide)
	fmt.Println("rId:", rId)

	relFile := "sample.xml.rels"
	relXML, err := ioutil.ReadFile(relFile)
	if err != nil {
		log.Fatal(err)
	}

	image := extractImage(relXML, rId)
	fmt.Println("Image:", image)
}

// Slide represents the main structure of a slide in the XML
type Slide struct {
	XMLName xml.Name `xml:"sld"`
	CSld    *CSld    `xml:"cSld"`
}

type CSld struct {
	XMLName xml.Name `xml:"cSld"`
	SpTree  *SpTree  `xml:"spTree"`
}

type SpTree struct {
	XMLName   xml.Name   `xml:"spTree"`
	NvGrpSpPr *NvGrpSpPr `xml:"nvGrpSpPr"`
	Shapes    []Shape    `xml:"any"` // Handle different shapes (sp, pic, etc.)
}

type NvGrpSpPr struct {
	XMLName xml.Name `xml:"nvGrpSpPr"`
	CNvPr   *CNvPr   `xml:"cNvPr"`
}

type CNvPr struct {
	XMLName xml.Name `xml:"cNvPr"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Descr   string   `xml:"descr,attr"`
}

type Shape struct {
	XMLName  xml.Name  `xml:",any"`
	NvSpPr   *NvSpPr   `xml:"nvSpPr"`
	BlipFill *BlipFill `xml:"blipFill"`
}

type NvSpPr struct {
	XMLName xml.Name `xml:"nvSpPr"`
	CNvPr   *CNvPr   `xml:"cNvPr"`
}

type BlipFill struct {
	XMLName xml.Name `xml:"blipFill"`
	Blip    *Blip    `xml:"blip"`
}

type Blip struct {
	XMLName xml.Name `xml:"blip"`
	Embed   string   `xml:"embed,attr"`
}

type Relationships struct {
	XMLName      xml.Name       `xml:"Relationships"`
	Relationship []Relationship `xml:"Relationship"`
}

type Relationship struct {
	XMLName xml.Name `xml:"Relationship"`
	ID      string   `xml:"Id,attr"`
	Target  string   `xml:"Target,attr"`
}

func extractDescription(slide *Slide) string {
	if slide.CSld != nil && slide.CSld.SpTree != nil && slide.CSld.SpTree.NvGrpSpPr != nil {
		return slide.CSld.SpTree.NvGrpSpPr.CNvPr.Descr
	}
	for _, shape := range slide.CSld.SpTree.Shapes {
		if shape.NvSpPr != nil && shape.NvSpPr.CNvPr != nil {
			return shape.NvSpPr.CNvPr.Descr
		}
	}
	return ""
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
		{" resource icon for.*", ""},
		{" instance icon for.*", ""},
		{" storage class icon for.*", ""},
		{" standard category icon.", ""},
		{"A representation of a.*", ""},
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

func extractEmbedID(slide *Slide) string {
	for _, shape := range slide.CSld.SpTree.Shapes {
		if shape.BlipFill != nil && shape.BlipFill.Blip != nil {
			return shape.BlipFill.Blip.Embed
		}
	}
	return ""
}

func extractImage(relXML []byte, rId string) string {
	relationships := &Relationships{}
	if err := xml.Unmarshal(relXML, relationships); err != nil {
		log.Fatal(err)
	}

	for _, rel := range relationships.Relationship {
		if rel.ID == rId {
			return strings.ReplaceAll(rel.Target, "../media/", "")
		}
	}
	return ""
}
