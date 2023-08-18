package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Command string

const (
	GenerateCollection Command = "generate-collection"
	GenerateEnvFile    Command = "generate-env-file"
)

func main() {
	// Define command-line flags
	var inputFile string
	var outputDir string
	flag.StringVar(&inputFile, "input", "", "Path to the input JSON collection file")
	flag.StringVar(&outputDir, "output", "generated-http-files", "Path to the output directory")
	flag.Parse()

	if flag.NArg() < 1 {
		displayDescription()
		return
	}

	command := flag.Arg(0)

	// Check if the input file is provided
	if inputFile == "" {
		fmt.Println("Please provide the path to the input JSON file using the -input flag.")
		return
	}

	if command != string(GenerateCollection) && command != string(GenerateEnvFile) {
		fmt.Println("please provide a valid command. Valid commands are: generate-collection, generate-env-file")
		return
	}

	// Read the JSON collection file
	jsonData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input JSON file:", err)
		return
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	switch command {
	case string(GenerateCollection):
		handleCollectionFile(jsonData, outputDir)
	case string(GenerateEnvFile):
		handleEnvFile(jsonData, outputDir)
	}

	fmt.Println("Conversion complete!")
}

func displayDescription() {
	description := `
	This is a CLI tool for converting Postman json files to HTTP files.
		
	Usage: generate-collection 
	Usage: generate-env-file

	Use --help for more information.
	`
	fmt.Println(description)
}

func handleCollectionFile(jsonData []byte, outputDir string) {
	var collectionData map[string]interface{}
	if err := json.Unmarshal(jsonData, &collectionData); err != nil {
		fmt.Println("Error parsing JSON data:", err)
		os.Exit(1)
	}

	processItems(collectionData["item"].([]interface{}), outputDir)
}

func handleEnvFile(jsonData []byte, outputDir string) {
	envData := Environment{}
	if err := json.Unmarshal(jsonData, &envData); err != nil {
		fmt.Println("Error parsing JSON data:", err)
		os.Exit(1)
	}

	processEnvFile(envData, outputDir)
}
