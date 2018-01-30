package main

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

const MANIFEST_NAME = "manifest-%v.yaml"

type ManifestDetails struct {
	Spacename        string
	UserConfig       *UserConfig
	T                *template.Template
	ManifestTemplate string
	Writer           io.Writer
	OutputDirectory  string
}

func (m *ManifestDetails) WriteManifestToBuffer() error {
	return m.T.ExecuteTemplate(m.Writer, m.ManifestTemplate, m)
}

func GenerateManifest(m *ManifestDetails) error {

	file, err := os.Create(m.OutputDirectory + "/" + fmt.Sprintf(MANIFEST_NAME, m.Spacename))
	if err != nil {
		return err
	}
	defer file.Close()

	m.T = template.Must(template.New("manifest").ParseFiles(m.ManifestTemplate))
	m.Writer = file

	return m.WriteManifestToBuffer()
}
