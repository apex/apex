package root

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"

	"github.com/Unknwon/goconfig"
	"github.com/apex/apex/dryrun"
	"github.com/apex/apex/project"
)

// chdir working directory.
var chdir string

// dryRun enabled.
var dryRun bool

// logLevel specified.
var logLevel string

// profile for AWS creds.
var profile string

// Env supplied.
var Env []string

// Session instance.
var Session *session.Session

// Project instance.
var Project *project.Project

// Register `cmd`.
func Register(cmd *cobra.Command) {
	Command.AddCommand(cmd)
}

// Command config.
var Command = &cobra.Command{
	Use:               "apex",
	PersistentPreRunE: preRun,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

// Initialize.
func init() {
	f := Command.PersistentFlags()

	f.StringVarP(&chdir, "chdir", "C", "", "Working directory")
	f.BoolVarP(&dryRun, "dry-run", "D", false, "Perform a dry-run")
	f.StringSliceVarP(&Env, "env", "e", nil, "Environment variable")
	f.StringVarP(&logLevel, "log-level", "l", "info", "Log severity level")
	f.StringVarP(&profile, "profile", "p", "", "AWS profile to use")

	Command.SetHelpTemplate("\n" + Command.HelpTemplate())
}

// PreRunNoop noop for other commands.
func PreRunNoop(c *cobra.Command, args []string) {
	// TODO: ew... better way to disable in cobra?
}

// PreRun sets up global tasks used for most commands, some use PreRunNoop
// to remove this default behaviour.
func preRun(c *cobra.Command, args []string) error {
	if l, err := log.ParseLevel(logLevel); err == nil {
		log.SetLevel(l)
	}

	config := aws.NewConfig()

	if profile != "" {
		config = config.WithCredentials(credentials.NewSharedCredentials("", profile))
	}

	region, _ := getAWSRegion()
	if region != "" {
		config = config.WithRegion(region)
	}

	Session = session.New(config)

	Project = &project.Project{
		Log:  log.Log,
		Path: ".",
	}

	if dryRun {
		log.SetLevel(log.WarnLevel)
		Project.Service = dryrun.New(Session)
		Project.Concurrency = 1
	} else {
		Project.Service = lambda.New(Session)
	}

	if chdir != "" {
		if err := os.Chdir(chdir); err != nil {
			return err
		}
	}

	return Project.Open()
}

func getAWSRegion() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	cfgPath := fmt.Sprintf("%s/.aws/config", currentUser.HomeDir)
	cfg, err := goconfig.LoadConfigFile(cfgPath)
	if err != nil {
		return "", err
	}

	var sectionName string
	if profile != "" {
		sectionName = fmt.Sprintf("profile %s", profile)
	} else {
		sectionName = "default"
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return "", errors.New("Did not find AWS region from config file")
	}

	return section["region"], nil
}
