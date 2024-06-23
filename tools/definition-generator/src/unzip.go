package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"log"
	"strings"
)

func unzip(data []byte) map[string][]byte {
	unzipped := make(map[string][]byte)
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer rc.Close()

		content, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Fatal(err)
		}

		unzipped[file.Name] = content
	}
	return unzipped
}

func findPPTXFile(files map[string][]byte) string {
	for name := range files {
		if strings.HasSuffix(name, ".pptx") {
			return name
		}
	}
	return ""
}
