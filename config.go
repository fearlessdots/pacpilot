package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	// External modules
	//color "github.com/gookit/color"
)

//
//// PROGRAM CONFIGURATION
//

type Program struct {
	name                string
	nameAscii           string
	version             string
	exec                string
	shortDescription    string
	longDescription     string
	defaultShell        string
	dataDir             string
	reposDir            string
	templatesDir        string
	targetsTemplatesDir string
	reposTemplatesDir   string
	indentLevel         int
}

func getDefaultShellAbsolutePath(shellName string) string {
	// Get shell absolute path using `which`
	cmd := exec.Command("which", shellName)

	output, err := cmd.CombinedOutput()
	outputString := string(output)
	outputString = strings.TrimLeft(outputString, "\n")
	outputString = strings.TrimRight(outputString, "\n")

	if err != nil {
		showError(fmt.Sprintf("Failed to get absolute path to default shell '%v' -> %v", shellName, outputString), 0)
		finishProgram(1)
	}

	return outputString
}

func initializeDefaultProgram(dataDir string) Program {
	// PROGRAM NAME
	programName := "pacpilot"

	// PROGRAM NAME ASCII
	programNameAscii := " ____            ____  _ _       _\n" +
		"|  _ \\ __ _  ___|  _ \\(_) | ___ | |_\n" +
		"| |_) | _` |/ __| |_) | | |/ _ \\| __|\n" +
		"|  __/ (_| | (__|  __/| | | (_) | |_\n" +
		"|_|   \\__,_|\\___|_|   |_|_|\\___/ \\__|\n"

	// PROGRAM VERSION
	programVersion := "0.1.0"

	// PROGRAM EXEC
	programExec := os.Args[0] // Path for program executable

	// DESCRIPTIONS (SHORT AND LONG)
	programShortDescription := "A user-friendly tool for efficiently managing and serving custom Pacman repositories. You can quickly create and manage repositories and targets using templates, and take advantage of hooks to customize the behavior of the program."
	programLongDescription := fmt.Sprintf("%v is a user-friendly tool for efficiently managing and serving custom Pacman repositories. You can quickly create and manage repositories and targets using templates, and take advantage of hooks to customize the behavior of the program.", programName)

	// DEFAULT SHELL
	programDefaultShellName := "sh" // Should work on all Unix systems (Linux, Android, ...)
	programDefaultShellPath := getDefaultShellAbsolutePath(programDefaultShellName)

	// DIRECTORIES
	reposDir := dataDir + "/repos"
	templatesDir := dataDir + "/templates"
	targetsTemplatesDir := dataDir + "/templates" + "/targets"
	reposTemplatesDir := dataDir + "/templates" + "/repos"

	// INDENT LEVEL
	indentLevel := 0

	return Program{
		name:                programName,
		nameAscii:           programNameAscii,
		version:             programVersion,
		exec:                programExec,
		shortDescription:    programShortDescription,
		longDescription:     programLongDescription,
		defaultShell:        programDefaultShellPath,
		dataDir:             dataDir,
		reposDir:            reposDir,
		templatesDir:        templatesDir,
		targetsTemplatesDir: targetsTemplatesDir,
		reposTemplatesDir:   reposTemplatesDir,
		indentLevel:         indentLevel,
	}
}

func getRootDirectory() string {
	// Check if the "PREFIX" environment variable is set
	prefix := os.Getenv("PREFIX")
	if prefix != "" {
		// Verify if the given prefix string ends with the directory suffix "/usr".
		// If it does, the program will proceed to strip it from the prefix.
		if strings.HasSuffix(prefix, "/usr") {
			return strings.TrimSuffix(prefix, "/usr")
		}
		return prefix
	} else {
		return "/"
	}
}

func displayProgramInfo(program Program) {
	showText(program.nameAscii, program.indentLevel)

	showText("Version: "+green.Sprintf(program.version), program.indentLevel)

	space()

	showText("Running on "+lightCopper.Sprintf(runtime.GOOS+"/"+runtime.GOARCH)+". Built with "+runtime.Version()+" using "+runtime.Compiler+" as compiler.", program.indentLevel)

	space()

	showText("This program is licensed under the GNU General Public License v3.0 (GPL-3.0).\nPlease refer to the LICENSE file for more information.", program.indentLevel)
}
