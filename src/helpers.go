package main

import "strings"

func sanitizeFileName(name string) string {
	// Replace slashes with underscores to avoid creating folders
	return strings.ReplaceAll(name, "/", "_")
}

func FirstOrDefault(items []KeyValueType, key string) *KeyValueType {
	for _, item := range items {
		if item.Key == key {
			return &item
		}
	}
	return nil
}
