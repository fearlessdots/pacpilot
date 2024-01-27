package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"os"
	"path/filepath"
	"strings"
	// External modules
)

//
//// FILES/DIRECTORIES
//

func fileIsHidden(name string) bool {
	_, file := filepath.Split(name)
	return strings.HasPrefix(file, ".")
}

func filterHiddenFilesAndDirectories(unfilteredFiles []os.FileInfo) []os.FileInfo {
	var filteredFiles []os.FileInfo

	for _, element := range unfilteredFiles {
		if !fileIsHidden(element.Name()) {
			filteredFiles = append(filteredFiles, element)
		}
	}

	return filteredFiles
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	var value float64
	var unit string

	switch {
	case bytes >= TB:
		value = float64(bytes) / TB
		unit = "T"
	case bytes >= GB:
		value = float64(bytes) / GB
		unit = "G"
	case bytes >= MB:
		value = float64(bytes) / MB
		unit = "M"
	case bytes >= KB:
		value = float64(bytes) / KB
		unit = "K"
	default:
		return fmt.Sprintf("%d", bytes) // bytes
	}

	// Using %.2f to only keep 2 decimal places
	return fmt.Sprintf("%.2f%s", value, unit)
}
