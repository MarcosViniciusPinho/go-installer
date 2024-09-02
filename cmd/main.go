package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/MarcosViniciusPinho/go-installer/internal"
)

func main() {

	versions, err := internal.FetchGoVersions()
	if err != nil {
		log.Fatalf("Error retrieving Go versions: %v", err)
	}

	filteredVersions := internal.FilterVersions(versions)

	fmt.Print("\nCurrently available only for the Linux operating system with x86-64 (amd64) architecture\n\n")

	fmt.Println("Available Go versions:")
	fmt.Println("------------------------------------")
	fmt.Printf("%-5s %-10s %-15s\n", "ID", "Version", "Installed")
	fmt.Println("------------------------------------")

	goInstallerDir, runtimeDir, cacheDir, err := internal.GetPathDirectories()
	if err != nil {
		log.Fatalf("Error retrieving the current user's directory paths: %v", err)
	}

	for i, version := range filteredVersions {

		fileName := fmt.Sprintf("go%s.linux-amd64.tar.gz", version)
		cacheFile := filepath.Join(*cacheDir, fileName)

		if _, err := os.Stat(cacheFile); !os.IsNotExist(err) {
			fmt.Printf("%-5d %-10s %-15s\n", i+1, version, "✅")
		} else {
			fmt.Printf("%-5d %-10s %-15s\n", i+1, version, "")
		}
	}
	fmt.Println("------------------------------------")

	var choice int
	fmt.Print("Choose the version to be installed (ID): ")
	fmt.Scan(&choice)
	if choice < 1 || choice > len(versions) {
		fmt.Println("Invalid choice")
		return
	}

	selectedVersion := filteredVersions[choice-1]

	fmt.Print("\n\n")

	if err := internal.DownloadAndInstallGo(*goInstallerDir, *runtimeDir, *cacheDir, selectedVersion); err != nil {
		log.Fatalf("Error installing Go: %v", err)
	}

	if err := internal.SetShellEnv(); err != nil {
		log.Fatalf("Error configuring environment variables in the shell: %v", err)
	}

	fmt.Println("Installation completed successfully!")
	fmt.Printf("The version %s is in use - Open a new terminal and type the command 'go version' to confirm the version\n", selectedVersion)
}
