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
	copy "github.com/otiai10/copy"
)

//
//// TARGETS
//

type Target struct {
	repo         Repo
	name         string
	path         string
	hooksDir     string
	tempDir      string
	poolDir      string
	disabledPath string
	environment  map[string]string
}

func getSelectedTargetsFromCLI(repoName string, targetNames []string, allTargets bool, interactiveSelection bool, multiple bool, program Program) (Repo, []Target, functionResponse) {
	var selectedRepo Repo
	var selectedTargets []Target

	if repoName == "" && interactiveSelection == false {
		return Repo{}, []Target{}, functionResponse{
			exitCode:    1,
			logLevel:    "error",
			message:     fmt.Sprintf("Flag '--repo/-r' or '--interactive/-i' should be specified"),
			indentLevel: program.indentLevel,
		}
	}

	if interactiveSelection == true {
		if allTargets != false || targetNames != nil || repoName != "" {
			return Repo{}, []Target{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Flag '--interactive/-i' cannot be used with flags '--repo/-r' and '--target/-t' or '--all/-a'"),
				indentLevel: program.indentLevel,
			}
		}

		availableRepos, response := getRepos(program)
		handleFunctionResponse(response, true)

		availableReposStrings := make([]string, len(availableRepos))
		for i, element := range availableRepos {
			availableReposStrings[i] = element.name
		}

		var selectedRepoIndex int
		promptRepo := &survey.Select{
			Message: "Select the repo",
			Options: availableReposStrings,
			Default: selectedRepoIndex,
		}
		err := survey.AskOne(promptRepo, &selectedRepoIndex, survey.WithPageSize(10))
		if err != nil {
			if err.Error() == "interrupt" {
				return Repo{}, []Target{}, functionResponse{
					exitCode:    1,
					message:     "Operation cancelled by user",
					logLevel:    "error",
					indentLevel: program.indentLevel,
				}
			}
		}

		selectedRepo = availableRepos[selectedRepoIndex]

		availableTargets, response := getRepoTargets(selectedRepo, program)
		handleFunctionResponse(response, true)

		availableTargetsStrings := make([]string, len(availableTargets))
		for i, element := range availableTargets {
			availableTargetsStrings[i] = element.name
		}

		var selectedTargetsIndices []int
		promptTarget := &survey.MultiSelect{
			Message: "Select the target(s)",
			Options: availableTargetsStrings,
			Default: selectedTargetsIndices,
		}
		err = survey.AskOne(promptTarget, &selectedTargetsIndices, survey.WithPageSize(10))
		if err != nil {
			if err.Error() == "interrupt" {
				return Repo{}, []Target{}, functionResponse{
					exitCode:    1,
					message:     "Operation cancelled by user",
					logLevel:    "error",
					indentLevel: program.indentLevel,
				}
			}
		}

		selectedTargets := make([]Target, len(selectedTargetsIndices))
		for i, index := range selectedTargetsIndices {
			selectedTargets[i] = availableTargets[index]
		}

		if len(selectedTargets) == 0 {
			return selectedRepo, selectedTargets, functionResponse{
				exitCode:    1,
				logLevel:    "attention",
				message:     fmt.Sprintf("No targets were selected"),
				indentLevel: program.indentLevel,
			}
		}

		return selectedRepo, selectedTargets, functionResponse{exitCode: 0}
	} else {
		if allTargets != false && targetNames != nil {
			return Repo{}, []Target{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Conflicting flags: both '--target/-t' and '--all/-a' flags cannot be specified at the same time"),
				indentLevel: program.indentLevel,
			}
		}

		if allTargets == false && targetNames == nil {
			return Repo{}, []Target{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Missing required flag: '--interactive/-i' or '--target/-t' or '--all/-a' flag must be specified"),
				indentLevel: program.indentLevel,
			}
		}

		repo := generateRepoObj(repoName, program)

		response := verifyRepoDirectory(repo, program)
		if response.exitCode != 0 {
			return Repo{}, []Target{}, functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Repo '%s' not found", repo.name),
				logLevel:    "error",
				indentLevel: program.indentLevel,
			}
		}

		if allTargets == true {
			selectedTargets, response := getRepoTargets(repo, program)
			handleFunctionResponse(response, true)

			return repo, selectedTargets, functionResponse{exitCode: 0}
		} else {
			selectedTargets = make([]Target, len(targetNames))
			for i, element := range targetNames {
				selectedTargets[i] = generateTargetObj(repo.name, element, program)
			}

			// Verify targets
			response := verifyTargetsDirectories(selectedTargets, program)
			handleFunctionResponse(response, true)

			return repo, selectedTargets, functionResponse{exitCode: 0}
		}
	}
}

func displayTargetTag(msg string, target Target) string {
	return fmt.Sprintf(msg) + fmt.Sprintf(" (") + salmonPink.Sprintf(target.repo.name) + fmt.Sprintf("/") + green.Sprintf(target.name) + fmt.Sprintf(")")
}

func getRepoTargets(repo Repo, program Program) ([]Target, functionResponse) {
	targetNames, err := ioutil.ReadDir(repo.targetsDir)
	if err != nil {
		return []Target{}, functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to read the repo's targets directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}
	targetNames = filterHiddenFilesAndDirectories(targetNames)

	// Generate an array of targets
	targets := make([]Target, len(targetNames))
	for i, element := range targetNames {
		targets[i] = generateTargetObj(repo.name, element.Name(), program)
	}

	if len(targets) == 0 {
		return targets, functionResponse{
			exitCode:    1,
			logLevel:    "attention",
			message:     "No targets found",
			indentLevel: program.indentLevel,
		}
	}

	return targets, functionResponse{
		exitCode: 0,
	}
}

func getEnabledTargets(repo Repo, program Program) []Target {
	availableTargets, response := getRepoTargets(repo, program)
	handleFunctionResponse(response, true)

	var enabledTargets []Target
	for _, target := range availableTargets {
		targetDisabled, response := isTargetDisabled(target, program)
		handleFunctionResponse(response, true)

		if targetDisabled == false {
			enabledTargets = append(enabledTargets, target)
		}
	}

	return enabledTargets
}

func generateTargetObj(repo string, target string, program Program) Target {
	defaultTargetEnv := map[string]string{
		"PROGRAM_NAME":          program.name,
		"DEFAULT_SHELL":         program.defaultShell,
		"PACPILOT_EXEC":         program.exec,
		"PACPILOT_UTILS":        fmt.Sprintf("%v utils", program.exec),
		"DATA_DIR":              program.dataDir,
		"REPOS_DIR":             program.reposDir,
		"TEMPLATES_DIR":         program.templatesDir,
		"REPOS_TEMPLATES_DIR":   program.reposTemplatesDir,
		"TARGETS_TEMPLATES_DIR": program.targetsTemplatesDir,
		"REPO_NAME":             repo,
		"REPO_DIR":              program.reposDir + "/" + repo,
		"REPO_HOOKS_DIR":        program.reposDir + "/" + repo + "/hooks",
		"REPO_TARGETS_DIR":      program.reposDir + "/" + repo + "/targets",
		"REPO_TEMP_DIR":         program.reposDir + "/" + repo + "/.tmp",
		"TARGET_NAME":           target,
		"TARGET_DIR":            program.reposDir + "/" + repo + "/targets" + "/" + target,
		"TARGET_HOOKS_DIR":      program.reposDir + "/" + repo + "/targets" + "/" + target + "/hooks",
		"TARGET_POOL_DIR":       program.reposDir + "/" + repo + "/targets" + "/" + target + "/pool",
		"TARGET_TEMP_DIR":       program.reposDir + "/" + repo + "/targets" + "/" + target + "/.tmp",
	}

	return Target{
		repo:         generateRepoObj(repo, program),
		name:         target,
		path:         program.reposDir + "/" + repo + "/targets" + "/" + target,
		hooksDir:     program.reposDir + "/" + repo + "/targets" + "/" + target + "/hooks",
		poolDir:      program.reposDir + "/" + repo + "/targets" + "/" + target + "/pool",
		tempDir:      program.reposDir + "/" + repo + "/targets" + "/" + target + "/.tmp",
		disabledPath: program.reposDir + "/" + repo + "/targets" + "/" + target + "/disabled",
		environment:  defaultTargetEnv,
	}
}

func isTargetDisabled(target Target, program Program) (bool, functionResponse) {
	if _, err := os.Stat(target.disabledPath); os.IsNotExist(err) {
		return false, functionResponse{
			exitCode:    0,
			indentLevel: program.indentLevel + 1,
		}
	} else if err != nil {
		return false, functionResponse{
			exitCode:    1,
			message:     "Failed to verify if target is disabled -> " + err.Error(),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		return true, functionResponse{
			exitCode: 0,
		}
	}
}

func enableTarget(target Target, program Program) functionResponse {
	if _, err := os.Stat(target.disabledPath); os.IsNotExist(err) {
		return functionResponse{
			exitCode:    0,
			message:     "Target already enabled",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	} else if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     "Failed to verify if target is disabled -> " + err.Error(),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		err = os.Remove(target.disabledPath)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     "Failed to remove disabled lock file -> " + err.Error(),
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
}

func targetsEnable(repo Repo, targets []Target, program Program) functionResponse {
	for _, target := range targets {
		space()
		showInfoSectionTitle(displayTargetTag("Enabling", target), program.indentLevel)

		response := enableTarget(target, program)
		if response.exitCode != 0 {
			return response
		}
		handleFunctionResponse(response, false)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func disableTarget(target Target, program Program) functionResponse {
	if _, err := os.Stat(target.disabledPath); os.IsNotExist(err) {
		file, err := os.Create(target.disabledPath)
		defer file.Close()
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     "Failed to create disabled lock file -> " + err.Error(),
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
	} else if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     "Failed to verify if target is enabled -> " + err.Error(),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		return functionResponse{
			exitCode:    0,
			message:     "Target already disabled",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	}
}

func targetsDisable(repo Repo, targets []Target, program Program) functionResponse {
	for _, target := range targets {
		space()
		showInfoSectionTitle(displayTargetTag("Disabling", target), program.indentLevel)

		response := disableTarget(target, program)
		if response.exitCode != 0 {
			return response
		}
		handleFunctionResponse(response, false)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func verifyTargetsDirectories(targets []Target, program Program) functionResponse {
	var failedTargets []Target

	for _, target := range targets {
		response := verifyTargetDirectory(target, program)
		if response.exitCode != 0 {
			failedTargets = append(failedTargets, target)
		}
	}

	var response functionResponse
	if len(failedTargets) > 0 {
		message := "The following target(s) was/were not found:\n"
		for _, target := range failedTargets {
			message = message + fmt.Sprintf("\n    - %s", target.name)
		}
		response = functionResponse{
			exitCode:    1,
			message:     message,
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	} else {
		response = functionResponse{
			exitCode: 0,
		}
	}

	return response
}

func verifyTargetDirectory(target Target, program Program) functionResponse {
	if _, err := os.Stat(target.path); os.IsNotExist(err) {
		return functionResponse{
			exitCode: 1,
		}
	}
	return functionResponse{
		exitCode: 0,
	}
}

func removeTargetTempDirectory(target Target, notRemoveTempDir bool, program Program) {
	var response functionResponse

	showInfoSectionTitle(fmt.Sprintf("Removing temporary directory"), program.indentLevel)

	if notRemoveTempDir == true {
		response = functionResponse{
			exitCode:    0,
			message:     "Skipping",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
		handleFunctionResponse(response, false)
		return
	}

	if _, err := os.Stat(target.tempDir); os.IsNotExist(err) {
		response = functionResponse{
			exitCode:    0,
			logLevel:    "attention",
			message:     fmt.Sprintf("Temporary directory not found"),
			indentLevel: program.indentLevel + 1,
		}
	} else {
		err = os.RemoveAll(target.tempDir)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				message:     "Failed to remove temporary directory -> " + err.Error(),
				logLevel:    "error",
				indentLevel: program.indentLevel + 1,
			}
		} else {
			response = functionResponse{
				exitCode:    0,
				message:     "Finished",
				logLevel:    "success",
				indentLevel: program.indentLevel + 1,
			}
		}
	}

	handleFunctionResponse(response, true)
}

func setupTargetTempDirectory(target Target, notCreateTempDir bool, program Program) {
	var response functionResponse

	showInfoSectionTitle(lightGray.Sprintf("Setting up temporary directory"), program.indentLevel)

	if notCreateTempDir == true {
		response = functionResponse{
			exitCode:    0,
			message:     "Skipping",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
		handleFunctionResponse(response, false)
		return
	}

	if _, err := os.Stat(target.tempDir); err == nil {
		response = functionResponse{
			exitCode:    0,
			message:     "Temporary directory already exists. Recreating it...",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 2,
		}
		handleFunctionResponse(response, false)

		err = os.RemoveAll(target.tempDir)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to recreate temporary directory -> '%v'", err.Error()),
				indentLevel: program.indentLevel + 3,
			}

			handleFunctionResponse(response, true)
		}
	}

	perm := os.FileMode(0755)
	err := os.Mkdir(target.tempDir, perm)
	if err != nil {
		response = functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to create temporary directory -> '%v'", err.Error()),
			indentLevel: program.indentLevel + 1,
		}
	} else {
		response = functionResponse{
			exitCode:    0,
			logLevel:    "success",
			message:     "Finished",
			indentLevel: program.indentLevel + 1,
		}
	}

	handleFunctionResponse(response, true)
}

func targetsCreate(program Program) functionResponse {
	var selectedRepo Repo
	var targetName string
	var targetTemplate string

	// Ask for a parent repo

	availableRepos, response := getRepos(program)
	handleFunctionResponse(response, true)

	availableReposStrings := make([]string, len(availableRepos))
	for i, element := range availableRepos {
		availableReposStrings[i] = element.name
	}

	var selectedIndex int
	promptRepo := &survey.Select{
		Message: "Select a repo",
		Options: availableReposStrings,
		Default: selectedIndex,
	}
	err := survey.AskOne(promptRepo, &selectedIndex, survey.WithPageSize(10))
	if err != nil {
		if err.Error() == "interrupt" {
			return functionResponse{
				exitCode:    1,
				message:     "Operation cancelled by user",
				logLevel:    "error",
				indentLevel: program.indentLevel,
			}
		}
	}

	selectedRepo = availableRepos[selectedIndex]

	//
	////
	//

	// Ask for a target name

	promptTargetName := &survey.Input{
		Message: "Target name:",
	}
	err = survey.AskOne(promptTargetName, &targetName, survey.WithValidator(survey.MinLength(2)))
	if err != nil {
		if err.Error() == "interrupt" {
			return functionResponse{
				exitCode:    1,
				message:     "Operation cancelled by user",
				logLevel:    "error",
				indentLevel: program.indentLevel,
			}
		}
	}

	// Generate a target object
	target := generateTargetObj(selectedRepo.name, targetName, program)

	// Ask for a target template
	availableTemplates, err := ioutil.ReadDir(program.targetsTemplatesDir)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to read the user's target templates directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}
	availableTemplates = filterHiddenFilesAndDirectories(availableTemplates)

	availableTemplatesStrings := make([]string, len(availableTemplates)+1)
	// Add a 'scratch' (empty) pseudo-template
	availableTemplatesStrings[0] = "scratch"
	// Add the available templates
	for i, element := range availableTemplates {
		availableTemplatesStrings[i+1] = element.Name()
	}

	promptTargetTemplate := &survey.Select{
		Message: "Target template:",
		Options: availableTemplatesStrings,
	}
	err = survey.AskOne(promptTargetTemplate, &targetTemplate)
	if err != nil {
		if err.Error() == "interrupt" {
			return functionResponse{
				exitCode:    1,
				message:     "Operation cancelled by user",
				logLevel:    "error",
				indentLevel: program.indentLevel,
			}
		}
	}

	// Verify if the scratch template was selected
	var targetTemplateDir string
	scratchTemplate := false
	if targetTemplate == "scratch" {
		scratchTemplate = true
	} else {
		targetTemplateDir = program.targetsTemplatesDir + "/" + targetTemplate
	}

	showInfoSectionTitle(displayTargetTag("Creating target", target), program.indentLevel)

	// Verify if target already exists
	response = verifyTargetDirectory(target, program)
	if response.exitCode == 0 {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Target '%s' already exists", target.name),
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	}

	// Create target directory
	space()
	showInfoSectionTitle("Creating target directory", program.indentLevel+1)
	err = os.Mkdir(target.path, 0755)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to create target directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel + 2,
		}
	}

	response = functionResponse{
		exitCode:    0,
		message:     "Finished",
		logLevel:    "success",
		indentLevel: program.indentLevel + 2,
	}
	handleFunctionResponse(response, false)

	// Create target's pool directory
	space()
	showInfoSectionTitle("Creating target's pool directory", program.indentLevel+1)
	err = os.Mkdir(target.poolDir, 0755)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to create target's pool directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel + 2,
		}
	}

	response = functionResponse{
		exitCode:    0,
		message:     "Finished",
		logLevel:    "success",
		indentLevel: program.indentLevel + 2,
	}
	handleFunctionResponse(response, false)

	// Copy template to target directory
	if scratchTemplate == false {
		space()
		showInfoSectionTitle("Copying template to target directory", program.indentLevel+1)
		copyOptions := copy.Options{
			PreserveTimes: true,
			PreserveOwner: true,
		}

		err = copy.Copy(targetTemplateDir, target.path, copyOptions)
		if err != nil {
			// Remove target directory
			_ = os.RemoveAll(target.path)

			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to copy template -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 2,
			}
		} else {
			response := functionResponse{
				exitCode:    0,
				message:     "Finished",
				logLevel:    "success",
				indentLevel: program.indentLevel + 2,
			}
			handleFunctionResponse(response, false)
		}
	}

	// Run post_create hook (if any)
	space()
	showInfoSectionTitle(lightGray.Sprintf("Running ")+orange.Sprintf("post_create")+lightGray.Sprintf(" hook"), program.indentLevel+1)
	if _, err := os.Stat(target.hooksDir + "/post_create"); err == nil {
		_, response := runHook(target.hooksDir+"/post_create", target.environment, true, true, true, true, true, incrementProgramIndentLevel(program, 1))

		if response.exitCode != 0 {
			handleFunctionResponse(response, false)

			space()

			// Remove target directory
			showInfoSectionTitle(lightGray.Sprintf("Removing target directory"), program.indentLevel+1)
			err = os.RemoveAll(target.path)

			if err != nil {
				response = functionResponse{
					exitCode:    response.exitCode,
					message:     "Failed to remove target directory",
					logLevel:    "attention",
					indentLevel: program.indentLevel + 2,
				}
			} else {
				response = functionResponse{
					exitCode:    response.exitCode,
					message:     "Removed",
					logLevel:    "attention",
					indentLevel: program.indentLevel + 2,
				}
			}

			return response
		} else {
			response := functionResponse{
				exitCode:    0,
				message:     "Finished",
				logLevel:    "success",
				indentLevel: program.indentLevel + 2,
			}
			handleFunctionResponse(response, false)
		}
	} else {
		response := functionResponse{
			exitCode:    0,
			message:     "Hook not found",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 2,
		}
		handleFunctionResponse(response, false)
	}

	space()

	return functionResponse{
		exitCode:    0,
		message:     "Finished",
		logLevel:    "success",
		indentLevel: program.indentLevel + 1,
	}
}

func targetsRm(repo Repo, targets []Target, program Program) functionResponse {
	var response functionResponse

	for index, target := range targets {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(targets)))
		showInfoSectionTitle(displayTargetTag("Removing", target), program.indentLevel)

		showInfoSectionTitle(lightGray.Sprintf("Running ")+orange.Sprintf("pre_rm")+lightGray.Sprintf(" hook"), program.indentLevel+1)

		if _, err := os.Stat(target.hooksDir + "/pre_rm"); os.IsNotExist(err) {
			response := functionResponse{
				exitCode:    0,
				message:     "Hook not found",
				logLevel:    "attention",
				indentLevel: program.indentLevel + 2,
			}
			handleFunctionResponse(response, false)
		} else {
			program = incrementProgramIndentLevel(program, 1)

			_, response := runHook(target.hooksDir+"/pre_rm", target.environment, true, true, true, true, true, program)
			response.indentLevel = program.indentLevel + 2
			handleFunctionResponse(response, true)
		}

		space()

		err := os.RemoveAll(target.path)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     "Failed to remove target -> " + err.Error(),
				indentLevel: program.indentLevel + 1,
			}
		} else {
			response = functionResponse{
				exitCode:    0,
				logLevel:    "success",
				message:     "Finished",
				indentLevel: program.indentLevel + 1,
			}
		}

		handleFunctionResponse(response, true)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func targetsLs(repos []Repo, program Program) functionResponse {
	for index, repo := range repos {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(repos)))
		showInfoSectionTitle(displayRepoTag("Listing targets", repo), program.indentLevel)

		targets, response := getRepoTargets(repo, program)
		if response.exitCode != 0 {
			response.indentLevel = program.indentLevel + 1
			return response
		}

		space()

		for _, target := range targets {
			// Get optional description (if hook exists)
			targetDescription, response := runHook(target.hooksDir+"/ls", target.environment, false, false, false, false, false, program)
			targetDescriptionString := targetDescription.Output

			var description string

			if response.exitCode == 0 {
				if len(targetDescriptionString) > 0 {
					description = fmt.Sprintf("(%s) ", blue.Sprintf(targetDescriptionString))
				} else {
					description = ""
				}
			}

			isTargetDisabled, response := isTargetDisabled(target, program)
			handleFunctionResponse(response, true)

			if isTargetDisabled == true {
				description += fmt.Sprintf("[%s]", red.Sprintf("disabled"))
			}

			showText(fmt.Sprintf(" - %s %s", target.name, description), program.indentLevel+1)
		}
	}

	return functionResponse{
		exitCode: 0,
	}
}

func targetsUpdate(repo Repo, targets []Target, repoPreHooks []string, repoPostHooks []string, program Program) functionResponse {
	hook := "update"
	response := targetsRunHooks(repo, targets, []string{hook}, repoPreHooks, repoPostHooks, false, false, false, false, false, program)

	return response
}

func targetsRunHooks(repo Repo, targets []Target, hooks []string, repoPreHooks []string, repoPostHooks []string, notCreateTempDir bool, notRemoveTempDir bool, notPrintOutput bool, notPrintEntryCmd bool, notPrintAlerts bool, program Program) functionResponse {
	var response functionResponse

	isRepoDisabled, response := isRepoDisabled(repo, program)
	if response.exitCode != 0 {
		response.indentLevel = program.indentLevel + 1

		return response
	}

	if isRepoDisabled == true {
		space()

		response = functionResponse{
			exitCode:    0,
			message:     "Repo is disabled",
			logLevel:    "attention",
			indentLevel: program.indentLevel,
		}
		return response
	}

	setupRepoTempDirectory(repo, true, notCreateTempDir, program)

	space()
	space()

	for _, hook := range repoPreHooks {
		showInfoSectionTitle(displayRepoTag(lightGray.Sprintf("Running ")+orange.Sprintf(hook)+lightGray.Sprintf(" hook"), repo), program.indentLevel)

		if _, err := os.Stat(repo.hooksDir + "/" + hook); os.IsNotExist(err) {
			response := functionResponse{
				exitCode:    0,
				message:     "Hook not found",
				logLevel:    "attention",
				indentLevel: program.indentLevel + 1,
			}
			handleFunctionResponse(response, true)
		} else {
			_, response := runHook(repo.hooksDir+"/"+hook, repo.environment, !notPrintOutput, true, true, !notPrintEntryCmd, true, program)
			response.indentLevel = program.indentLevel + 1

			handleFunctionResponse(response, false)

			if response.exitCode != 0 {
				space()
				removeRepoTempDirectory(repo, true, false, program)

				space()

				finishProgram(response.exitCode)
			}
		}

		space()
		space()
	}

	for index, target := range targets {
		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(targets)))

		showInfoSectionTitle(displayTargetTag("Running hook(s)", target), program.indentLevel)

		targetDisabled, response := isTargetDisabled(target, program)
		if response.exitCode != 0 {
			response.indentLevel = program.indentLevel + 1

			return response
		}

		if targetDisabled == true {
			response := functionResponse{
				exitCode:    0,
				message:     "Target is disabled",
				logLevel:    "attention",
				indentLevel: program.indentLevel + 1,
			}
			handleFunctionResponse(response, false)

			space()
			space()

			continue
		}

		program = incrementProgramIndentLevel(program, 1)

		space()

		setupTargetTempDirectory(target, notCreateTempDir, program)

		response = func(repo Repo, target Target, hooks []string, notRemoveTempDir bool, notPrintOutput bool, notPrintEntryCmd bool, program Program) functionResponse {
			for _, hook := range hooks {
				space()
				space()

				showInfoSectionTitle(lightGray.Sprintf("Running ")+orange.Sprintf(hook)+lightGray.Sprintf(" hook"), program.indentLevel)

				// Run hook
				if _, err := os.Stat(target.hooksDir + "/" + hook); os.IsNotExist(err) {
					response = functionResponse{
						exitCode:    1,
						message:     fmt.Sprintf("No '%v' hook found", hook),
						logLevel:    "error",
						indentLevel: program.indentLevel + 1,
					}
					handleFunctionResponse(response, false)
				} else {
					_, hookResponse := runHook(target.hooksDir+"/"+hook, target.environment, !notPrintOutput, true, true, !notPrintEntryCmd, true, program)

					if hookResponse.exitCode != 0 {
						hookResponse.indentLevel = program.indentLevel + 1
						return hookResponse
					} else {
						response = functionResponse{
							exitCode:    0,
							message:     "Finished",
							logLevel:    "success",
							indentLevel: program.indentLevel + 1,
						}
						handleFunctionResponse(response, false)
					}
				}
			}

			return functionResponse{
				exitCode: response.exitCode,
			}
		}(repo, target, hooks, notRemoveTempDir, notPrintOutput, notPrintEntryCmd, program)

		if response.exitCode != 0 {
			handleFunctionResponse(response, false)

			space()
			space()

			removeTargetTempDirectory(target, notRemoveTempDir, program)

			space()

			removeRepoTempDirectory(repo, true, notRemoveTempDir, decrementProgramIndentLevel(program, 1))

			space()
			finishProgram(response.exitCode)
		} else {
			handleFunctionResponse(response, false)
		}

		space()
		space()

		removeTargetTempDirectory(target, notRemoveTempDir, program)

		program = decrementProgramIndentLevel(program, 1)

		space()
		space()
	}

	for _, hook := range repoPostHooks {
		showInfoSectionTitle(displayRepoTag(lightGray.Sprintf("Running ")+orange.Sprintf(hook)+lightGray.Sprintf(" hook"), repo), program.indentLevel)

		if _, err := os.Stat(repo.hooksDir + "/" + hook); os.IsNotExist(err) {
			response = functionResponse{
				exitCode:    1,
				message:     "Hook not found",
				logLevel:    "error",
				indentLevel: program.indentLevel + 1,
			}
			handleFunctionResponse(response, true)
		} else {
			_, response := runHook(repo.hooksDir+"/"+hook, repo.environment, !notPrintOutput, true, true, !notPrintEntryCmd, true, program)
			response.indentLevel = program.indentLevel + 1

			handleFunctionResponse(response, false)

			if response.exitCode != 0 {
				space()
				removeRepoTempDirectory(repo, true, false, program)

				space()

				finishProgram(response.exitCode)
			}
		}

		space()
	}

	removeRepoTempDirectory(repo, true, notRemoveTempDir, program)

	return functionResponse{
		exitCode: 0,
	}
}

func targetsHooksLs(repo Repo, targets []Target, program Program) functionResponse {
	for index, target := range targets {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(targets)))
		showInfoSectionTitle(displayTargetTag("Listing hooks", target), program.indentLevel)

		hooks, err := ioutil.ReadDir(target.hooksDir)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Error reading the hooks directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 1,
			}
		}

		// Filter out .entry files
		filteredHooks := make([]os.FileInfo, 0)
		for _, element := range hooks {
			if !strings.HasSuffix(element.Name(), ".entry") {
				filteredHooks = append(filteredHooks, element)
			}
		}
		filteredHooks = filterHiddenFilesAndDirectories(filteredHooks)

		if len(filteredHooks) == 0 {
			showAttention("No hooks found", program.indentLevel)
		}

		for _, element := range filteredHooks {
			// Verify if hook has custom entry command
			var entryCommand string
			customEntryFilePath := target.hooksDir + "/" + element.Name() + ".entry"
			if _, err = os.Stat(customEntryFilePath); err == nil {
				contents, err := ioutil.ReadFile(customEntryFilePath)
				if err != nil {
					return functionResponse{
						exitCode:    1,
						message:     fmt.Sprintf("Failed to read custom entry configuration file for hook " + element.Name() + " -> " + err.Error()),
						logLevel:    "error",
						indentLevel: program.indentLevel,
					}
				}
				entryCommand = string(contents)
				entryCommand = strings.TrimLeft(entryCommand, "\n")
				entryCommand = strings.TrimRight(entryCommand, "\n")
			} else {
				entryCommand = program.defaultShell
			}
			showText(fmt.Sprintf("- %s (%s)", element.Name(), coral.Sprintf(entryCommand)), program.indentLevel+1)
		}
	}

	return functionResponse{
		exitCode: 0,
	}
}
