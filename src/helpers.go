package main

import "strings"

func sanitizeFileName(name string) string {
	// Replace slashes with underscores to avoid creating folders
	temp := strings.ReplaceAll(name, "/", "_")

	// Remove colons to avoid creating files with invalid names
	return strings.ReplaceAll(temp, ":", "")
}

func FirstOrDefault(items []KeyValueType, key string) *KeyValueType {
	for _, item := range items {
		if item.Key == key {
			return &item
		}
	}
	return nil
}
