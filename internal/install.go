package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func GetPathDirectories() (*string, *string, *string, error) {
	homeDir, err := getUserDir()
	if err != nil {
		return nil, nil, nil, err
	}

	goInstallerDir := filepath.Join(homeDir, "goinstaller")
	runtimeDir := filepath.Join(goInstallerDir, "runtime")
	cacheDir := filepath.Join(runtimeDir, "cache")

	return &goInstallerDir, &runtimeDir, &cacheDir, nil
}

func DownloadAndInstallGo(goInstallerDir, runtimeDir, cacheDir, version string) error {
	installDirReturn, cacheDirReturn, err := createDirectories(
		goInstallerDir,
		runtimeDir,
		cacheDir,
	)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("go%s.linux-amd64.tar.gz", version)
	cacheFile := filepath.Join(*cacheDirReturn, fileName)

	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		url := fmt.Sprintf("https://golang.org/dl/%s", fileName)
		err = downloadFile(url, cacheFile)
		if err != nil {
			return err
		}
	}

	err = extractTarGz(cacheFile, *installDirReturn)
	if err != nil {
		return err
	}

	return nil
}

func getUserDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func createDirectories(goInstallerDir, runtimeDir, cacheDir string) (*string, *string, error) {
	homeDir, err := getUserDir()
	if err != nil {
		return nil, nil, err
	}

	if _, err := os.Stat(goInstallerDir); os.IsNotExist(err) {
		fmt.Println("Creating directory:", goInstallerDir)
		if err := os.Mkdir(goInstallerDir, 0755); err != nil {
			return nil, nil, fmt.Errorf("error creating directory goInstaller: %w", err)
		}
	}

	if _, err := os.Stat(runtimeDir); os.IsNotExist(err) {
		fmt.Println("Creating directory:", runtimeDir)
		if err := os.Mkdir(runtimeDir, 0755); err != nil {
			return nil, nil, fmt.Errorf("error creating directory runtime: %w", err)
		}
	}

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Println("Creating directory:", cacheDir)
		if err := os.Mkdir(cacheDir, 0755); err != nil {
			return nil, nil, fmt.Errorf("error creating directory cache: %w", err)
		}
	}

	installDir := filepath.Join(homeDir, "goinstaller", "runtime")
	goDir := filepath.Join(installDir, "go")

	if _, err := os.Stat(goDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(goDir); err != nil {
			return nil, nil, fmt.Errorf("error removing directory go: %w", err)
		}
	}

	return &installDir, &cacheDir, nil
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	bar := progressbar.NewOptions(int(resp.ContentLength),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("Downloading..."),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetRenderBlankState(true),
	)

	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		return err
	}

	bar.Finish()

	fmt.Print("\n")

	return nil
}

func extractTarGz(src, dest string) error {
	cmd := exec.Command("tar", "-xzf", src, "-C", dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error trying to unzip the file: %s: %w", string(output), err)
	}
	return nil
}
