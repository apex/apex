package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Unknwon/goconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/rliebling/gitignorer"
)

// Sha256 returns a base64 encoded SHA256 hash of `b`.
func Sha256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// LoadFiles return filtered map of relative to 'root' file paths;
// for filtering it uses shell file name pattern matching
func LoadFiles(root string, ignoreFile []byte) (files []string, err error) {
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

		matched, err := gitignorer.GitIgnore(bytes.NewReader(ignoreFile), rel)
		if err != nil {
			return err
		}

		if matched {
			return nil
		}

		files = append(files, rel)

		return nil
	})

	return
}

// GetRegion attempts loading the AWS region from ~/.aws/config.
func GetRegion(profile string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, ".aws", "config")
	cfg, err := goconfig.LoadConfigFile(path)
	if err != nil {
		return "", err
	}

	sectionName := "default"
	if profile != "" && profile != "default" {
		sectionName = fmt.Sprintf("profile %s", profile)
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return "", fmt.Errorf("Could not find AWS region in %s", path)
	}

	return section["region"], nil
}

// ReadIgnoreFile reads .apexignore in `dir` when present and returns a list of patterns.
func ReadIgnoreFile(dir string) ([]byte, error) {
	path := filepath.Join(dir, ".apexignore")

	b, err := ioutil.ReadFile(path)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return b, nil
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
