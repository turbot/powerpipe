package bats_generator

import (
	"fmt"
	"log"
)

func CallGenerateBatsFile() {
	batsTemplatePath := "tests/acceptance/test_data/templates/mod_test_template.bats.tmpl"
	jsonPath := "tests/acceptance/test_data/templates/test_cases.json"
	outputPath := "tests/acceptance/test_files/mod.bats"

	err := GenerateBatsFile(batsTemplatePath, jsonPath, outputPath)
	if err != nil {
		log.Fatalf("Failed to generate Bats test file: %s", err)
	}

	fmt.Println("Bats test file generated successfully!")
}
