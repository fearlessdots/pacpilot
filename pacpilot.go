package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"errors"
	"fmt"
	"os"
	"strconv"

	// External modules
	cobra "github.com/spf13/cobra"
	cobra_doc "github.com/spf13/cobra/doc"
)

//
////
//

func main() {
	var program Program
	var dataDir string

	// Initialize program
	program = initializeDefaultProgram("")

	var rootCmd = &cobra.Command{
		Use:   fmt.Sprintf("%v [command]", program.name),
		Short: program.shortDescription,
		Long:  program.shortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			displayProgramInfo(program)

			space()

			showText(fmt.Sprintf("Run %v to get started.", blue.Sprintf("%v --help/-h", program.name)), program.indentLevel)

			finishProgram(0)
		},
	}
	rootCmd.PersistentFlags().StringVarP(&dataDir, "directory", "D", "", "Data directory")

	var showVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the program's version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(program.nameAscii)
			space()
			showInfoSectionTitle(fmt.Sprintf("Version: %s", green.Sprintf(program.version)), program.indentLevel)
		},
	}

	var userInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Create user data directory",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if dataDir == "" {
				return errors.New("A data directory should be specified using the '-D' flag")
			}

			showAttention(salmonPink.Sprintf("Running %v using data directory at: %v", program.name, dataDir), program.indentLevel)
			space()

			program = initializeDefaultProgram(dataDir)

			space()

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(program.nameAscii)

			space()

			showText(fmt.Sprintf("Creating user data directory at %v", green.Sprintf(dataDir)), program.indentLevel)

			if _, err := os.Stat(program.dataDir); os.IsNotExist(err) {
				verifyDataDirectory(false, decrementProgramIndentLevel(program, 1))
			} else {
				showAttention("> User data directory already exists", program.indentLevel+1)
			}
		},
	}

	var docsCmd = &cobra.Command{
		Use:   "docs",
		Short: "Program documentation",
	}

	//
	//// DOCUMENTATION
	//

	var docsOutDir string

	var docsGenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate program documentation (markdown files)",
		Run: func(cmd *cobra.Command, args []string) {
			var response functionResponse

			space()
			fmt.Println(program.nameAscii)
			space()

			showInfoSectionTitle("Generating program documentation", program.indentLevel)

			_, err := os.Stat(docsOutDir)
			if os.IsNotExist(err) {
				err := os.MkdirAll(docsOutDir, 0755)
				if err != nil {
					response = functionResponse{
						exitCode:    1,
						message:     fmt.Sprintf("Failed to create output directory -> " + err.Error()),
						logLevel:    "error",
						indentLevel: program.indentLevel + 1,
					}
				} else {
					response = functionResponse{
						exitCode:    0,
						message:     fmt.Sprintf("Created output directory at '%v'", docsOutDir),
						logLevel:    "success",
						indentLevel: program.indentLevel + 1,
					}
				}
			} else {
				if err != nil {
					response = functionResponse{
						exitCode:    1,
						message:     fmt.Sprintf("Failed to verify if output directory exists -> " + err.Error()),
						logLevel:    "error",
						indentLevel: program.indentLevel + 1,
					}
				} else {
					response = functionResponse{
						exitCode:    1,
						message:     fmt.Sprintf("Output directory '%v' already exists. Remove it before generating the documentation", docsOutDir),
						logLevel:    "attention",
						indentLevel: program.indentLevel + 1,
					}
				}
			}
			handleFunctionResponse(response, true)

			err = cobra_doc.GenMarkdownTree(rootCmd, "./docs")
			if err != nil {
				response = functionResponse{
					exitCode:    1,
					message:     fmt.Sprintf("Failed to build documentation -> " + err.Error()),
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

			handleFunctionResponse(response, true)
		},
	}

	docsGenerateCmd.Flags().StringVarP(&docsOutDir, "out", "o", "./docs", "Output directory")

	//
	//// UTILITIES
	//

	var utilitiesCmd = &cobra.Command{
		Use:   "utils",
		Short: "Utilites for hooks execution",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if dataDir == "" {
				return errors.New("A data directory should be specified using the '-D' flag")
			}

			showAttention(salmonPink.Sprintf("Running %v using data directory at: %v", program.name, dataDir), program.indentLevel)
			space()

			program = initializeDefaultProgram(dataDir)

			space()

			return nil
		},
	}

	var utilityMsgCmd = &cobra.Command{
		Use:   "msg <color_hex> <message>",
		Short: "Print a message in a specific color given a HEX code",
		Long: `The 'msg' command allows you to print messages in a color
		specified by a HEX code.

		Arguments:
		1. color_hex: This should be a valid HEX color code. The message will be
		printed in this color.
		2. message: This is the message that you want to print in the console. The
		message will be printed in the color specified by color_hex.`,
		Example: "utils msg '#FF5733' 'This is a colored message'",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			utilsMsg(args[0], args[1], program)
		},
	}

	var utilityAttentionCmd = &cobra.Command{
		Use:   "attention <message>",
		Short: "Display attention message",
		Long: `The 'attention' command allows you to display a message
		to attract user's attention.

		Argument:
		1. message: This is the attention-grabbing message you want to display.`,
		Example: "utils attention 'This is an attention message'",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utilsShowAttention(args[0], program)
		},
	}

	var utilityErrorCmd = &cobra.Command{
		Use:   "error <message>",
		Short: "Display error message",
		Long: `The 'error' command allows you to display an error
		message.

		Argument:
		1. message: This is the error message you want to display.`,
		Example: "utils error 'This is an error message'",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utilsShowError(args[0], program)
		},
	}

	var utilitySuccessCmd = &cobra.Command{
		Use:   "success <message>",
		Short: "Display success message",
		Long: `The 'success' command allows you to display a success
		message.

		Argument:
		1. message: This is the success message you want to display.`,
		Example: "utils success 'This is a success message'",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utilsShowSuccess(args[0], program)
		},
	}

	var utilitySectionCmd = &cobra.Command{
		Use:   "section <message>",
		Short: "Display section title",
		Long: `The 'section' command allows you to display a section
		title.

		Argument:
		1. message: This is the section title you want to display.`,
		Example: "utils section 'This is a section title'",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utilsShowSection(args[0], program)
		},
	}

	var utilityHrCmd = &cobra.Command{
		Use:   "hr <char> <factor>",
		Short: "Display a horizontal line",
		Long: `The 'hr' command allows you to display a horizontal line.

		Arguments:
		1. char: This character is used to construct the horizontal line.
		2. factor: This is the length factor of the line. The factor represents the
		percentage of the terminal's width.`,
		Example: "utils hr '-' 0.45",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			factor, _ := strconv.ParseFloat(args[1], 64)
			utilsHr(args[0], factor, program)
		},
	}

	var utilityConfirmCmd = &cobra.Command{
		Use:   "confirm <message>",
		Short: "Ask for confirmation",
		Long: `The 'confirm' command prompts the user for a confirmation
		based on a provided message.

		Argument:
		1. message: This is the message to display when asking for confirmation.

		The command will then wait for user input. If the user confirms (by typing
		'y' or 'yes'),
		the command will terminate the program with an exit status of 0.
		If the user rejects (by typing 'n' or 'no'), the command will terminate the
		program with an exit status of 1.`,
		Example: "utils confirm 'Are you sure you want to continue?'",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utilsAskConfirmation(args[0], program)
		},
	}

	var utilitySSHAgentStartCmd = &cobra.Command{
		Use:   "sshagent-start <privateKeyPath> <tempDir>",
		Short: "Starts an ssh-agent process and adds a private key",
		Long: `The 'sshagent-start' command starts an ssh-agent process
		and adds a private key to it.

		Arguments:
		1. privateKeyPath: The path to the private key that will be added to the
		ssh-agent.
		2. tempDir: The temporary directory in which the ssh-agent will create its
		socket and pid files. For example, it can be set with environment variables
		$CRATE_TEMP_DIR or $TARGET_TEMP_DIR.`,
		Example: "utils sshagent-start '/path/to/my/key' $TARGET_TEMP_DIR",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			response := utilsSSHAgentStart(args[0], args[1], program)
			handleFunctionResponse(response, true)
		},
	}

	var utilitySSHAgentStopCmd = &cobra.Command{
		Use:   "sshagent-stop <tempDir>",
		Short: "Stops the ssh-agent process",
		Long: `The 'sshagent-stop' command stops the ssh-agent process.

		Arguments:
		1. tempDir: The directory where the ssh-agent created its socket and pid
		files. For example, it can be set with environment variables $CRATE_TEMP_DIR or
		$TARGET_TEMP_DIR.`,
		Example: "utils sshagent-stop $TARGET_TEMP_DIR",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			response := utilsSSHAgentStop(args[0], program)
			handleFunctionResponse(response, true)
		},
	}

	var utilitySSHAgentGetPIDCmd = &cobra.Command{
		Use:   "sshagent-getpid <tempDir>",
		Short: "Get the process ID of the ssh-agent",
		Long: `The 'sshagent-getpid' command retrieves the process ID of
		the ssh-agent.

		Arguments:
		1. tempDir: The directory where the ssh-agent created its pid file. For
		example, it can be set with environment variables $CRATE_TEMP_DIR or
		$TARGET_TEMP_DIR.`,
		Example: "utils sshagent-getpid $TARGET_TEMP_DIR",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utilsSSHAgentGetPID(args[0])
		},
	}

	var utilitySSHAgentGetSockCmd = &cobra.Command{
		Use:   "sshagent-getsock <tempDir>",
		Short: "Get the socket path of the ssh-agent",
		Long: `The 'sshagent-getsock' command retrieves the socket path
		of the ssh-agent.

		Arguments:
		1. tempDir: The directory where the ssh-agent created its socket file. For
		example, it can be set with environment variables $CRATE_TEMP_DIR or
		$TARGET_TEMP_DIR.`,
		Example: "utils sshagent-getsock $TARGET_TEMP_DIR",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			utilsSSHAgentGetSock(args[0])
		},
	}

	//
	////
	//

	var interactiveSelection bool
	var notCreateTempDir bool
	var notRemoveTempDir bool
	var notPrintOutput bool

	var repoPreHooks []string
	var repoPostHooks []string

	//
	//// REPOS
	//

	var repoName string
	var repoNames []string
	var repoHooksNames []string
	var allRepos bool

	var reposCmd = &cobra.Command{
		Use:   "repos",
		Short: "Manage repos",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if dataDir == "" {
				return errors.New("A data directory should be specified using the '-D' flag")
			}

			showAttention(salmonPink.Sprintf("Running %v using data directory at: %v", program.name, dataDir), program.indentLevel)
			space()

			program = initializeDefaultProgram(dataDir)

			// Verify user data directory
			response := verifyDataDirectory(true, program)
			handleFunctionResponse(response, true)

			return nil
		},
	}

	var reposCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create repos",
		Run: func(cmd *cobra.Command, args []string) {
			response := reposCreate(program)
			handleFunctionResponse(response, true)
		},
	}

	var reposRmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove repos",
		Run: func(cmd *cobra.Command, args []string) {
			selectedRepos, response := getSelectedReposFromCLI(repoNames, allRepos, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = reposRm(selectedRepos, program)
			handleFunctionResponse(response, true)
		},
	}

	reposRmCmd.Flags().StringSliceVarP(&repoNames, "repo", "r", nil, "Repo(s) name(s)")
	reposRmCmd.Flags().BoolVarP(&allRepos, "all", "a", false, "Include all repos")
	reposRmCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	reposRmCmd.Flags().SetInterspersed(false)

	var reposLsCmd = &cobra.Command{
		Use:   "ls",
		Short: "List all repos",
		Run: func(cmd *cobra.Command, args []string) {
			response := reposLs(program)
			handleFunctionResponse(response, true)
		},
	}

	var serverPort string
	var debugMode bool
	var trustedProxies []string

	var reposServeCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve repos",
		Run: func(cmd *cobra.Command, args []string) {
			response := reposServe(serverPort, debugMode, trustedProxies, program)
			handleFunctionResponse(response, true)
		},
	}

	reposServeCmd.Flags().StringVarP(&serverPort, "port", "p", "8080", "Server port")
	reposServeCmd.Flags().StringSliceVarP(&trustedProxies, "proxy", "x", nil, "Trusted proxy (by default, all proxies are trusted)")
	reposServeCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "Debug mode")
	reposServeCmd.Flags().SetInterspersed(false)

	var reposEnableCmd = &cobra.Command{
		Use:   "enable",
		Short: "Enable repos",
		Run: func(cmd *cobra.Command, args []string) {
			selectedRepos, response := getSelectedReposFromCLI(repoNames, allRepos, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = reposEnable(selectedRepos, program)
			handleFunctionResponse(response, true)
		},
	}

	reposEnableCmd.Flags().StringSliceVarP(&repoNames, "repo", "r", nil, "Repo name")
	reposEnableCmd.Flags().BoolVarP(&allRepos, "all", "a", false, "Include all repos")
	reposEnableCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	reposEnableCmd.Flags().SetInterspersed(false)

	var reposDisableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disable repos",
		Run: func(cmd *cobra.Command, args []string) {
			selectedRepos, response := getSelectedReposFromCLI(repoNames, allRepos, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = reposDisable(selectedRepos, program)
			handleFunctionResponse(response, true)
		},
	}

	reposDisableCmd.Flags().StringSliceVarP(&repoNames, "repo", "r", nil, "Repo name")
	reposDisableCmd.Flags().BoolVarP(&allRepos, "all", "a", false, "Include all repos")
	reposDisableCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	reposDisableCmd.Flags().SetInterspersed(false)

	var reposHooksCmd = &cobra.Command{
		Use:   "hooks",
		Short: "Manage repo hooks",
	}

	var reposHooksLsCmd = &cobra.Command{
		Use:   "ls",
		Short: "List repo hooks",
		Run: func(cmd *cobra.Command, args []string) {
			selectedRepos, response := getSelectedReposFromCLI(repoNames, allRepos, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = reposHooksLs(selectedRepos, program)
			handleFunctionResponse(response, true)
		},
	}

	reposHooksLsCmd.Flags().StringSliceVarP(&repoNames, "repo", "r", nil, "Repo(s) name(s)")
	reposHooksLsCmd.Flags().BoolVarP(&allRepos, "all", "a", false, "Include all repos")
	reposHooksLsCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	reposHooksLsCmd.Flags().SetInterspersed(false)

	var reposHooksRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run repo hook(s)",
		Run: func(cmd *cobra.Command, args []string) {
			if len(repoHooksNames) == 0 {
				response := functionResponse{
					exitCode:    1,
					logLevel:    "error",
					message:     fmt.Sprintf("Flag '--hook/-k' should be specified"),
					indentLevel: program.indentLevel,
				}
				handleFunctionResponse(response, true)
			}

			selectedRepos, response := getSelectedReposFromCLI(repoNames, allRepos, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = reposRunHooks(selectedRepos, repoHooksNames, notCreateTempDir, notRemoveTempDir, notPrintOutput, false, true, program)
			handleFunctionResponse(response, true)
		},
	}

	reposHooksRunCmd.Flags().StringSliceVarP(&repoNames, "repo", "r", nil, "Repo(s) name(s)")
	reposHooksRunCmd.Flags().BoolVarP(&allRepos, "all", "a", false, "Include all repos")
	reposHooksRunCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	reposHooksRunCmd.Flags().StringSliceVarP(&repoHooksNames, "hook", "k", nil, "Hook(s) name(s)")
	reposHooksRunCmd.Flags().BoolVarP(&notCreateTempDir, "nocreatetemp", "", false, "Do not create the temporary directory before running the hook(s) (by default, it is created)")
	reposHooksRunCmd.Flags().BoolVarP(&notRemoveTempDir, "noremovetemp", "", false, "Do not remove the temporary directory after the hook(s) has/have finished running (by default, it is removed)")
	reposHooksRunCmd.Flags().BoolVarP(&notPrintOutput, "quiet", "q", false, "Do not print command output (silent)")
	reposHooksRunCmd.Flags().SetInterspersed(false)

	//
	//// TARGETS
	//

	var targetHooksNames []string
	var targetNames []string
	var allTargets bool

	var targetsCmd = &cobra.Command{
		Use:   "targets",
		Short: "Manage targets",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if dataDir == "" {
				return errors.New("A data directory should be specified using the '-D' flag")
			}

			showAttention(salmonPink.Sprintf("Running %v using data directory at: %v", program.name, dataDir), program.indentLevel)
			space()

			program = initializeDefaultProgram(dataDir)

			// Verify user data directory
			response := verifyDataDirectory(true, program)
			handleFunctionResponse(response, true)

			return nil
		},
	}

	var targetsLsCmd = &cobra.Command{
		Use:   "ls",
		Short: "List targets",
		Run: func(cmd *cobra.Command, args []string) {
			selectedRepos, response := getSelectedReposFromCLI(repoNames, allRepos, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = targetsLs(selectedRepos, program)
			handleFunctionResponse(response, true)
		},
	}

	targetsLsCmd.Flags().StringSliceVarP(&repoNames, "repo", "r", nil, "Repo(s) name(s)")
	targetsLsCmd.Flags().BoolVarP(&allRepos, "all", "a", false, "Include all repos")
	targetsLsCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	targetsLsCmd.Flags().SetInterspersed(false)

	var targetsUpdateCmd = &cobra.Command{ // TODO:
		Use:   "update",
		Short: "Update targets",
		Run: func(cmd *cobra.Command, args []string) {
			repo, selectedTargets, response := getSelectedTargetsFromCLI(repoName, targetNames, allTargets, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = targetsUpdate(repo, selectedTargets, repoPreHooks, repoPostHooks, program)
			handleFunctionResponse(response, true)
		},
	}

	targetsUpdateCmd.Flags().StringVarP(&repoName, "repo", "r", "", "Repo name")
	targetsUpdateCmd.Flags().StringSliceVarP(&targetNames, "target", "t", nil, "Target(s) name(s)")
	targetsUpdateCmd.Flags().BoolVarP(&allTargets, "all", "a", false, "Include all targets")
	targetsUpdateCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	targetsUpdateCmd.Flags().StringSliceVarP(&repoPreHooks, "repopre", "", nil, "Repo pre hook(s)")
	targetsUpdateCmd.Flags().StringSliceVarP(&repoPostHooks, "repopost", "", nil, "Repo post hook(s)")
	targetsUpdateCmd.Flags().SetInterspersed(false)

	var targetsEnableCmd = &cobra.Command{
		Use:   "enable",
		Short: "Enable targets",
		Run: func(cmd *cobra.Command, args []string) {
			repo, selectedTargets, response := getSelectedTargetsFromCLI(repoName, targetNames, allTargets, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = targetsEnable(repo, selectedTargets, program)
			handleFunctionResponse(response, true)
		},
	}

	targetsEnableCmd.Flags().StringVarP(&repoName, "repo", "r", "", "Repo name")
	targetsEnableCmd.Flags().StringSliceVarP(&targetNames, "target", "t", nil, "Target(s) name(s)")
	targetsEnableCmd.Flags().BoolVarP(&allTargets, "all", "a", false, "Include all targets")
	targetsEnableCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	targetsEnableCmd.Flags().SetInterspersed(false)

	var targetsDisableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disable targets",
		Run: func(cmd *cobra.Command, args []string) {
			repo, selectedTargets, response := getSelectedTargetsFromCLI(repoName, targetNames, allTargets, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = targetsDisable(repo, selectedTargets, program)
			handleFunctionResponse(response, true)
		},
	}

	targetsDisableCmd.Flags().StringVarP(&repoName, "repo", "r", "", "Repo name")
	targetsDisableCmd.Flags().StringSliceVarP(&targetNames, "target", "t", nil, "Target(s) name(s)")
	targetsDisableCmd.Flags().BoolVarP(&allTargets, "all", "a", false, "Include all targets")
	targetsDisableCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	targetsDisableCmd.Flags().SetInterspersed(false)

	var targetsCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create targets",
		Run: func(cmd *cobra.Command, args []string) {
			response := targetsCreate(program)
			handleFunctionResponse(response, true)
		},
	}

	var targetsRmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove targets",
		Run: func(cmd *cobra.Command, args []string) {
			repo, selectedTargets, response := getSelectedTargetsFromCLI(repoName, targetNames, allTargets, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = targetsRm(repo, selectedTargets, program)
			handleFunctionResponse(response, true)
		},
	}

	targetsRmCmd.Flags().StringVarP(&repoName, "repo", "r", "", "Repo name")
	targetsRmCmd.Flags().StringSliceVarP(&targetNames, "target", "t", nil, "Target(s) name(s)")
	targetsRmCmd.Flags().BoolVarP(&allTargets, "all", "a", false, "Include all targets")
	targetsRmCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	targetsRmCmd.Flags().SetInterspersed(false)

	var targetsHooksCmd = &cobra.Command{
		Use:   "hooks",
		Short: "Manage target hooks",
	}

	var targetsHooksLsCmd = &cobra.Command{
		Use:   "ls",
		Short: "List target hooks",
		Run: func(cmd *cobra.Command, args []string) {
			repo, selectedTargets, response := getSelectedTargetsFromCLI(repoName, targetNames, allTargets, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = targetsHooksLs(repo, selectedTargets, program)
			handleFunctionResponse(response, true)
		},
	}

	targetsHooksLsCmd.Flags().StringVarP(&repoName, "repo", "r", "", "Repo name")
	targetsHooksLsCmd.Flags().StringSliceVarP(&targetNames, "target", "t", nil, "Target(s) name(s)")
	targetsHooksLsCmd.Flags().BoolVarP(&allTargets, "all", "a", false, "Include all targets")
	targetsHooksLsCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	targetsHooksLsCmd.Flags().SetInterspersed(false)

	var targetsHooksRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run target hook(s)",
		Run: func(cmd *cobra.Command, args []string) {
			if len(targetHooksNames) == 0 {
				response := functionResponse{
					exitCode:    1,
					logLevel:    "error",
					message:     fmt.Sprintf("Flag '--hook/-k' should be specified"),
					indentLevel: program.indentLevel,
				}
				handleFunctionResponse(response, true)
			}

			repo, selectedTargets, response := getSelectedTargetsFromCLI(repoName, targetNames, allTargets, interactiveSelection, true, program)
			handleFunctionResponse(response, true)

			response = targetsRunHooks(repo, selectedTargets, targetHooksNames, repoPreHooks, repoPostHooks, notCreateTempDir, notRemoveTempDir, notPrintOutput, false, true, program)
			handleFunctionResponse(response, true)
		},
	}

	targetsHooksRunCmd.Flags().StringVarP(&repoName, "repo", "r", "", "Repo name")
	targetsHooksRunCmd.Flags().StringSliceVarP(&targetNames, "target", "t", nil, "Target(s) name(s)")
	targetsHooksRunCmd.Flags().BoolVarP(&allTargets, "all", "a", false, "Include all targets")
	targetsHooksRunCmd.Flags().BoolVarP(&interactiveSelection, "interactive", "i", false, "Interactive selection")
	targetsHooksRunCmd.Flags().StringSliceVarP(&targetHooksNames, "hook", "k", nil, "Hook(s) name(s)")
	targetsHooksRunCmd.Flags().BoolVarP(&notCreateTempDir, "nocreatetemp", "", false, "Do not create the temporary directory before running the hook(s) (by default, it is created)")
	targetsHooksRunCmd.Flags().BoolVarP(&notRemoveTempDir, "noremovetemp", "n", false, "Do not remove the temporary directory after the hook(s) has/have finished running (by default, it is removed)")
	targetsHooksRunCmd.Flags().StringSliceVarP(&repoPreHooks, "repopre", "", nil, "Repo pre hook(s)")
	targetsHooksRunCmd.Flags().StringSliceVarP(&repoPostHooks, "repopost", "", nil, "Repo post hook(s)")
	targetsHooksRunCmd.Flags().BoolVarP(&notPrintOutput, "quiet", "q", false, "Do not print command output (silent)")
	targetsHooksRunCmd.Flags().SetInterspersed(false)

	// Add Cobra commands
	rootCmd.AddCommand(reposCmd)
	rootCmd.AddCommand(targetsCmd)
	rootCmd.AddCommand(utilitiesCmd)
	rootCmd.AddCommand(showVersionCmd)
	rootCmd.AddCommand(userInitCmd)
	rootCmd.AddCommand(docsCmd)

	docsCmd.AddCommand(docsGenerateCmd)

	utilitiesCmd.AddCommand(utilityMsgCmd)
	utilitiesCmd.AddCommand(utilityAttentionCmd)
	utilitiesCmd.AddCommand(utilityErrorCmd)
	utilitiesCmd.AddCommand(utilitySuccessCmd)
	utilitiesCmd.AddCommand(utilitySectionCmd)
	utilitiesCmd.AddCommand(utilityHrCmd)
	utilitiesCmd.AddCommand(utilityConfirmCmd)
	utilitiesCmd.AddCommand(utilitySSHAgentStartCmd)
	utilitiesCmd.AddCommand(utilitySSHAgentStopCmd)
	utilitiesCmd.AddCommand(utilitySSHAgentGetPIDCmd)
	utilitiesCmd.AddCommand(utilitySSHAgentGetSockCmd)

	reposCmd.AddCommand(reposCreateCmd)
	reposCmd.AddCommand(reposRmCmd)
	reposCmd.AddCommand(reposLsCmd)
	reposCmd.AddCommand(reposServeCmd)
	reposCmd.AddCommand(reposEnableCmd)
	reposCmd.AddCommand(reposDisableCmd)
	reposCmd.AddCommand(reposHooksCmd)

	reposHooksCmd.AddCommand(reposHooksRunCmd)
	reposHooksCmd.AddCommand(reposHooksLsCmd)

	targetsCmd.AddCommand(targetsUpdateCmd)
	targetsCmd.AddCommand(targetsEnableCmd)
	targetsCmd.AddCommand(targetsDisableCmd)
	targetsCmd.AddCommand(targetsCreateCmd)
	targetsCmd.AddCommand(targetsRmCmd)
	targetsCmd.AddCommand(targetsLsCmd)
	targetsCmd.AddCommand(targetsHooksCmd)

	targetsHooksCmd.AddCommand(targetsHooksRunCmd)
	targetsHooksCmd.AddCommand(targetsHooksLsCmd)

	if err := rootCmd.Execute(); err != nil {
		showError("Error: "+err.Error(), program.indentLevel)
		finishProgram(1)
	}
}
