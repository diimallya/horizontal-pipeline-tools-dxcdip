package main

import (
	"bytes"
	"strings"
	"testing"
	"text/template"
)

var testTemplate = `---
applications:
- name: {{.Spacename}}-admin
  memory: 128M
  no-hostname: true
  no-route: true
  services:
  - {{.Spacename}}-cloudant
`

var testTemplateWithDb = `---
applications:
- name: {{.Spacename}}-admin
  memory: 128M
  no-hostname: true
  no-route: true
	services:
	{{if .UserConfig.Database }}- {{.Spacename}}-{{ .UserConfig.Database.Name }}{{end}}
  - {{.Spacename}}-cloudant
`

func TestGenerateManifest(t *testing.T) {

	b := bytes.NewBuffer([]byte{})
	templateName := "testTemplate"

	template := template.Must(template.New(templateName).Parse(testTemplate))
	manifestDetails := &ManifestDetails{
		Spacename:        "test123",
		T:                template,
		ManifestTemplate: templateName,
		Writer:           b,
	}

	err := manifestDetails.WriteManifestToBuffer()
	if err != nil {
		t.Fatal(err)
	}

	if strings.Count(b.String(), "test123") != 2 {
		t.Fatal("Expected test123 to be mentioned twice but was not")
	}
}

func TestGenerateDbManifest(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	templateName := "testTemplate"

	template := template.Must(template.New(templateName).Parse(testTemplateWithDb))
	manifestDetails := &ManifestDetails{
		Spacename:        "test1234",
		T:                template,
		ManifestTemplate: templateName,
		Writer:           b,
		UserConfig: &UserConfig{
			Database: &RelationalDBConfig{
				Name: "compose-for-mysql",
			},
		},
	}

	err := manifestDetails.WriteManifestToBuffer()
	if err != nil {
		t.Fatal(err)
	}

	if strings.Count(b.String(), "test1234") != 3 {
		t.Fatal("Expected test1234 to be mentioned thrice but was not")
	}

	if !strings.Contains(b.String(), "test1234-compose-for-mysql") {
		t.Fatalf("Expected to see string 'test1234-compose-for-mysql' but did not get it: \n%v", b.String())
	}
}
