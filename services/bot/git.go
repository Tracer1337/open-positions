package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func CloneRepo() {
	tmpDir := MakeTempDir()

	cmd := exec.Command("git", "clone", "--depth=1", "--single-branch", os.Getenv("GITHUB_REPO"))
	cmd.Dir = tmpDir
	cmd.Run()
}

func MakeTempDir() string {
	tmpDir := filepath.Join(".", ".tmp")

	if err := os.RemoveAll(tmpDir); err != nil {
		log.Fatalln("Error removing temp directory")
	}

	if err := os.MkdirAll(tmpDir, os.ModeDir); err != nil {
		log.Fatalln("Error creating temp directory")
	}

	return tmpDir
}
