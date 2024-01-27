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
//// REPOS
//

type Repo struct {
	name         string
	path         string
	hooksDir     string
	targetsDir   string
	tempDir      string
	disabledPath string
	environment  map[string]string
}

func getSelectedReposFromCLI(repoNames []string, allRepos bool, interactiveSelection bool, multiple bool, program Program) ([]Repo, functionResponse) {
	var selectedRepos []Repo

	if interactiveSelection == true {
		if allRepos != false || repoNames != nil {
			return []Repo{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Flag '--interactive/-i' cannot be used together with flags '--repo/-r' or '--all/-a'"),
				indentLevel: program.indentLevel,
			}
		}

		availableRepos, response := getRepos(program)
		handleFunctionResponse(response, true)

		availableReposStrings := make([]string, len(availableRepos))
		for i, element := range availableRepos {
			availableReposStrings[i] = element.name
		}

		var selectedIndices []int
		prompt := &survey.MultiSelect{
			Message: "Select the repo(s)",
			Options: availableReposStrings,
			Default: selectedIndices,
		}
		err := survey.AskOne(prompt, &selectedIndices, survey.WithPageSize(10))
		if err != nil {
			if err.Error() == "interrupt" {
				return []Repo{}, functionResponse{
					exitCode:    1,
					message:     "Operation cancelled by user",
					logLevel:    "error",
					indentLevel: program.indentLevel,
				}
			}
		}

		selectedRepos := make([]Repo, len(selectedIndices))
		for i, index := range selectedIndices {
			selectedRepos[i] = availableRepos[index]
		}

		if len(selectedRepos) == 0 {
			return selectedRepos, functionResponse{
				exitCode:    1,
				logLevel:    "attention",
				message:     fmt.Sprintf("No repos were selected"),
				indentLevel: program.indentLevel,
			}
		}

		return selectedRepos, functionResponse{exitCode: 0}
	} else {
		if allRepos != false && repoNames != nil {
			return []Repo{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Conflicting flags: both '--repo/-r' and '--all/-a' flags cannot be specified at the same time"),
				indentLevel: program.indentLevel,
			}
		}

		if allRepos == false && repoNames == nil {
			return []Repo{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Missing required flag: '--interactive/--i' or '--repo/-r' or '--all/-a' flag must be specified"),
				indentLevel: program.indentLevel,
			}
		}

		if allRepos == true {
			selectedRepos, response := getRepos(program)
			handleFunctionResponse(response, true)

			return selectedRepos, functionResponse{exitCode: 0}
		} else {
			selectedRepos = make([]Repo, len(repoNames))
			for i, element := range repoNames {
				selectedRepos[i] = generateRepoObj(element, program)
			}

			// Verify repos
			response := verifyReposDirectories(selectedRepos, program)
			handleFunctionResponse(response, true)

			return selectedRepos, functionResponse{exitCode: 0}
		}
	}
}

func displayRepoTag(msg string, repo Repo) string {
	return fmt.Sprintf(msg) + fmt.Sprintf(" (") + salmonPink.Sprintf(repo.name) + fmt.Sprintf(")")
}

func getRepos(program Program) ([]Repo, functionResponse) {
	repoNames, err := ioutil.ReadDir(program.reposDir)
	if err != nil {
		return []Repo{}, functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Error reading the repos directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}
	repoNames = filterHiddenFilesAndDirectories(repoNames)

	// Generate a repos array
	repos := make([]Repo, len(repoNames))
	for i, element := range repoNames {
		repos[i] = generateRepoObj(element.Name(), program)
	}

	if len(repos) == 0 {
		return repos, functionResponse{
			exitCode:    1,
			logLevel:    "attention",
			message:     "No repos found",
			indentLevel: program.indentLevel,
		}
	}

	return repos, functionResponse{
		exitCode: 0,
	}
}

func getEnabledRepos(program Program) []Repo {
	availableRepos, response := getRepos(program)
	handleFunctionResponse(response, true)

	var enabledRepos []Repo
	for _, repo := range availableRepos {
		repoDisabled, response := isRepoDisabled(repo, program)
		handleFunctionResponse(response, true)

		if repoDisabled == false {
			enabledRepos = append(enabledRepos, repo)
		}
	}

	return enabledRepos
}

func generateRepoObj(repo string, program Program) Repo {
	defaultRepoEnv := map[string]string{
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
	}

	return Repo{
		name:         repo,
		path:         program.reposDir + "/" + repo,
		hooksDir:     program.reposDir + "/" + repo + "/hooks",
		targetsDir:   program.reposDir + "/" + repo + "/targets",
		tempDir:      program.reposDir + "/" + repo + "/.tmp",
		disabledPath: program.reposDir + "/" + repo + "/disabled",
		environment:  defaultRepoEnv,
	}
}

func isRepoDisabled(repo Repo, program Program) (bool, functionResponse) {
	if _, err := os.Stat(repo.disabledPath); os.IsNotExist(err) {
		return false, functionResponse{
			exitCode:    0,
			indentLevel: program.indentLevel + 1,
		}
	} else if err != nil {
		return false, functionResponse{
			exitCode:    1,
			message:     "Failed to verify if repo is disabled -> " + err.Error(),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		return true, functionResponse{
			exitCode: 0,
		}
	}
}

func enableRepo(repo Repo, program Program) functionResponse {
	if _, err := os.Stat(repo.disabledPath); os.IsNotExist(err) {
		return functionResponse{
			exitCode:    0,
			message:     "Repo already enabled",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	} else if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     "Failed to verify if repo is disabled -> " + err.Error(),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		err = os.Remove(repo.disabledPath)
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

func reposEnable(repos []Repo, program Program) functionResponse {
	for _, repo := range repos {
		space()
		showInfoSectionTitle(displayRepoTag("Enabling", repo), program.indentLevel)

		response := enableRepo(repo, program)
		if response.exitCode != 0 {
			return response
		}
		handleFunctionResponse(response, false)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func disableRepo(repo Repo, program Program) functionResponse {
	if _, err := os.Stat(repo.disabledPath); os.IsNotExist(err) {
		file, err := os.Create(repo.disabledPath)
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
			message:     "Failed to verify if repo is enabled -> " + err.Error(),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		return functionResponse{
			exitCode:    0,
			message:     "Repo already disabled",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	}
}

func reposDisable(repos []Repo, program Program) functionResponse {
	for _, repo := range repos {
		space()
		showInfoSectionTitle(displayRepoTag("Disabling", repo), program.indentLevel)

		response := disableRepo(repo, program)
		if response.exitCode != 0 {
			return response
		}
		handleFunctionResponse(response, false)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func verifyReposDirectories(repos []Repo, program Program) functionResponse {
	var failedRepos []Repo

	for _, repo := range repos {
		response := verifyRepoDirectory(repo, program)
		if response.exitCode != 0 {
			failedRepos = append(failedRepos, repo)
		}
	}

	var response functionResponse
	if len(failedRepos) > 0 {
		message := "The following repo(s) was/were not found:\n"
		for _, repo := range failedRepos {
			message = message + fmt.Sprintf("\n    - %s", repo.name)
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

func verifyRepoDirectory(repo Repo, program Program) functionResponse {
	if _, err := os.Stat(repo.path); os.IsNotExist(err) {
		return functionResponse{
			exitCode: 1,
		}
	}
	return functionResponse{
		exitCode: 0,
	}
}

func removeRepoTempDirectory(repo Repo, showRepoName bool, notRemoveTempDir bool, program Program) {
	var response functionResponse

	if showRepoName == true {
		showInfoSectionTitle(displayRepoTag("Removing temporary directory", repo), program.indentLevel)
	} else {
		showInfoSectionTitle("Removing temporary directory", program.indentLevel)
	}

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

	if _, err := os.Stat(repo.tempDir); os.IsNotExist(err) {
		response = functionResponse{
			exitCode:    0,
			logLevel:    "attention",
			message:     fmt.Sprintf("Temporary directory not found"),
			indentLevel: program.indentLevel + 1,
		}
	} else {
		err = os.RemoveAll(repo.tempDir)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Failed to remove temporary directory -> %v", err.Error()),
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

func setupRepoTempDirectory(repo Repo, showRepoName bool, notCreateTempDir bool, program Program) {
	var response functionResponse

	if showRepoName == true {
		showInfoSectionTitle(displayRepoTag("Setting up temporary directory", repo), program.indentLevel)
	} else {
		showInfoSectionTitle("Setting up temporary directory", program.indentLevel)
	}

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

	if _, err := os.Stat(repo.tempDir); err == nil {
		response = functionResponse{
			exitCode:    0,
			message:     "Temporary directory already exists. Recreating it...",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 2,
		}
		handleFunctionResponse(response, false)

		err = os.RemoveAll(repo.tempDir)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Failed to recreate temporary directory -> '%v'", err.Error()),
				indentLevel: program.indentLevel + 3,
			}

			handleFunctionResponse(response, true)
		}
	}

	perm := os.FileMode(0755)
	err := os.Mkdir(repo.tempDir, perm)
	if err != nil {
		response = functionResponse{
			exitCode:    1,
			logLevel:    "error",
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

func reposCreate(program Program) functionResponse {
	// Ask for a repo name
	var repoName string
	var repoTemplate string

	promptName := &survey.Input{
		Message: "Repo name:",
	}
	err := survey.AskOne(promptName, &repoName, survey.WithValidator(survey.MinLength(2)))
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

	// Generate a repo object
	repo := generateRepoObj(repoName, program)

	// Ask for a repo template
	availableTemplates, err := ioutil.ReadDir(program.reposTemplatesDir)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to read the user's repo templates directory -> " + err.Error()),
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

	promptTemplate := &survey.Select{
		Message: "Repo template:",
		Options: availableTemplatesStrings,
	}
	err = survey.AskOne(promptTemplate, &repoTemplate)
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
	var repoTemplateDir string
	scratchTemplate := false
	if repoTemplate == "scratch" {
		scratchTemplate = true
	} else {
		repoTemplateDir = program.reposTemplatesDir + "/" + repoTemplate
	}

	showInfoSectionTitle(displayRepoTag("Creating repo", repo), program.indentLevel)

	// Verify if repo already exists
	response := verifyRepoDirectory(repo, program)
	if response.exitCode == 0 {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Repo '%s' already exists", repo.name),
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	}

	// Create repo directory
	space()
	showInfoSectionTitle("Creating repo directory", program.indentLevel+1)
	err = os.Mkdir(repo.path, 0755)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to create repo directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel + 2,
		}
	}

	err = os.Mkdir(repo.targetsDir, 0755)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to create repo's targets directory -> " + err.Error()),
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

	// Copy template to repo directory
	if scratchTemplate == false {
		space()
		showInfoSectionTitle("Copying template to repo directory", program.indentLevel+1)
		copyOptions := copy.Options{
			PreserveTimes: true,
			PreserveOwner: true,
		}

		err = copy.Copy(repoTemplateDir, repo.path, copyOptions)
		if err != nil {
			// Remove repo directory
			_ = os.RemoveAll(repo.path)

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
	if _, err := os.Stat(repo.hooksDir + "/post_create"); err == nil {
		_, response := runHook(repo.hooksDir+"/post_create", repo.environment, true, true, true, true, true, incrementProgramIndentLevel(program, 1))

		if response.exitCode != 0 {
			handleFunctionResponse(response, false)

			space()

			// Remove repo directory
			showInfoSectionTitle(lightGray.Sprintf("Removing repo directory"), program.indentLevel+1)
			err = os.RemoveAll(repo.path)

			if err != nil {
				response = functionResponse{
					exitCode:    response.exitCode,
					message:     "Failed to remove repo directory",
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

func reposRm(repos []Repo, program Program) functionResponse {
	for index, repo := range repos {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(repos)))
		showInfoSectionTitle(displayRepoTag("Removing", repo), program.indentLevel)

		// Show number of targets and ask for confirmation
		targets, response := getRepoTargets(repo, program)
		if response.exitCode != 0 && response.logLevel != "attention" {
			handleFunctionResponse(response, true)
		}

		if len(targets) > 0 {
			showAttention(fmt.Sprintf("This action will delete %v targets from this repo", len(targets)), program.indentLevel+1)

			userConfirmation := askConfirmation("Enter 'yes/y' to confirm or 'no/n' to cancel the operation", incrementProgramIndentLevel(program, 2))
			if userConfirmation == false {
				return functionResponse{
					exitCode:    1,
					logLevel:    "error",
					message:     "Operation cancelled by user",
					indentLevel: program.indentLevel + 1,
				}
			}

			space()

			// Run the pre_rm hook for each target
			for _, target := range targets {
				showInfoSectionTitle(displayTargetTag(lightGray.Sprintf("Running ")+gray.Sprintf("pre_rm")+lightGray.Sprintf(" hook"), target), program.indentLevel+1)

				if _, err := os.Stat(target.hooksDir + "/pre_rm"); os.IsNotExist(err) {
					response = functionResponse{
						exitCode:    0,
						message:     "Hook not found",
						logLevel:    "attention",
						indentLevel: program.indentLevel + 2,
					}
					handleFunctionResponse(response, false)
				} else {
					program := incrementProgramIndentLevel(program, 1)

					_, response := runHook(target.hooksDir+"/pre_rm", target.environment, true, true, true, true, true, program)
					response.indentLevel = program.indentLevel + 2
					handleFunctionResponse(response, true)
				}

				space()
			}
		} else {
			response = functionResponse{
				exitCode:    0,
				message:     "Repo has no targets. Keeping on...",
				logLevel:    "attention",
				indentLevel: program.indentLevel + 1,
			}
			handleFunctionResponse(response, false)
		}

		// Run pre_rm hook for the repo
		space()

		showInfoSectionTitle(displayRepoTag(lightGray.Sprintf("Running ")+gray.Sprintf("pre_rm")+lightGray.Sprintf(" hook"), repo), program.indentLevel+1)

		if _, err := os.Stat(repo.hooksDir + "/pre_rm"); os.IsNotExist(err) {
			response = functionResponse{
				exitCode:    0,
				message:     "Hook not found",
				logLevel:    "attention",
				indentLevel: program.indentLevel + 2,
			}
			handleFunctionResponse(response, false)
		} else {
			program = incrementProgramIndentLevel(program, 1)

			_, response := runHook(repo.hooksDir+"/pre_rm", repo.environment, true, true, true, true, true, program)
			response.indentLevel = program.indentLevel + 2
			handleFunctionResponse(response, true)
		}

		space()

		err := os.RemoveAll(repo.path)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     "Failed to remove repo -> " + err.Error(),
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

func reposLs(program Program) functionResponse {
	space()
	showInfoSectionTitle("Listing repos", program.indentLevel)

	repos, response := getRepos(program)
	if response.exitCode != 0 {
		response.indentLevel = program.indentLevel + 1
		return response
	}

	space()

	for _, repo := range repos {
		// Get optional description (if hook exists)
		repoDescription, response := runHook(repo.hooksDir+"/ls", repo.environment, false, false, false, false, false, program)
		repoDescriptionString := repoDescription.Output

		var description string

		if response.exitCode == 0 {
			if len(repoDescriptionString) > 0 {
				description = fmt.Sprintf("(%s) ", blue.Sprintf(repoDescriptionString))
			} else {
				description = ""
			}
		}

		isRepoDisabled, response := isRepoDisabled(repo, program)
		handleFunctionResponse(response, true)

		if isRepoDisabled == true {
			description += fmt.Sprintf("[%s]", red.Sprintf("disabled"))
		}

		showText(fmt.Sprintf(" - %s %s", repo.name, description), program.indentLevel+1)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func reposRunHooks(repos []Repo, hooks []string, notCreateTempDir bool, notRemoveTempDir bool, notPrintOutput bool, notPrintEntryCmd bool, notPrintAlerts bool, program Program) functionResponse {
	for index, repo := range repos {
		space()
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(repos)))

		showInfoSectionTitle(displayRepoTag("Running hook(s)", repo), program.indentLevel)

		isRepoDisabled, response := isRepoDisabled(repo, program)
		if response.exitCode != 0 {
			response.indentLevel = program.indentLevel + 1

			return response
		}

		if isRepoDisabled == true {
			response = functionResponse{
				exitCode:    0,
				message:     "Repo is disabled",
				logLevel:    "attention",
				indentLevel: program.indentLevel + 1,
			}
			handleFunctionResponse(response, false)

			continue
		}

		program = incrementProgramIndentLevel(program, 1)

		setupRepoTempDirectory(repo, false, notCreateTempDir, program)

		response = func(repo Repo, hooks []string, notRemoveTempDir bool, notPrintOutput bool, notPrintEntryCmd bool, program Program) functionResponse {
			for _, hook := range hooks {
				space()
				space()

				showInfoSectionTitle(lightGray.Sprintf("Running ")+orange.Sprintf(hook)+lightGray.Sprintf(" hook"), program.indentLevel)

				// Run hook
				if _, err := os.Stat(repo.hooksDir + "/" + hook); os.IsNotExist(err) {
					response = functionResponse{
						exitCode:    1,
						message:     fmt.Sprintf("No '%v' hook found", hook),
						logLevel:    "error",
						indentLevel: program.indentLevel + 1,
					}
					handleFunctionResponse(response, false)
				} else {
					_, hookResponse := runHook(repo.hooksDir+"/"+hook, repo.environment, !notPrintOutput, true, true, !notPrintEntryCmd, true, program)

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
				exitCode: 0,
			}
		}(repo, hooks, notRemoveTempDir, notPrintOutput, notPrintEntryCmd, program)

		if response.exitCode != 0 {
			handleFunctionResponse(response, false)

			space()
			space()

			removeRepoTempDirectory(repo, false, notRemoveTempDir, program)

			space()
			finishProgram(response.exitCode)
		} else {
			handleFunctionResponse(response, false)
		}

		space()
		space()

		removeRepoTempDirectory(repo, false, notRemoveTempDir, program)

		program = decrementProgramIndentLevel(program, 1)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func reposHooksLs(repos []Repo, program Program) functionResponse {
	for index, repo := range repos {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(repos)))
		showInfoSectionTitle(displayRepoTag("Listing hooks", repo), program.indentLevel)

		hooks, err := ioutil.ReadDir(repo.hooksDir)
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
			customEntryFilePath := repo.hooksDir + "/" + element.Name() + ".entry"
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
