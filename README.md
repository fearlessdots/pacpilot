# PacPilot

```
 ____            ____  _ _       _
|  _ \ __ _  ___|  _ \(_) | ___ | |_
| |_) / _` |/ __| |_) | | |/ _ \| __|
|  __/ (_| | (__|  __/| | | (_) | |_
|_|   \__,_|\___|_|   |_|_|\___/ \__|
```

`pacpilot` is a user-friendly wrapper that simplifies the management and serving of custom Pacman repositories. With `pacpilot`, you can easily create and manage customized repositories, which act as containers for organizing and distributing specific sets of packages. Within each repository, you can define hooks, which allow you to run custom scripts or commands at specific stages of various operations within pacpilot.

At the heart of `pacpilot` lies its modular approach, enabled by hooks. Hooks are an essential component that allows you to implement your custom logic to the program. This powerful feature empowers you to extend and customize the behavior of not only the serving process but also other operations such as creating and removing repositories and targets. Hooks provide the flexibility to tailor the behavior of `pacpilot` to suit your requirements.

With `pacpilot`, you have the freedom to customize the serving process and other aspects of the program to suit your needs. By leveraging hooks, you can incorporate additional functionality, perform complex transformations, or integrate with external systems. The modular nature of pacpilot ensures that it can adapt to various use cases and provide a flexible solution for managing and serving custom Pacman repositories.

Bearing a name that combines `Pacman` and `pilot`, `pacpilot` aims to guide you through the process of managing and serving custom Pacman repositories. It's about taking control of your package management and making it work for you.

## Key Features

- **Repos and Targets:** `pacpilot` provides a flexible and organized structure for managing and serving custom Pacman repositories. Repos represent a collection of targets that can be used to control the serving process via hooks. It can also hold configuration and other files used across most or all of its targets. The targets also provide their own set of hooks and are the structural elements actually used to group and serve packages.

- **Hooks Support:** Both repos and targets support hooks, allowing you to run custom scripts or commands before or after syncing. The sync process itself is also defined by a hook. This enables you to perform additional actions or customize the synchronization process according to your specific needs. These hooks provide a modular way to extend and customize the synchronization process according to your specific needs. With hooks, you can seamlessly integrate additional functionality, perform complex transformations, or interact with external systems.

- **Templates Support:** Creating new repos and targets is made easier with template support. Templates provide a convenient way to create consistent configurations by predefining common settings and hooks.

- **Utilities for Hooks:** `pacpilot` includes a set of utilities that can be used with hooks to simplify common tasks. These utilities allow you to display messages, show section titles, print colored output, ask for user confirmation, and more. They enhance the functionality of hooks and provide a convenient way to interact with the user during the syncing process.

- **Command-Line Interface and User-Friendly Structure:** With its command-line interface, `pacpilot` offers a straightforward and efficient way to manage your synchronization tasks. The program is designed with a user-friendly structure, providing intuitive commands and comprehensive documentation.

## Installation

To build the program, make sure that Go is installed on your system. Clone the repository or download an archive for a specific version and run the following command in the terminal:

```bash
make build
```

This will create a binary file called `pacpilot`; autocompletion files for `bash`, `zsh`, and `fish` in the directory `./autocompletions`; and markdown documentation files in the directory `./docs`. To install the binary and program files (autocompletion files and documentation):

```bash
sudo make install
```

And to uninstall:

```bash
sudo make uninstall
```

To remove files and directories created by this operation:

```bash
make clean
```

### Custom Destination Directory

By default, all files will be installed to their respectives subdirectories under the `/usr` directory in your root filesystem. However, if you want to set a custom destination directory for the installation, you can use the `DESTDIR` variable when running the `make install` command (also valid for `make uninstall`). For example, here's how you could build this program for Termux:

```bash
make install DESTDIR=/data/data/com.termux/files/usr
```

### PKGBUILD

To simplify the building and installation process on Arch Linux, a PKGBUILD is available in this repository.

To build the package, run the following command in the terminal:
```bash
makepkg -sf
```

If you want to build a package specifically for Termux when using `pacman` as the package manager instead of the default one (`apt`), use the following command:
```bash
TERMUX_BUILD= makepkg -sf
```

To ensure that the PKGBUILD correctly identifies Termux as the build target and sets the installation destination directory correctly, the `TERMUX_BUILD` environment variable should be set. It doesn't require any specific value; simply setting the variable is sufficient for the PKGBUILD to recognize Termux as the build target and handle the installation directory accordingly.

> For more information on how to switch package managers on Termux, refer to the [Termux Wiki](https://wiki.termux.com/wiki/Switching_package_manager).

### Initializing User Data Directory

Before you begin using `pacpilot`, you need to create the user data directory by executing the following command:

```bash
pacpilot init -D <data_dir>
```

The flag `-D` is used to specify a path for the user data directory.

## Usage

The general usage of PacPilot is as follows:

```
pacpilot [command]
```

To get started, run the following command:

```
pacpilot
```

This will display information about the program and provide instructions on how to proceed. To learn more about the available commands and options, you can use the `--help/-h` flag:

```
pacpilot --help
```

This command will provide detailed information about the available commands, their usage, and the available options for each command. It is a useful reference when you need more information on how to use a specific command or what options are available for customization.

Feel free to explore the available commands and options using the `--help` flag to get a better understanding of the functionality provided by PacPilot.

### Documentation

To generate the program's documentation in markdown format, you can use the following command:

```shell
pacpilot docs generate
```

The documentation will be generated in the `./docs` directory by default. You can specify a custom location by using:

```shell
pacpilot docs generate -o <output_dir>
```

### Autocompletion Files

It is possible to generate autocompletion files for this program to be used with the following shells:

- `bash`
- `zsh`
- `fish`
- `powershell`

For example, to build the autocompletion file for `fish`, run:

```shell
pacpilot completion fish > <output_file>
```

> By default, the `Makefile` and, consequentially, the `PKGBUILD` automatically build the completion files for `bash`, `zsh`, and `fish`.

## Documentation

### Available Subcommands

- pacpilot: The main command for the program.
  - version: Show the program's version.
  - init: Create user data directory
  - completion: Generate autocompletion files (`bash`, `zsh`, `fish`, and `powershell`)
  - docs: Program documentation.
    - generate: Generate program documentation (markdown files).
  - utils: Utilities for hooks execution.
    - msg: Print a message in a specific color given a HEX code.
    - attention: Display attention message.
    - error: Display error message.
    - success: Display success message.
    - section: Display section title.
    - hr: Display a horizontal line.
    - confirm: Ask for confirmation.
    - sshagent-start: Starts an ssh-agent process and adds a private key.
    - sshagent-stop: Stops the ssh-agent process.
    - sshagent-getpid: Get the process ID of the ssh-agent.
    - sshagent-getsock: Get the socket path of the ssh-agent.
  - repos: Manage repos.
    - enable: Enable repos.
    - disable: Disable repos.
    - create: Create repos.
    - rm: Remove repos.
    - ls: List repos.
    - hooks: Manage repo hooks.
      - run: Run repo hook(s).
      - ls: List repo hooks.
	- serve: Serve repo targets.
  - targets: Manage targets.
    - ls: List targets.
    - enable: Enable targets.
    - disable: Disable targets.
    - create: Create targets.
    - rm: Remove targets.
    - hooks: Manage target hooks.
      - run: Run target hook(s).
      - ls: List target hooks.
	- update: Update targets.

### User Data Directory

There is no default path for the user data directory. Instead, it is specified via the `-D` flag in the command line. However, its default tree structure is as follows:

- `templates`: This directory contains templates used for creating repos and targets. It provides a starting point with pre-configured setups for common scenarios. The templates are organized into subdirectories based on their type, such as `repos` and `targets`.

- `templates/repos`: This subdirectory within the `templates` directory contains templates specifically designed for creating repos. Each template may include a set of pre-defined hook scripts and configuration files to streamline the repo creation process.

- `templates/targets`: This subdirectory within the `templates` directory contains templates specifically designed for creating targets. Similar to the repo templates, each target template may include hook scripts and configuration files tailored for specific needs.

- `repos`: This directory holds the configurations and settings for all created repos. Each repo has its own subdirectory within the `repos` directory. The subdirectories are named after the respective repo and contain the associated configuration files, hooks, and any other necessary files.

- `repos/<repo>/targets`: Within each repo's subdirectory, there is a `targets` directory. This directory holds the configurations and hooks for all the targets associated with that particular repo. Each target has its own subdirectory within the `targets` directory, containing the target-specific configuration files, hooks, and any other necessary files.

### Repos

> I created a separation between `repos` and `targets` because it allows for the creation of different repositories for different operating systems. For example, in my use case, I have one repository for Arch Linux and another for Termux, both of which use pacman as the package manager (**note:** on Termux, `pacman` is not used by default). By creating separate targets for each supported architecture (such as x86_64 and aarch64), I can ensure that packages are built and distributed correctly for each platform. This separation also makes it easier to manage and maintain the repositories, as each one can be tailored to the specific needs of its associated operating system and architecture.

#### Creating a Repo

To create a new repo, you can utilize the `repos create` command. This command will guide you through the process by prompting for a repo name and allowing you to choose a template. Templates provide pre-configured setups for common scenarios.

If there are no repo templates available or if you prefer to have a minimal setup with just the repo directory and a `targets` subdirectory, you can select the `scratch` template. This template will only create the necessary directories and will not include any additional configuration files or hooks.

After creating the repo, a new directory will be generated specifically for that repo. This directory will contain the relevant configuration files and hooks, based on the selected template.

#### Managing Repos

Once you have created repos, you can perform various operations on them. The following commands are available for managing repos:

- `repos rm`: Remove repos.
- `repos ls`: List repos.
- `repos enable`: Enable repos.
- `repos disable`: Disable repos.
- `repos hooks`: Manage repo hooks.
 - `repos hooks run`: Run repo hook(s).
 - `repos hooks ls`: List repo hooks.

#### Disabled Repos

When executing the `repos hooks run` command or other commands that depend on hooks to function, `pacpilot` will check if the repo is disabled before running hooks for each repo. If any disabled repos are encountered, `pacpilot` will skip them without generating an error. It will proceed to the remaining selected repos, if any are available.

#### Hooks

Hooks are an essential component of repos, allowing you to define custom actions or scripts to execute. The following default hooks can be associated with a repo:

- `post_create`: Runs after a repo is created. Can be used to further configure the repo configuration beyond only creating its directory (which is done automatically by the program).
- `pre_rm`: Executes before removing a repo.
- `ls`: Displays custom information for the repo when running `repos ls`.

##### Environment Variables

When running repo hooks, the following environment variables are available for your use:

- **PROGRAM_NAME**: The name of the program (`pacpilot`).
- **DEFAULT_SHELL**: The default shell used by the program.
- **PACPILOT_EXEC**: The executable path of `pacpilot`.
- **PACPILOT_UTILS**: The executable path of `pacpilot utils`, providing access to the utility commands.
- **USER_DATA_DIR**: The user data directory used by `pacpilot`.
- **USER_REPOS_DIR**: The directory where user repos are stored.
- **USER_TEMPLATES_DIR**: The directory where user templates are stored.
- **USER_REPOS_TEMPLATES_DIR**: The directory where user repo templates are stored.
- **USER_TARGETS_TEMPLATES_DIR**: The directory where user target templates are stored.
- **REPO_NAME**: The name of the current repo.
- **REPO_DIR**: The directory path of the current repo.
- **REPO_HOOKS_DIR**: The directory path of the hooks within the current repo.
- **REPO_TARGETS_DIR**: The directory path of the targets within the current repo.
- **REPO_TEMP_DIR**: The temporary directory path specific to the current repo.

These environment variables provide useful information and paths that can be utilized within your repo hooks to customize the behavior and perform specific actions based on the current context.

##### Custom entry command

A custom entry command for a specific hook can be defined by creating a file called `<hook_name>.entry` in the repo's hooks directory. For example, if you want to run the `post_create` hook with the `fish` shell, you could do so by creating a `post_create.entry` file with the following content:

```bash
/usr/bin/fish
```

##### Managing Hooks

When managing hooks at the repo level, you can use the following subcommands:

- `repos hooks ls`: List all hooks available for a specific repo. This command displays the names of the hooks and their associated entry commands (if any).

- `repos hooks run`: Run one or more hooks for the specified repos. You can also run hooks in a specific sequence by running, for example:

```bash
repos hooks run --repo <repo_name> --hook <first_hook> --hook <second_hook>
```

##### Temporary Directory

When running hooks for repos, a temporary directory is created for each repo. This temporary directory serves as a workspace for performing actions or modifications during the hook execution process. It is called `.tmp` and created within the repo directory itself. This repo-specific temporary directory provides a separate workspace for any temporary files or data specific to the individual repo and can be used for any shared temporary files or data required by the hook scripts across multiple targets within the same repo. It allows the hooks to operate within a context that is isolated to the repo directory.

Once the hook execution is complete, the temporary directory and its contents are typically cleaned up automatically. When running hooks via `repos hooks run`, this behavior can be modified using two options (non-exclusive):

- `nocreatetemp`: By using this option, temporary directories will not be created before running the hook(s). This can be useful if you prefer to handle temporary directory creation manually or if your hook scripts do not require a separate workspace.

- `noremovetemp`: With this option, temporary directories will not be automatically removed after running the hook(s). This allows you to inspect or access the temporary directories and their contents after the hook execution has completed. It can be beneficial for debugging purposes or if you need to access the temporary files generated during the hook execution.

### Targets

> I created a separation between `repos` and `targets` because it allows for the creation of different repositories for different operating systems. For example, in my use case, I have one repository for Arch Linux and another for Termux, both of which use pacman as the package manager (**note:** on Termux, `pacman` is not used by default). By creating separate targets for each supported architecture (such as x86_64 and aarch64), I can ensure that packages are built and distributed correctly for each platform. This separation also makes it easier to manage and maintain the repositories, as each one can be tailored to the specific needs of its associated operating system and architecture.

#### Creating a Target

To create a new target, you can use the `targets create` command. This command will guide you through the process by allowing you to select a parent repo and provide a name for the target. Additionally, you can choose a template to pre-configure the target's setup according to your requirements.

If there are no target templates available or if you prefer a minimal setup for your target, you can select the `scratch` template during the target creation process. This template will create only the target's directory. It will not include any additional configuration files or hooks, providing you with a clean slate to customize according to your needs.

#### Managing Targets

Once you have created targets, you can perform various operations on them. The following commands are available for managing targets:

- `targets ls`: List targets.
- `targets enable`: Enable targets.
- `targets disable`: Disable targets.
- `targets rm`: Remove targets.
- `targets hooks`: Manage target hooks.
 - `targets hooks run`: Run target hook(s).
 - `targets hooks ls`: List target hooks.
- `targets update`: Update targets.

#### Disabled Targets

When executing the `targets hooks run` command or other commands that depend on hooks to function (like `targets edit` and `targets view`), `pacpilot` will check if the target is disabled before initiating the hooks execution for each target. It will also verify if the repo is disabled before iterating through the targets. If any disabled targets are encountered, `pacpilot` will skip them without generating an error. It will proceed to execute hooks for the remaining selected targets, if any are available (only if the repo is enabled).

#### Hooks

Similar to repos, targets also support hooks that allow you to define custom actions or scripts to execute at specific stages of hooks execution. The default hooks available for targets include:

- `post_create`: Runs after a target is created. Can be used to further configure the target configuration beyond only creating its directory (which is done automatically by the program).
- `pre_rm`: Executes before removing a target.
- `ls`: Displays custom information for the target when running `targets ls`.
- `update`: Displays custom information for the target when running `targets update`.

##### Environment Variables

When running target hooks, the following environment variables are available for your use:

- **PROGRAM_NAME**: The name of the program (`pacpilot`).
- **DEFAULT_SHELL**: The default shell used by the program.
- **PACPILOT_EXEC**: The executable path of `pacpilot`.
- **PACPILOT_UTILS**: The executable path of `pacpilot utils`, providing access to the utility commands.
- **USER_DATA_DIR**: The user data directory used by `pacpilot`.
- **USER_REPOS_DIR**: The directory where user repos are stored.
- **USER_TEMPLATES_DIR**: The directory where user templates are stored.
- **USER_REPOS_TEMPLATES_DIR**: The directory where user repo templates are stored.
- **USER_TARGETS_TEMPLATES_DIR**: The directory where user target templates are stored.
- **REPO_NAME**: The name of the current repo.
- **REPO_DIR**: The directory path of the current repo.
- **REPO_HOOKS_DIR**: The directory path of the hooks within the current repo.
- **REPO_TARGETS_DIR**: The directory path of the targets within the current repo.
- **REPO_TEMP_DIR**: The temporary directory path specific to the current repo.
- **TARGET_NAME**: The name of the current target.
- **TARGET_DIR**: The directory path of the current target.
- **TARGET_HOOKS_DIR**: The directory path of the hooks within the current target.
- **TARGET_TEMP_DIR**: The temporary directory path specific to the current target.
- **TARGET_POOL_DIR**: The directory path where the packages (e.g., `*.pkg.tar.xz`) are stored and served from.

These environment variables provide useful information and paths that can be utilized within your target hooks to customize the behavior and perform specific actions based on the current context.

##### Custom entry command

A custom entry command for a specific hook can be defined by creating a file called `<hook_name>.entry` in the target's hooks directory. For example, if you want to run the `post_create` hook with the `fish` shell, you could do so by creating a `post_create.entry` file with the following content:

```bash
/usr/bin/fish
```
##### Managing Hooks

When working with targets, you have the option to use the following subcommands to manage their hooks:

- `targets hooks ls`: This command lists all available hooks for the specified targets. It displays the names of the hooks and their associated entry commands (if any).

- `targets hooks run`: Use this command to run one or more hooks for the specified targets. You can also specify a specific sequence for the hooks by running the command as follows:

```bash
targets hooks run --repo <repo_name> --target <target_name> --hook <first_hook> --hook <second_hook>
```

Additionally, the `targets hooks run` command allows you to specify an array of repo hooks to be executed before and after running the targets' hooks. You can achieve this by using the following two flags:

- `repopre`: Specifies the hooks to be run before iterating through the targets.
- `repopost`: Specifies the hooks to be run after iterating through the targets.

Similar to the `--hooks/-k` flag, you can add multiple flags for both options. The hooks will be executed in the sequence they were specified in the command line. For example, you could run the following command:

```bash
targets hooks run --repo <repo_name> --target <target_name> --repopre pre_build --repopost post_build --hook build --hook install --hook clean
```

##### Temporary Directory

When running hooks for targets, two temporary directories are created. These temporary directories serve as a workspace for performing actions or modifications during the hook execution process.

- The first one is the `.tmp` directory within the repo directory, similar to the repo hooks. This temporary directory is used for any shared temporary files or data required by the hook scripts across multiple targets within the same repo.

- Additionally, a second temporary directory called `.tmp` is created within the target directory itself. This target-specific temporary directory provides a separate workspace for any temporary files or data specific to the individual target. It allows the hooks to operate within a context that is isolated to the target directory.

Once the hook execution is complete, the temporary directories and their contents are typically cleaned up automatically. When running hooks via `targets hooks run`, this behavior can be modified using two options (non-exclusive):

- `nocreatetemp`: By using this option, temporary directories will not be created before running the hook(s). This can be useful if you prefer to handle temporary directory creation manually or if your hook scripts do not require a separate workspace.

- `noremovetemp`: With this option, temporary directories will not be automatically removed after running the hook(s). This allows you to inspect or access the temporary directories and their contents after the hook execution has completed. It can be beneficial for debugging purposes or if you need to access the temporary files generated during the hook execution.

### Serving Packages

#### Starting the Server

To start the server, run:

```
pacpilot -D <data_dir> repos serve
```

This will call the `reposServe` function with the following parameters:

- `serverPort`: The port number on which the server should listen for incoming requests. Default: **8080**.
- `debugMode`: A boolean value indicating whether the server should run in debug mode or not. In debug mode, the server will provide more detailed error messages and stack traces. Default: **false**.
- `trustedProxies`: An optional list of trusted proxy IP addresses or CIDR blocks. If provided, the server will trust the X-Forwarded-For header from these proxies when determining the client's IP address. Default: **all**.

The function prints some server information, including the port number, debug mode, and trusted proxies. It then initializes a new Gin router and configures it with the provided settings.

#### Server Routes

The server has several routes that handle different types of requests:

##### Root Route (`/`)

This route returns a simple HTML page that lists the available repositories. It gets the list of enabled repositories and generates an HTML response with links to each repository's directory.

##### Repositories Route (`/repos`)

This route returns an HTML page that informs the user that they are viewing the `/repos` directory. It provides a link back to the root directory.

##### Repository Route (`/repos/:repo`)

This route returns an HTML page that lists the available targets for a specific repository. It first verifies that the repository exists and is enabled. If the repository is valid, it gets the list of enabled targets and generates an HTML response with links to each target's directory.

##### Repository Target Route (`/repos/:repo/:target`)

This route returns an HTML page that provides information about the available subdirectories for a specific target in a repository. It lists the tree and api subdirectories and their purposes.

##### Repository Target API Route (`/repos/:repo/:target/api`)

This route is used to handle API requests for a specific target in a repository. It only supports POST requests and returns a `400 Bad Request` response for GET requests.

##### Repository Target API Action Route (`/repos/:repo/:target/api/:action`)

This route is used to handle specific API actions for a target in a repository. It supports POST requests only and returns a `400 Bad Request` response for GET requests. The `:action` parameter specifies the action to be performed.

#### API Actions

##### Upload Action

The `upload` action allows users to upload files to a specific target in a repository. When a POST request is made to the `/repos/:repo/:target/api/upload` route, the server expects a `multipart/form-data` request with one or more files in the `upload[]` field. The server saves the uploaded files to the target's pool directory and returns a JSON response indicating the number of files uploaded.

##### Running Target Hooks

To run a target hook using the API, you can send a POST request to the `/repos/:repo/:target/api/:hook` route, where `:hook` is the name of the hook you want to run. The server will then execute the hook script located at `target.hooksDir + "/" + action` and return a JSON response indicating the exit code and output of the command.

###### Example

Suppose you have a target hook named `update` that updated the database files for your repository. To run this hook using the API, you can use the `curl` command with the `-X POST` option:

```
curl -X POST http://localhost:8080/repos/your-repo/your-target/api/update
```

Replace `your-repo` and `your-target` with the names of your repository and target, respectively.

The server will then execute the `update` hook script located at `target.hooksDir + "/" + "build"` and return a JSON response indicating the exit code and output of the command. If the hook finishes successfully, the response will look something like this:

```json
{
  "message": "Hook finished running successfully",
  "command": {
    "exitCode": 0,
    "output": "Hook output\n"
  }
}
```

If the hook finishes with a non-zero exit code, the response will look something like this:

```json
{
  "message": "Hook finished with the following exit code: 1",
  "command": {
    "exitCode": 1,
    "output": "Hook output\n"
  }
}
```

**Note:** the hook scripts must have the appropriate permissions to be executable. You can set the permissions using the `chmod` command.

By following these steps, you can easily run target hooks using the API provided by the server.

###### Example

To upload a file using the `upload` action, you can use the `curl` command with the `-F` option to specify the file to be uploaded:

```
curl -X POST -F "upload[]=@/path/to/your/file.pkg.tar.zst" http://localhost:8080/repos/your-repo/your-target/api/upload
```

Replace `/path/to/your/file.pkg.tar.zst` with the path to the file you want to upload, and replace `your-repo` and `your-target` with the names of your repository and target, respectively.

#### Serving Repository Files

The server also handles requests for individual files in the repositories. When a GET request is made to a URL in the format `/repos/:repo/:target/tree/*filepath`, the server checks if the requested resource exists and serves it if it does. The `*filepath` parameter specifies the path to the file or directory relative to the target's pool directory.

If the requested resource is a directory, the server generates an HTML directory listing that includes links to the files and subdirectories within the directory. If the requested resource is a file, the server serves the file directly.

When serving repository files, the server provides some useful information about each file in the directory listing. Here's what each piece of information means:

1. **File Name**: The name of the file or directory. If the entry is a directory, it will be indicated by a trailing slash (`/`).
2. **Last Modified Time**: The date and time when the file was last modified, in the format `Mon, 02 Jan 2006 15:04:05 MST`. This information is displayed in the HTTP date format.
3. **Size**: The size of the file in bytes. If the file is a directory, the size will be displayed as `-`.
4. **MD5 Hash**: The MD5 hash of the file, displayed as a 32-character hexadecimal string. This information can be used to verify the integrity of the file.
5. **SHA256 Hash**: The SHA256 hash of the file, displayed as a 64-character hexadecimal string. This information can be used to verify the integrity of the file with a higher level of security than the MD5 hash.

By providing this information, the server makes it easy for users to identify and verify the files in their repositories. Note that the server calculates the MD5 and SHA256 hashes on the fly, so there may be a slight delay when serving large files or directories with many files.

### Utilities

`pacpilot` provides a set of utilities designed to be used within the hooks of repos and targets, allowing you to perform additional actions or execute custom logic during operations, though they can be used wherever and whenever you want. The main difference is that when running repo and target hooks, an environment variable called `$PACPILOT_UTILS` is automatically created, pointing to `pacpilot utils`.

#### Available Utilities

- **msg**: Print a message in a specific color given a HEX code.
  ```
  pacpilot utils msg '#FF5733' 'This is a colored message'
  ```

- **attention**: Display an attention message.
  ```
  pacpilot utils attention 'This is an attention message'
  ```

- **error**: Display an error message.
  ```
  pacpilot utils error 'This is an error message'
  ```

- **success**: Display a success message.
  ```
  pacpilot utils success 'This is a success message'
  ```

- **section**: Display a section title.
  ```
  pacpilot utils section 'This is a section title'
  ```

- **hr**: Display a horizontal line.
  ```
  pacpilot utils hr '-' 0.45
  ```

- **confirm**: Ask for confirmation from the user.
  ```
  pacpilot utils confirm 'Are you sure you want to continue?'
  ```

- **sshagent-start**: Start an ssh-agent process and add a private key.
  ```
  pacpilot utils sshagent-start '/path/to/my/key' $TARGET_TEMP_DIR
  ```

- **sshagent-stop**: Stop the ssh-agent process.
  ```
  pacpilot utils sshagent-stop $TARGET_TEMP_DIR
  ```

- **sshagent-getpid**: Get the process ID of the ssh-agent.
  ```
  pacpilot utils sshagent-getpid $TARGET_TEMP_DIR
  ```

- **sshagent-getsock**: Get the socket path of the ssh-agent.
  ```
  pacpilot utils sshagent-getsock $TARGET_TEMP_DIR
  ```

#### Using Utilities in Hooks

To use any of the utilities within a repo or target hook, you can access them using the `$PACPILOT_UTILS` environment variable, which points to the command `pacpilot utils`. For example, to display an attention message within a repo hook:

```bash
$PACPILOT_UTILS attention 'This is an attention message'
```

This will print the specified attention message to the console, attracting the user's attention during the operation.

### Templates

Templates play a crucial role in customizing the creation of repos and targets. `pacpilot` provides a template-based approach to create repos and targets, allowing you to quickly set up and configure your project structure. When creating a repo or target, `pacpilot` automatically generates the corresponding repo or target directory and copies the selected template structure to it.
To simplify template usage, you can find a collection of example repo and target templates that I personally use in the `./templates` directory within the source code. These templates serve as starting points and can be customized to suit your specific project requirements. Additionally, you have the flexibility to create your own templates and store them in the following directories within the user data directory:

- `templates/repos`: This directory is dedicated to repo templates.
- `templates/targets`: This directory is dedicated to target templates.

By placing your custom templates in these directories, they become readily available for selection during the repo and target creation process. You can leverage these templates to expedite the setup of your projects and tailor them to your specific needs.

## License

PacPilot is licensed under the GPL-3.0 license.
