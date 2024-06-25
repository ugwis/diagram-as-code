package main

import (
	"encoding/xml"
	"strings"
)

type Slide struct {
	XMLName xml.Name `xml:"sld"`
	CSld    *CSld    `xml:"cSld"`
}

type CSld struct {
	XMLName xml.Name `xml:"cSld"`
	SpTree  *SpTree  `xml:"spTree"`
}

type SpTree struct {
	XMLName xml.Name `xml:"spTree"`
	Pics    []Pic    `xml:"pic"`
	Sps     []Sp     `xml:"sp"`
}

type Pic struct {
	XMLName  xml.Name  `xml:"pic"`
	NvPicPr  *NvPicPr  `xml:"nvPicPr"`
	BlipFill *BlipFill `xml:"blipFill"`
}

type NvPicPr struct {
	XMLName xml.Name `xml:"nvPicPr"`
	CNvPr   *CNvPr   `xml:"cNvPr"`
}

type CNvPr struct {
	XMLName xml.Name `xml:"cNvPr"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Descr   string   `xml:"descr,attr"`
}

type BlipFill struct {
	XMLName xml.Name `xml:"blipFill"`
	Blip    *Blip    `xml:"blip"`
}

type Blip struct {
	XMLName xml.Name `xml:"blip"`
	Embed   string   `xml:"embed,attr"`
}

type Sp struct {
	XMLName xml.Name `xml:"sp"`
	TxBody  *TxBody  `xml:"txBody"`
}

type TxBody struct {
	XMLName xml.Name `xml:"txBody"`
	P       *P       `xml:"p"`
}

type P struct {
	XMLName xml.Name `xml:"p"`
	Rs      []R      `xml:"r"`
}

type R struct {
	XMLName xml.Name `xml:"r"`
	T       string   `xml:"t"`
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
