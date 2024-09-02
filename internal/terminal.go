package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ensureGoPathInConfig(filePath, path string) error {
	exists, err := pathExists(filePath, path)
	if err != nil {
		return fmt.Errorf("error verifying the configuration file: %w", err)
	}

	if exists {
		fmt.Println("The PATH environment variable is already present in the configuration file:", filePath)
		return nil
	}

	return appendToConfigFile(filePath, path)
}

func pathExists(filePath, path string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("error opening the configuration file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := fmt.Sprintf("export PATH=$PATH:%s", path+"/bin")
		content := scanner.Text()
		if strings.Contains(content, line) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading the configuration file: %w", err)
	}

	return false, nil
}

func appendToConfigFile(filePath, path string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening the configuration file: %w", err)
	}
	defer file.Close()

	line := fmt.Sprintf("export PATH=$PATH:%s\n", path+"/bin")
	_, err = file.WriteString(line)
	if err != nil {
		return fmt.Errorf("error writing to the configuration file: %w", err)
	}

	fmt.Println("PATH environment variable added to the configuration file:", filePath)
	return nil
}

func SetShellEnv() error {
	shell := os.Getenv("SHELL")
	var configFile string
	var path string

	home := os.Getenv("HOME")

	if shell != "" {
		fmt.Println("Current shell according to $SHELL:", shell)
		if strings.Contains(shell, "bash") {
			fmt.Println("You are using bash")
			configFile = home + "/.bashrc"
		} else if strings.Contains(shell, "zsh") {
			fmt.Println("You are using zsh")
			configFile = home + "/.zshrc"
		} else {
			return errors.New("you are using a different shell")
		}

		path = home + "/goinstaller/runtime/go"
		if err := ensureGoPathInConfig(configFile, path); err != nil {
			return fmt.Errorf("error processing the configuration file: %v", err)
		} else {
			fmt.Println("Configuration file processing completed successfully")
			return nil
		}
	}

	fmt.Println("Environment variable $SHELL not found")
	return nil
}
