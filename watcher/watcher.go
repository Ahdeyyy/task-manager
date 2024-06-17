package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Define the directory to watch and the command to run
	watchDir := "./"
	command := "sh build.sh; go run ." // Replace with your command

	// Store the initial state of the files
	filesState := make(map[string]time.Time)
	err := filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".templ") || strings.HasSuffix(path, ".go")) {
			filesState[path] = info.ModTime()
			fmt.Println("Watching", path)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	runCommand(command)
	// Continuously monitor the directory for changes
	for {
		err = filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (strings.HasSuffix(path, ".templ") || strings.HasSuffix(path, ".go")) {
				lastModTime, exists := filesState[path]
				if !exists || lastModTime != info.ModTime() {
					filesState[path] = info.ModTime()
					fmt.Printf("Detected change in %s, running command...\n", path)
					runCommand(command)
				}
			}
			return nil
		})
		if err != nil {
			log.Println("Error walking the directory:", err)
		}

		time.Sleep(100 * time.Millisecond) // Adjust the polling interval as needed
	}
}

// Run the specified command and print its output
func runCommand(command string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println("Command execution failed:", err)
	}
}
