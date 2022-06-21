package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

func InitGitRepo() (func(string, ...string), string) {
	tmpDir := MakeTempDir()
	repoName := CloneRepo(tmpDir)
	repoDir := filepath.Join(tmpDir, repoName)
	execCommand := func(name string, args ...string) {
		cmd := exec.Command(name, args...)
		cmd.Dir = repoDir
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error running command %s %s\n", name, args)
		}
	}
	return execCommand, repoDir
}

func CloneRepo(path string) string {
	cmd := exec.Command("git", "clone", "--depth=1", "--single-branch", "--branch", os.Getenv("GITHUB_BRANCH"), CreateGitUrl())
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		log.Fatalln("Error cloning git repository")
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Error reading directory %s\n", path)
	}

	file := files[0]
	if file == nil || !file.IsDir() {
		log.Fatalln("Path is not a directory")
	}

	repoName := file.Name()
	repoPath := filepath.Join(path, repoName)

	cmd = exec.Command("git", "config", "user.email", os.Getenv("GITHUB_EMAIL"))
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		log.Fatalln("Error setting user.email")
	}

	cmd = exec.Command("git", "config", "user.name", os.Getenv("GITHUB_NAME"))
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		log.Fatalln("Error setting user.name")
	}

	return repoName
}

func CreateGitUrl() string {
	result, err := url.Parse(os.Getenv("GITHUB_REPO"))
	if err != nil {
		log.Fatalln("Error parsing GITHUB_REPO url")
	}
	result.User = url.UserPassword(os.Getenv("GITHUB_NAME"), os.Getenv("GITHUB_TOKEN"))
	return result.String()
}

func MakeTempDir() string {
	tmpDir := filepath.Join(".", ".tmp")

	if err := os.RemoveAll(tmpDir); err != nil {
		log.Fatalln("Error removing temp directory")
	}

	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		log.Fatalln("Error creating temp directory")
	}

	return tmpDir
}
