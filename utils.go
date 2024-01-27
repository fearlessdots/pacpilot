package main

// ____            ____  _ _       _
// |  _ \ __ _  ___|  _ \(_) | ___ | |_
// | |_) / _` |/ __| |_) | | |/ _ \| __|
// |  __/ (_| | (__|  __/| | | (_) | |_
// |_|   \__,_|\___|_|   |_|_|\___/ \__|
//
//
// This code implements a versatile set of utilities meticulously designed to simplify
// tasks during hook execution with PacPilot. These purpose-built utilities enhance
// PacPilot's functionality by providing convenient methods for common operations,
// making hook execution more seamless and efficient.

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	// External modules
	ptywrapper "github.com/fearlessdots/ptywrapper"
	color "github.com/gookit/color"
)

//
//// DISPLAY
//

func utilsMsg(colorHEX string, msg string, program Program) {
	// Create the color object
	color := color.HEX(colorHEX)

	// Unquote the message (to render new lines)
	msg, _ = strconv.Unquote(`"` + msg + `"`)

	// Print the message using the specified color
	printColoredMessage(color, msg, program.indentLevel)
}

func utilsShowAttention(msg string, program Program) {
	showAttention(msg, program.indentLevel)
}

func utilsShowError(msg string, program Program) {
	showError(msg, program.indentLevel)
}

func utilsShowSuccess(msg string, program Program) {
	showSuccess(msg, program.indentLevel)
}

func utilsShowSection(msg string, program Program) {
	showInfoSectionTitle(msg, program.indentLevel)
}

func utilsHr(char string, factor float64, program Program) {
	hr(char, factor, program)
}

//
//// USER INPUT
//

func utilsAskConfirmation(msg string, program Program) {
	var confirm bool

	confirm = askConfirmation(msg, program)

	if confirm {
		finishProgram(0)
	} else {
		finishProgram(1)
	}
}

//
//// SSH
//

func utilsSSHAgentStart(privateKeyPath string, tempDir string, program Program) functionResponse {
	pidFile := filepath.Join(tempDir, "sshagent.pid")
	sockFile := filepath.Join(tempDir, "sshauth.sock")

	if _, err := os.Stat(pidFile); err == nil {
		return functionResponse{
			exitCode:    1,
			logLevel:    "error",
			message:     fmt.Sprintf("File 'sshagent.pid' already found at %v. Please, remove it manually.", pidFile),
			indentLevel: program.indentLevel,
		}
	}
	if _, err := os.Stat(sockFile); err == nil {
		return functionResponse{
			exitCode:    1,
			logLevel:    "error",
			message:     fmt.Sprintf("File 'sshauth.sock' already found at %v. Please, remove it manually.", pidFile),
			indentLevel: program.indentLevel,
		}
	}

	//
	////
	//

	showInfo("Starting ssh-agent process", program.indentLevel)

	cmd := &ptywrapper.Command{
		Entry:   "ssh-agent",
		Args:    []string{"-s"},
		Discard: true,
	}

	completedCmd, err := cmd.RunInPTY()
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to execute command -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	}

	if completedCmd.ExitCode != 0 {
		return functionResponse{
			exitCode:    completedCmd.ExitCode,
			message:     fmt.Sprintf("Failed to execute command: exit code %v", completedCmd.ExitCode),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		response := functionResponse{
			exitCode:    0,
			message:     "Finished",
			logLevel:    "success",
			indentLevel: program.indentLevel + 1,
		}
		handleFunctionResponse(response, true)
	}

	outputLines := strings.Split(string(completedCmd.Output), "\n")
	sshAuthSock := strings.Split(strings.Split(outputLines[0], "=")[1], ";")[0]
	sshAgentPid := strings.Split(strings.Split(outputLines[1], "=")[1], ";")[0]

	sshAgentEnvVars := make(map[string]string)
	sshAgentEnvVars[pidFile] = sshAgentPid
	sshAgentEnvVars[sockFile] = sshAuthSock

	for key, value := range sshAgentEnvVars {
		err := ioutil.WriteFile(key, []byte(value), 0644)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to write to file %v -> %v", key, err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 1,
			}
		}
	}

	os.Setenv("SSH_AGENT_PID", sshAgentPid)
	os.Setenv("SSH_AUTH_SOCK", sshAuthSock)

	//
	////
	//

	space()

	showInfo("Adding private key (a passphrase may be needed)", program.indentLevel)

	cmd = &ptywrapper.Command{
		Entry:   "ssh-add",
		Args:    []string{privateKeyPath},
		Discard: false,
	}

	completedCmd, err = cmd.RunInPTY()
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to execute command -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	}

	if completedCmd.ExitCode != 0 {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to execute command: exit code %v", completedCmd.ExitCode),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		return functionResponse{
			exitCode:    0,
			message:     "Finished",
			logLevel:    "success",
			indentLevel: program.indentLevel + 1,
		}
	}
}

func utilsSSHAgentStop(tempDir string, program Program) functionResponse {
	pidFile := filepath.Join(tempDir, "sshagent.pid")
	sockFile := filepath.Join(tempDir, "sshauth.sock")

	content, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to read PID file -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	}
	sshAgentPid := string(content)

	showInfo("Killing ssh-agent process", program.indentLevel)
	err = exec.Command("kill", sshAgentPid).Run()
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to kill process -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		response := functionResponse{
			exitCode:    0,
			message:     "Finished",
			logLevel:    "success",
			indentLevel: program.indentLevel + 1,
		}
		handleFunctionResponse(response, true)
	}

	space()

	showInfo("Removing temporary files", program.indentLevel)
	for _, file := range []string{pidFile, sockFile} {
		err = os.Remove(file)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to remove file %v -> %v", file, err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 1,
			}
		}
	}

	return functionResponse{
		exitCode:    0,
		message:     "Finished",
		logLevel:    "success",
		indentLevel: program.indentLevel + 1,
	}
}

func utilsSSHAgentGetPID(tempDir string) {
	pidFile := filepath.Join(tempDir, "sshagent.pid")

	content, _ := ioutil.ReadFile(pidFile)
	sshAgentPid := string(content)

	fmt.Println(sshAgentPid)
}

func utilsSSHAgentGetSock(tempDir string) {
	sockFile := filepath.Join(tempDir, "sshauth.sock")

	content, _ := ioutil.ReadFile(sockFile)
	sshAgentSock := string(content)

	fmt.Println(sshAgentSock)
}
