package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/mitchellh/go-homedir"
	gitignorer "github.com/rliebling/gitignorer"
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
		if !(mode.IsDir() || mode.IsRegular() || mode&os.ModeSymlink == os.ModeSymlink) {
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

		if mode.IsDir() && matched {
			return filepath.SkipDir
		} else if mode.IsDir() || matched {
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

// ParseEnv accepts an `env` slice from the command-line and returns a map.
func ParseEnv(env []string) (map[string]string, error) {
	m := make(map[string]string)

	for _, s := range env {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		} else {
			return nil, fmt.Errorf("environment variable %q is missing a value", parts[0])
		}
	}

	return m, nil
}

// ProfileAndRegionFromConfig attempts to load the .profile setting from `environment`'s config.
func ProfileAndRegionFromConfig(environment string) (string, string, error) {
	configFile := "project.json"

	if environment != "" {
		configFile = fmt.Sprintf("project.%s.json", environment)
	}

	f, err := os.Open(configFile)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	var v struct {
		Profile string `json:"profile"`
		Region  string `json:"region"`
	}

	if err := json.NewDecoder(f).Decode(&v); err != nil {
		return "", "", err
	}

	return v.Profile, v.Region, nil
}

// AssumeRole uses STS to assume the given `role`.
func AssumeRole(role string, config *aws.Config) (*aws.Config, error) {
	stscreds := sts.New(session.New(config))

	params := &sts.AssumeRoleInput{
		RoleArn:         &role,
		RoleSessionName: aws.String("apex"),
		DurationSeconds: aws.Int64(1800),
	}

	res, err := stscreds.AssumeRole(params)
	if err != nil {
		return nil, err
	}

	id := *res.Credentials.AccessKeyId
	secret := *res.Credentials.SecretAccessKey
	token := *res.Credentials.SessionToken

	return &aws.Config{
		Region:      config.Region,
		Credentials: credentials.NewStaticCredentials(id, secret, token),
	}, nil
}
