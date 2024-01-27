package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// External modules
	survey "github.com/AlecAivazis/survey/v2"
	ptywrapper "github.com/fearlessdots/ptywrapper"
	terminal "golang.org/x/crypto/ssh/terminal"
	//pty "github.com/creack/pty" => indirectly imported by the 'ptywrapper' module
)

//
//// USER INPUT
//

func askConfirmation(message string, program Program) bool {
	confirm := false

	prompt := &survey.Confirm{
		Message: fmt.Sprintf(returnText(message, program.indentLevel)),
	}
	survey.AskOne(prompt, &confirm)

	return confirm
}

//
//// TERMINAL MANAGEMENT
//

type terminalDimensions struct {
	height int
	width  int
}

func getTerminalDimensions(program Program) (terminalDimensions, functionResponse) {
	// Get the file descriptor for the standard output
	fd := int(os.Stdout.Fd())

	// Check if the file descriptor is associated with a terminal
	if !terminal.IsTerminal(fd) {
		return terminalDimensions{}, functionResponse{
			exitCode:    1,
			message:     "Standard output is not a terminal",
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}

	// Retrieve the terminal size
	width, height, err := terminal.GetSize(fd)
	if err != nil {
		return terminalDimensions{}, functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to get terminal dimensions -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}

	return terminalDimensions{height: height, width: width}, functionResponse{exitCode: 0}
}

//
//// COMMAND EXECUTION
//

func runHook(hookPath string, env map[string]string, printOutput bool, printFinished bool, showRulers bool, printEntryCmd bool, printAlerts bool, program Program) (ptywrapper.Command, functionResponse) {
	// Verify if hook exists
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return ptywrapper.Command{}, functionResponse{
			exitCode:    1,
			message:     "Hook not found",
			logLevel:    "attention",
			indentLevel: program.indentLevel,
		}
	}

	// Get the current environment
	currentEnv := os.Environ()

	// Modify the environment variables
	if env != nil {
		for key, value := range env {
			currentEnv = append(currentEnv, key+"="+value)
		}
	}

	// Verify if hook has custom entry command
	var entryCommand string
	var entryArgs []string

	customEntryFilePath := hookPath + ".entry"
	if _, err := os.Stat(customEntryFilePath); err == nil {
		contents, err := ioutil.ReadFile(customEntryFilePath)
		if err != nil {
			return ptywrapper.Command{}, functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to read custom entry configuration file -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel,
			}
		}
		entryCommandString := string(contents)
		entryCommandString = strings.TrimLeft(entryCommandString, "\n")
		entryCommandString = strings.TrimRight(entryCommandString, "\n")
		entryCommandSlice := strings.Split(entryCommandString, " ")

		if len(entryCommandSlice) > 1 {
			entryCommand = entryCommandSlice[0]
			entryArgs = append(entryCommandSlice[1:], hookPath)
		} else {
			entryCommand = entryCommandSlice[0]
			entryArgs = []string{hookPath}
		}
	} else {
		entryCommand = program.defaultShell
		entryArgs = []string{hookPath}
	}

	// Run the hook (using the ptywrapper module)
	if printEntryCmd == true {
		showInfoSectionTitle(fmt.Sprintf("Entry command: %s", paleLime.Sprintf(entryCommand)), program.indentLevel+1)
	}

	if printOutput == false && printAlerts == true {
		showText(gray.Sprintf("> The command will run silently. Interactive commands may not function properly. If necessary, press Ctrl+C or use the program-specific shortcut to quit."), program.indentLevel+1)
	}

	cmd := &ptywrapper.Command{
		Entry:   entryCommand,
		Args:    entryArgs,
		Env:     currentEnv,
		Discard: !printOutput,
	}

	if showRulers == true {
		hr("-", 0.5, incrementProgramIndentLevel(program, 1))
	}

	completedCmd, err := cmd.RunInPTY()
	if err != nil {
		return ptywrapper.Command{}, functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to execute hook -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}

	if showRulers == true {
		hr("-", 0.5, incrementProgramIndentLevel(program, 1))
	}

	var logLevel string
	var message string

	if completedCmd.ExitCode != 0 {
		logLevel = "error"
		message = fmt.Sprintf("Failed to execute hook: exit code %v", completedCmd.ExitCode)
	} else {
		logLevel = "success"

		if printFinished == true {
			message = "Finished"
		} else {
			message = ""
		}
	}

	return completedCmd, functionResponse{
		exitCode:    completedCmd.ExitCode,
		message:     message,
		indentLevel: program.indentLevel,
		logLevel:    logLevel,
	}
}
