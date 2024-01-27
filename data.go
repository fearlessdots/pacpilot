package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"os"
	// External modules
)

//
//// DATA/CONFIGURATION
//

func verifyDataDirectory(printSectionTitle bool, program Program) functionResponse {
	if printSectionTitle == true {
		showInfoSectionTitle("Verifying data directory", program.indentLevel)
	}

	createdDirectories := false
	if _, err := os.Stat(program.dataDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> Data directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.dataDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create data directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}
	if _, err := os.Stat(program.reposDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> Repos directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.reposDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create repos directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}
	if _, err := os.Stat(program.templatesDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> Templates directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.templatesDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create templates directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}
	if _, err := os.Stat(program.reposTemplatesDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> Repos templates directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.reposTemplatesDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create repos templates directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}
	if _, err := os.Stat(program.targetsTemplatesDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> Targets templates directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.targetsTemplatesDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create targets templates directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}

	response := functionResponse{
		exitCode:    0,
		logLevel:    "success",
		indentLevel: program.indentLevel + 1,
	}

	if createdDirectories == true {
		response.message = "Finished"
	} else {
		response.message = "Passed"
	}

	return response
}
