package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func processItems(items []interface{}, currentDir string) {
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		requestName, _ := itemMap["name"].(string)
		nestedItems, nestedOK := itemMap["item"].([]interface{})

		// Check if the item contains a "request" field
		if request, ok := itemMap["request"].(map[string]interface{}); ok {
			// This is a request item
			if requestName != "" {
				httpFileName := sanitizeFileName(requestName) + ".http"
				httpFilePath := filepath.Join(currentDir, httpFileName)
				processRequestItem(requestName, request, httpFilePath)
				fmt.Printf("Converted '%s' to '%s'\n", requestName, httpFilePath)
			}
		} else if nestedOK {
			// This is a folder with nested items
			// Create a subdirectory for the folder
			subDir := filepath.Join(currentDir, sanitizeFileName(requestName))
			if err := os.MkdirAll(subDir, 0755); err != nil {
				fmt.Println("Error creating subdirectory:", err)
				return
			}
			processItems(nestedItems, subDir)
		}
	}
}

func processRequestItem(requestName string, request map[string]interface{}, filePath string) {
	requestMethod, _ := request["method"].(string)
	requestURL, _ := request["url"].(map[string]interface{})["raw"].(string)
	requestQuery, queryOK := request["url"].(map[string]interface{})["query"].([]interface{})
	requestVariables, _ := request["url"].(map[string]interface{})["variable"].([]interface{})
	requestBody, bodyOK := request["body"].(map[string]interface{})

	requestURL = useCorrectInterpolationForPlaceholders(requestURL)

	// Construct .http file content
	httpFileContent := fmt.Sprintf("### %s\n", requestName) // Start region
	httpFileContent += generateDescription(requestQuery)
	httpFileContent += generateDescription(requestVariables)
	httpFileContent += fmt.Sprintf("\n\n%s %s", requestMethod, removeQueryParamsFromRequestURL(requestURL))

	if queryOK && len(requestQuery) > 0 {
		httpFileContent += formatQueryParameters(requestQuery)
	}

	authData, authOK := request["auth"].(map[string]interface{})
	if authOK {
		httpFileContent += generateAuthLine(authData)
	}

	httpFileContent += genereateHeaderLines(request["header"].([]interface{}))
	if bodyOK {
		httpFileContent += generateBodyContent(requestBody)
	}

	// Write the .http file
	if err := os.WriteFile(filePath, []byte(httpFileContent), 0644); err != nil {
		fmt.Println("Error writing .http file:", err)
		return
	}
}

func useCorrectInterpolationForPlaceholders(url string) string {
	re := regexp.MustCompile(`/(?::(\w+))`)
	return re.ReplaceAllString(url, `/{{${1}}}`)
}

func generateDescription(parameters []interface{}) string {
	if (len(parameters)) == 0 {
		return ""
	}

	descriptionLines := []string{}
	for _, param := range parameters {
		paramMap := param.(map[string]interface{})
		paramKey, _ := paramMap["key"].(string)
		paramValue, _ := paramMap["value"].(string)
		paramDescription, _ := paramMap["description"].(string)
		paramDescription = strings.ReplaceAll(paramDescription, "\n", " ")
		descriptionLines = append(descriptionLines, fmt.Sprintf("\n// %s", paramDescription))
		descriptionLines = append(descriptionLines, fmt.Sprintf("\n@%s = %s\n", paramKey, paramValue))
	}

	return strings.Join(descriptionLines, "")
}

func removeQueryParamsFromRequestURL(rawURL string) string {
	lines := strings.Split(rawURL, "?")
	return lines[0]
}

func formatQueryParameters(queryParameters []interface{}) string {
	if len(queryParameters) > 0 {
		queryLines := []string{}
		for index, param := range queryParameters {
			paramMap := param.(map[string]interface{})
			paramKey, _ := paramMap["key"].(string)
			queryPrefix := getQueryPrefixFromIndex(index)
			queryLines = append(queryLines, fmt.Sprintf("\n  %s%s={{%s}}", queryPrefix, paramKey, paramKey))
		}
		return strings.Join(queryLines, "")
	}
	return ""
}

func getQueryPrefixFromIndex(index int) string {
	if index == 0 {
		return "?"
	}
	return "&"
}

func genereateHeaderLines(headers []interface{}) string {
	if len(headers) == 0 {
		return ""
	}

	headerLines := []string{}
	for _, header := range headers {
		headerMap := header.(map[string]interface{})
		headerKey, _ := headerMap["key"].(string)
		headerValue, _ := headerMap["value"].(string)
		headerLines = append(headerLines, fmt.Sprintf("\n%s: %s", headerKey, headerValue))
	}

	return strings.Join(headerLines, "")
}

func generateAuthLine(authData map[string]interface{}) string {
	var auth Auth
	authJSON, err := json.Marshal(authData)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(authJSON, &auth)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if auth.Type == "bearer" {
		return fmt.Sprintf("\nAuthorization: Bearer %s", auth.Bearer[0].Value)
	}

	if auth.Type == "basic" {
		username := FirstOrDefault(auth.Basic, "username")
		password := FirstOrDefault(auth.Basic, "password")
		return fmt.Sprintf("\nAuthorization: Basic %s %s", username.Value, password.Value)
	}

	if auth.Type == "apikey" {
		apiKeyPlacementItem := FirstOrDefault(auth.ApiKey, "in")
		key := FirstOrDefault(auth.ApiKey, "key")
		value := FirstOrDefault(auth.ApiKey, "value")
		if apiKeyPlacementItem != nil && apiKeyPlacementItem.Value == "header" {
			return fmt.Sprintf("\n%s: %s", key.Value, value.Value)
		}

		// for we don't know where the api key should be placed, so we comment it out
		return fmt.Sprintf("\n//TODO: Verify placement of apiKey in header or query: %s=%s", key.Value, value.Value)
	}

	return ""
}

func generateBodyContent(body map[string]interface{}) string {
	bodyMode, _ := body["mode"].(string)
	if bodyMode == "raw" {
		bodyRaw, _ := body["raw"].(string)
		return fmt.Sprintf("\n\n%s", bodyRaw)
	}
	return ""
}
