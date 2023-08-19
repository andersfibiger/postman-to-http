package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Command string

const (
	ConvertCollection  Command = "convert-collection"
	ConvertEnvironment Command = "convert-environment"
)

func main() {
	// Define command-line flags
	var outputDir string
	flag.StringVar(&outputDir, "output", "out", "Path to the output directory")
	flag.Parse()

	if flag.NArg() < 1 {
		displayDescription()
		return
	}

	command := flag.Arg(0)
	inputFile := flag.Arg(1)

	if command != string(ConvertCollection) && command != string(ConvertEnvironment) {
		fmt.Println("please provide a valid command. Valid commands are: convert-collection, convert-environment")
		return
	}

	if inputFile == "" {
		fmt.Println("Please provide the path to the input JSON file as an argument")
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
	case string(ConvertCollection):
		handleCollectionFile(jsonData, outputDir)
	case string(ConvertEnvironment):
		handleEnvFile(jsonData, outputDir)
	}

	fmt.Println("Conversion complete!")
}

func displayDescription() {
	description := `
	This is a CLI tool for converting Postman json files to HTTP files.
		
	Usage: postmanTohttp [-output] <command> [args]
	
	Commands:
		convert-collection		Converts a Postman collection file to HTTP files 
		convert-environment		Converts a Postman environment file to a .env file

	Args:
		Path to the input JSON file

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
