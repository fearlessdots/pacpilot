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
//// FINISH/RESPONSE/STATE CONTROL
//

func finishProgram(code int) {
	os.Exit(code)
}

type functionResponse struct {
	exitCode    int
	message     string
	logLevel    string
	indentLevel int
}

func incrementProgramIndentLevel(program Program, count int) Program {
	program.indentLevel = program.indentLevel + 1

	return program
}

func decrementProgramIndentLevel(program Program, count int) Program {
	program.indentLevel = program.indentLevel - 1

	return program
}

func handleFunctionResponse(response functionResponse, finishProgramAfter bool) {
	if response.exitCode != 0 {
		if response.logLevel == "attention" {
			showAttention(fmt.Sprintf("> "+response.message), response.indentLevel)

			if finishProgramAfter == true {
				space()

				finishProgram(response.exitCode)
			}
		} else if response.logLevel == "error" {
			showError(fmt.Sprintf("> Error: "+response.message), response.indentLevel)

			if finishProgramAfter == true {
				space()

				finishProgram(response.exitCode)
			}
		}
	} else {
		if response.logLevel == "attention" {
			showAttention(fmt.Sprintf("> "+response.message), response.indentLevel)
		} else if response.logLevel == "success" {
			showSuccess(fmt.Sprintf("> "+response.message), response.indentLevel)
		}
	}
}
