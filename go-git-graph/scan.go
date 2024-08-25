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
func scan(path string) {
	fmt.Println("Repo found:")
	repos := getGitFolders(path)
	dotFile := getDotFilePath()

	addNewRepos(dotFile, repos)
	fmt.Println("\n\nRepos saved successfully")
}

// Reset ~/.gogitgraph file
func clearGraph() {
	dotFile := getDotFilePath()
	removeFile(dotFile)
	fmt.Println("\nRepos cleared successfully")
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
func getGitFolders(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

// Get paths for folder containing .git
func scanGitFolders(folders []string, folder string) []string {
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

		if file.Name() == "node_modules" {
			continue
		}

		folders = scanGitFolders(folders, path)
	}

	return folders
}
