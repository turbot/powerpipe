package bats_generator

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"text/template"
)

type TestCase struct {
	Description             string `json:"description"`
	Name                    string `json:"name"`
	SetupComment            string `json:"setup_comment"`
	Setup                   string `json:"setup"`
	OperationComment        string `json:"operation_comment"`
	Operation               string `json:"operation"`
	Cmd                     string `json:"cmd"`
	ExpectedOutput          string `json:"expected_output"`
	ExpectedFolderStructure string `json:"expected_folder_structure"`
	ExpectedModFiles        string `json:"expected_mod_files"`
}

func GenerateBatsFile(batsTemplatePath string, jsonPath string, outputPath string) error {
	// Read JSON file
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var testCases []TestCase
	if err := json.Unmarshal(byteValue, &testCases); err != nil {
		return err
	}

	// Read the template file
	templateContent, err := ioutil.ReadFile(batsTemplatePath)
	if err != nil {
		return err
	}

	// Create and parse the template
	t := template.Must(template.New("bats").Parse(string(templateContent)))

	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Execute the template with the test cases
	if err := t.Execute(outFile, testCases); err != nil {
		return err
	}

	return nil
}
