package internal

import (
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const baseURL = "https://golang.org/dl/"

func FetchGoVersions() ([]string, error) {
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`href="/dl/go([0-9]+\.[0-9]+\.[0-9]+)\.linux-amd64\.tar\.gz"`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var versions []string
	for _, match := range matches {
		versions = append(versions, match[1])
	}

	return versions, nil
}

func FilterVersions(versions []string) []string {
	versionMap := make(map[string]string)

	for _, version := range versions {
		parts := strings.Split(version, ".")
		if len(parts) < 3 {
			continue
		}
		majorMinor := parts[0] + "." + parts[1]
		if existingVersion, ok := versionMap[majorMinor]; !ok || isNewerVersion(version, existingVersion) {
			versionMap[majorMinor] = version
		}
	}

	var filteredVersions []string
	for _, version := range versionMap {
		filteredVersions = append(filteredVersions, version)
	}

	sort.Slice(filteredVersions, func(i, j int) bool {
		return compareVersions(filteredVersions[i], filteredVersions[j]) > 0
	})

	return filteredVersions
}

func isNewerVersion(version1, version2 string) bool {
	parts1 := strings.Split(version1, ".")
	parts2 := strings.Split(version2, ".")

	if len(parts1) < 3 || len(parts2) < 3 {
		return false
	}

	minor1, _ := strconv.Atoi(parts1[1])
	minor2, _ := strconv.Atoi(parts2[1])

	return minor1 > minor2
}

func compareVersions(version1, version2 string) int {
	parts1 := strings.Split(version1, ".")
	parts2 := strings.Split(version2, ".")

	if len(parts1) < 3 || len(parts2) < 3 {
		return 0
	}

	major1, _ := strconv.Atoi(parts1[0])
	minor1, _ := strconv.Atoi(parts1[1])
	patch1, _ := strconv.Atoi(parts1[2])

	major2, _ := strconv.Atoi(parts2[0])
	minor2, _ := strconv.Atoi(parts2[1])
	patch2, _ := strconv.Atoi(parts2[2])

	if major1 != major2 {
		return major1 - major2
	}
	if minor1 != minor2 {
		return minor1 - minor2
	}
	return patch1 - patch2
}
