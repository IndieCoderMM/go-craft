package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

// Save all git repos within a folder into ~/.gogitgraph file
func scan(path string, ignoreFolders string) {
	fmt.Println("🔍 Scanning repos...")
	repos := getGitFolders(path, strings.Split(ignoreFolders, ","))
	dotFile := getDotFilePath()

	addNewRepos(dotFile, repos)
	fmt.Printf("\n\n💾 %d Repos registered", len(repos))
}

// List all registerd repos
func listRepos() {
	dotFile := getDotFilePath()
	repos := readRepoOnly(dotFile)

	if len(repos) == 0 {
		fmt.Println("🗂 No repos registered")
	} else {
		fmt.Println("🗂 Registered repos:")
		for _, repo := range repos {
			fmt.Println(repo)
		}
	}
}

// Reset ~/.gogitgraph file
func clearGraph() {
	dotFile := getDotFilePath()
	removeFile(dotFile)
	fmt.Println("\n 🗑 All repos cleared")
}

// Remove a given file
func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

// Append new repos into a file
func addNewRepos(path string, repos []string) {
	existingRepos := readFile(path)
	newRepos := joinSlices(repos, existingRepos)
	saveRepos(newRepos, path)
}

// Write repos into a file line by line
func saveRepos(repos []string, path string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(path, []byte(content), 0755)
}

// Save default configs into ~/.gogitgraph file
func saveConfigs(email string) {
	dotFile := getDotFilePath()

	// First line is email
	// Second line is ignore folders
	// TODO: Add configuration for ignore folders
	var ignoreFolders = []string{"node_modules"}

	configs := []string{"#email:" + email, "#ignore:" + strings.Join(ignoreFolders, ",")}
	repos := readRepoOnly(dotFile)
	saveRepos(append(configs, repos...), dotFile)
}

// Show current configs
func getConfigs(print bool) (map[string]string, error) {
	dotFile := getDotFilePath()
	lines := readFile(dotFile)
	configs := make(map[string]string)

	if len(lines) == 0 {
		fmt.Println("📝 No configs found")
		panic("🛠 Please set your email with '-email'")
	}

	for _, line := range lines {
		typ, value := parseConfig(line)

		configs[typ] = value

		if !print {
			continue
		}

		if typ == "email" {
			fmt.Println("📧 Email:", value)
		} else if typ == "ignore" {
			fmt.Println("🚫 Ignore folders:", value)
		}
	}

	return configs, nil
}

// Parse config line
func parseConfig(line string) (string, string) {
	if line[0] != '#' {
		return "", ""
	}

	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return "", ""
	}

	return parts[0][1:], parts[1]
}

func readRepoOnly(path string) []string {
	lines := readFile(path)
	repos := []string{}

	for _, line := range lines {
		if line[0] == '#' {
			continue
		}
		repos = append(repos, line)
	}

	return repos
}

// Read file content into a slice
func readFile(path string) []string {
	f := openFile(path)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	return lines
}

// Append new slice into current slice
func joinSlices(new []string, curr []string) []string {
	for _, n := range new {
		if !includes(curr, n) {
			curr = append(curr, n)
		}
	}

	return curr
}

// Check if a slice contains a value
func includes(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}

// Open a file, create if not exist
func openFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	return f
}

// Get ~/.gogitgraph filepath
func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "/.gogitgraph"

	return dotFile
}

// Get path for all git repos within a folder
func getGitFolders(folder string, ignoreFolders []string) []string {
	return scanGitFolders(make([]string, 0), folder, ignoreFolders)
}

// Get paths for folder containing .git
func scanGitFolders(folders []string, folder string, ignoreFolders []string) []string {
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		path = folder + "/" + file.Name()
		if file.Name() == ".git" {
			path = strings.TrimSuffix(path, "/.git")
			fmt.Println(path)
			folders = append(folders, path)
			continue
		}

		if includes(ignoreFolders, file.Name()) {
			continue
		}

		folders = scanGitFolders(folders, path, ignoreFolders)
	}

	return folders
}
