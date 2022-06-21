package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

func initGitRepo() (func(string, ...string), string) {
	tmpDir := makeTempDir()
	repoName := cloneRepo(tmpDir)
	repoDir := filepath.Join(tmpDir, repoName)
	execCommand := func(name string, args ...string) {
		cmd := exec.Command(name, args...)
		cmd.Dir = repoDir
		if err := cmd.Run(); err != nil {
			panic(fmt.Sprintf("Error running command %s %s\n", name, args))
		}
	}
	return execCommand, repoDir
}

func cloneRepo(path string) string {
	cmd := exec.Command("git", "clone", "--depth=1", "--single-branch", "--branch", os.Getenv("GIT_BRANCH"), createGitUrl())
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		panic("Error cloning git repository")
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(fmt.Sprintf("Error reading directory %s\n", path))
	}

	file := files[0]
	if file == nil || !file.IsDir() {
		panic("Path is not a directory")
	}

	repoName := file.Name()
	repoPath := filepath.Join(path, repoName)

	cmd = exec.Command("git", "config", "user.email", os.Getenv("GIT_EMAIL"))
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		panic("Error setting user.email")
	}

	cmd = exec.Command("git", "config", "user.name", os.Getenv("GIT_NAME"))
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		panic("Error setting user.name")
	}

	return repoName
}

func createGitUrl() string {
	result, err := url.Parse(os.Getenv("GIT_REPO"))
	if err != nil {
		panic("Error parsing GITHUB_REPO url")
	}
	result.User = url.UserPassword(os.Getenv("GIT_NAME"), os.Getenv("GIT_PASSWORD"))
	return result.String()
}

func makeTempDir() string {
	tmpDir := filepath.Join(".", ".tmp")

	if err := os.RemoveAll(tmpDir); err != nil {
		panic("Error removing temp directory")
	}

	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		panic("Error creating temp directory")
	}

	return tmpDir
}
