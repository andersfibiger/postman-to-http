package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment struct {
	Name   string `json:"name"`
	Values []struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Type    string `json:"type"`
		Enabled bool   `json:"enabled"`
	} `json:"values"`
}

func processEnvFile(envData Environment, outputDir string) {
	envContent := ""
	for _, value := range envData.Values {
		envContent += fmt.Sprintf("%s=%s\n", value.Key, value.Value)
	}

	envFileName := fmt.Sprintf("%s.env", strings.ToLower(envData.Name))
	if err := os.WriteFile(filepath.Join(outputDir, envFileName), []byte(envContent), 0644); err != nil {
		fmt.Println("Error writing .env file:", err)
		return
	}
}
