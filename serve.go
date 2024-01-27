package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// External modules
	gin "github.com/gin-gonic/gin"
)

//
//// REPOS (SERVER FUNCTIONALITY)
//

func reposServe(serverPort string, debugMode bool, trustedProxies []string, program Program) functionResponse {
	//
	//// PRINT SERVER INFORMATION
	//

	fmt.Println(program.nameAscii)
	fmt.Println(blue.Sprintf("Starting server"))
	fmt.Println(fmt.Sprintf("Port: %s", paleLime.Sprintf(serverPort)))
	fmt.Println(fmt.Sprintf("Debug mode: %s", paleLime.Sprintf(strconv.FormatBool(debugMode))))
	fmt.Println("Trusted proxies:")
	if len(trustedProxies) > 0 {
		for _, proxy := range trustedProxies {
			fmt.Println(fmt.Sprintf("  - %s", paleLime.Sprintf(proxy)))
		}
	} else {
		showAttention("  - All proxies are trusted", program.indentLevel)
	}
	space()

	//
	//// GIN STARTUP/CONFIGURATION
	//

	if debugMode == false {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.SetTrustedProxies(trustedProxies)

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	//
	//// GIN ROUTES
	//

	router.GET("/", func(c *gin.Context) {
		// Get enabled repos
		repos := getEnabledRepos(program)

		// Create an HTML response
		c.Header("Content-Type", "text/html")
		c.Writer.Write([]byte("<h1> Index of " + "/" + " </h1>"))
		c.Writer.Write([]byte("<pre>\n"))

		for _, repo := range repos {
			// For each repo, create a link and append to the response
			c.Writer.Write([]byte("<a href=\"" + "/repos/" + repo.name + "\">" + repo.name + "</a>\n"))
		}

		c.Writer.Write([]byte("</pre>\n"))
	})

	router.GET("/repos", func(c *gin.Context) {
		// Create a HTML response
		c.Header("Content-Type", "text/html")
		c.Writer.Write([]byte("<h1> Index of /repos </h1>"))
		c.Writer.Write([]byte("<p> To view available repos, click <a href=\"" + "../" + "\">" + "here" + "</a>.</p>"))
	})

	router.GET("/repos/:repo", func(c *gin.Context) {
		repoName := c.Param("repo")
		repo := generateRepoObj(repoName, program)

		// Verify if repo exists or is enabled
		repoVerify := verifyRepoDirectory(repo, program)
		status, response := isRepoDisabled(repo, program)
		handleFunctionResponse(response, true)

		if repoVerify.exitCode != 0 || status == true {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "The requested repo could not be found.",
			})
			return
		}

		// Get enabled targets
		targets := getEnabledTargets(repo, program)

		// Create an HTML response
		c.Header("Content-Type", "text/html")
		c.Writer.Write([]byte("<h1> Index of " + "/repos/" + repo.name + " </h1>"))
		c.Writer.Write([]byte("<pre>\n"))
		c.Writer.Write([]byte("<a href=\"" + "../" + "\">" + "../" + "</a>\n"))

		for _, target := range targets {
			// For each target, create a link and append to the response
			c.Writer.Write([]byte("<a href=\"" + "/repos/" + repo.name + "/" + target.name + "/tree" + "\">" + target.name + "</a>\n"))
		}

		c.Writer.Write([]byte("</pre>\n"))
	})

	router.GET("/repos/:repo/:target", func(c *gin.Context) {
		repoName := c.Param("repo")
		repo := generateRepoObj(repoName, program)
		targetName := c.Param("target")
		target := generateTargetObj(repoName, targetName, program)

		// Create an HTML response
		c.Header("Content-Type", "text/html")
		c.Writer.Write([]byte("<h1> Index of /repos/" + repo.name + "/" + target.name + "</h1>"))
		c.Writer.Write([]byte("<p>These are the available subdirectories for targets:</p>"))
		c.Writer.Write([]byte("<ul>"))
		c.Writer.Write([]byte("<li>" + "<a href=\"" + "/repos/" + repo.name + "/" + target.name + "/tree" + "\">" + "tree" + "</a>" + "</li>"))
		c.Writer.Write([]byte("<li>" + "<a href=\"" + "/repos/" + repo.name + "/" + target.name + "/api" + "\">" + "api" + "</a>" + " (requires an action as a subdirectory, e.g., `api/upload`)</li>"))
		c.Writer.Write([]byte("</ul>"))
	})

	router.GET("/repos/:repo/:target/api", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The API exclusively supports the use of the POST method.",
		})
	})

	router.POST("/repos/:repo/:target/api", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No action was specified",
		})
	})

	router.GET("/repos/:repo/:target/api/:action", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The API exclusively supports the use of the POST method.",
		})
	})

	router.POST("/repos/:repo/:target/api/:action", func(c *gin.Context) {
		repoName := c.Param("repo")
		repo := generateRepoObj(repoName, program)
		targetName := c.Param("target")
		target := generateTargetObj(repoName, targetName, program)
		action := c.Param("action")

		// Verify if repo exists or is enabled
		repoVerify := verifyRepoDirectory(repo, program)
		status, response := isRepoDisabled(repo, program)
		handleFunctionResponse(response, true)

		if repoVerify.exitCode != 0 || status == true {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "The requested repo could not be found.",
			})
			return
		}

		// Verify if target exists or is enabled
		targetVerify := verifyTargetDirectory(target, program)
		status, response = isTargetDisabled(target, program)
		handleFunctionResponse(response, true)

		if targetVerify.exitCode != 0 || status == true {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "The requested target could not be found.",
			})
			return
		}

		//
		////
		//

		if action == "upload" {
			form, err := c.MultipartForm()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "No files were sent.",
				})
				return
			}
			files := form.File["upload[]"]

			for _, file := range files {
				showAttention(fmt.Sprintf("=> Uploading file '%s' to '%s'", file.Filename, target.poolDir+"/"+file.Filename), program.indentLevel)

				c.SaveUploadedFile(file, target.poolDir+"/"+file.Filename)
			}

			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("%d file(s) uploaded.", len(files)),
			})
			return
		} else {
			// Run hook
			if _, err := os.Stat(target.hooksDir + "/" + action); os.IsNotExist(err) {
				c.JSON(http.StatusNotFound, gin.H{
					"message": fmt.Sprintf("Hook '%s' not found", action),
				})
			} else {
				completedCmd, hookResponse := runHook(target.hooksDir+"/"+action, target.environment, false, false, false, false, false, program)

				if hookResponse.exitCode != 0 {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": fmt.Sprintf("Hook finished with the following exit code: %v", hookResponse.exitCode),
						"command": gin.H{
							"exitCode": hookResponse.exitCode,
							"output":   completedCmd.Output,
						},
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"message": "Hook finished running successfully",
						"command": gin.H{
							"exitCode": hookResponse.exitCode,
							"output":   completedCmd.Output,
						},
					})
				}
			}
		}
	})

	router.GET("/repos/:repo/:target/tree/*filepath", func(c *gin.Context) {
		repoName := c.Param("repo")
		repo := generateRepoObj(repoName, program)
		targetName := c.Param("target")
		target := generateTargetObj(repoName, targetName, program)
		resourcePath := c.Param("filepath")
		//resourcePath = strings.TrimPrefix(resourcePath, "/")
		resourcePath = strings.TrimSuffix(resourcePath, "/")

		// Verify if repo exists or is enabled
		repoVerify := verifyRepoDirectory(repo, program)
		status, response := isRepoDisabled(repo, program)
		handleFunctionResponse(response, true)

		if repoVerify.exitCode != 0 || status == true {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "The requested repo could not be found.",
			})
			return
		}

		// Verify if target exists or is enabled
		targetVerify := verifyTargetDirectory(target, program)
		status, response = isTargetDisabled(target, program)
		handleFunctionResponse(response, true)

		if targetVerify.exitCode != 0 || status == true {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "The requested target could not be found.",
			})
			return
		}

		//
		////
		//

		// Join with the base directory to get full file path
		fullPath := filepath.Join(target.poolDir, resourcePath)

		// Check if the path is a directory or a file
		info, err := os.Stat(fullPath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "The requested resource could not be found.",
			})
			return
		}

		if info.IsDir() {
			// If it's a directory, generate a directory listing
			files, err := ioutil.ReadDir(fullPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error: failed to read directory",
				})
				return
			}

			// Create an HTML response
			c.Header("Content-Type", "text/html")
			c.Writer.Write([]byte("<h1> Index of " + "/repos/" + repo.name + "/" + target.name + resourcePath + " </h1>"))
			c.Writer.Write([]byte("<pre>\n"))
			c.Writer.Write([]byte("<a href=\"" + "../" + "\">" + "../" + "</a>\n"))

			for _, file := range files {
				fileName := strings.ReplaceAll(file.Name(), " ", "%20")
				modTime := file.ModTime().Format(http.TimeFormat) // Format as HTTP date

				var fileInfoMD5 string
				var fileInfoSHA256 string
				var size string

				if file.IsDir() {
					size = "-" // Use '-' for directories
					fileInfoMD5 = "-"
					fileInfoSHA256 = "-"

					linkData := fmt.Sprintf("<a href=\"/repos/%s/%s/tree%s/%s\">%s</a>",
						repo.name, target.name, resourcePath, fileName, file.Name())
					fileData := fmt.Sprintf("<p>     %s     %s</p>", modTime, size)
					//fileData = strings.ReplaceAll(fileData, " ", "&nbsp;")
					linedata := linkData + fileData + "\n"
					c.Writer.Write([]byte(linedata))
				} else {
					size = formatBytes(file.Size())

					fileInfoMD5, err = calculateMD5(fullPath + "/" + fileName)
					if err != nil {
						fileInfoMD5 = "<i>Failed to calculate</i>"
					}

					fileInfoSHA256, err = calculateSHA256(fullPath + "/" + fileName)
					if err != nil {
						fileInfoSHA256 = "<i>Failed to calculate</i>"
					}

					linkData := fmt.Sprintf("<a href=\"/repos/%s/%s/tree%s/%s\">%s</a>",
						repo.name, target.name, resourcePath, fileName, file.Name())
					fileData := fmt.Sprintf("<p>     %s     %s</p><p>     <b>MD5:</b>%s</p><p>     <b>SHA256:</b>%s</p>", modTime, size, fileInfoMD5, fileInfoSHA256)
					//fileData = strings.ReplaceAll(fileData, " ", "&nbsp;")
					linedata := linkData + fileData + "\n"
					c.Writer.Write([]byte(linedata))
				}
			}

			c.Writer.Write([]byte("</pre>\n"))
		} else {
			// If it's a file, serve the file
			c.File(fullPath)
		}
	})

	// Listen and serve
	router.Run(fmt.Sprintf(":%s", serverPort))

	return functionResponse{
		exitCode: 0,
	}
}
