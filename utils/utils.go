package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
)

// Sha256 returns a base64 encoded SHA256 hash of `b`.
func Sha256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// LoadFiles return filtered map of relative to 'root' file paths;
// for filtering it uses shell file name pattern matching
func LoadFiles(root string, ignoredPatterns []string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		mode := info.Mode()
		if !(mode.IsRegular() || mode&os.ModeSymlink == os.ModeSymlink) {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		for _, pattern := range ignoredPatterns {
			matched, err := filepath.Match(pattern, rel)
			if err != nil {
				return err
			}

			if matched {
				return nil
			}
		}

		files = append(files, rel)

		return nil
	})

	return
}

// GetProfile attempts to load the profile from AWS_PROFILE otherwise defaults to "default"
func GetProfile() string {
	profile := os.Getenv("AWS_PROFILE")

	if profile == "" {
		return "default"
	}

	return profile
}

// GetRegion attempts loading the AWS region from ~/.aws/config.
func GetRegion(profile string) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	path := filepath.Join(u.HomeDir, ".aws", "config")
	cfg, err := goconfig.LoadConfigFile(path)
	if err != nil {
		return "", err
	}

	sectionName := "default"
	if profile != "" {
		sectionName = fmt.Sprintf("profile %s", profile)
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return "", fmt.Errorf("Could not find AWS region in %s", path)
	}

	return section["region"], nil
}

// ReadIgnoreFile reads .apexignore in `dir` when present and returns a list of patterns.
func ReadIgnoreFile(dir string) ([]string, error) {
	path := filepath.Join(dir, ".apexignore")

	b, err := ioutil.ReadFile(path)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return strings.Split(string(b), "\n"), nil
}

// ContainsString checks if array contains string
func ContainsString(array []string, element string) bool {
	for _, e := range array {
		if element == e {
			return true
		}
	}
	return false
}
