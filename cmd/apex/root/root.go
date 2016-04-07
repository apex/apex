package root

import (
	"os"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"

	"github.com/apex/apex/dryrun"
	"github.com/apex/apex/project"
	"github.com/apex/apex/utils"
)

// environment for project.
var environment string

// chdir working directory.
var chdir string

// dryRun enabled.
var dryRun bool

// logLevel specified.
var logLevel string

// profile for AWS.
var profile string

// region for AWS.
var region string

// Session instance.
var Session *session.Session

// Project instance.
var Project *project.Project

// Config for AWS.
var Config *aws.Config

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

	f.StringVarP(&environment, "env", "e", "", "Infastructure environment (default: \"defaultEnvironment\" from project.json)")
	f.StringVarP(&chdir, "chdir", "C", "", "Working directory")
	f.BoolVarP(&dryRun, "dry-run", "D", false, "Perform a dry-run")
	f.StringVarP(&logLevel, "log-level", "l", "info", "Log severity level")
	f.StringVarP(&profile, "profile", "p", "", "AWS profile")
	f.StringVarP(&region, "region", "r", "", "AWS region")

	Command.SetHelpTemplate("\n" + Command.HelpTemplate())
}

// PreRunNoop noop for other commands.
func PreRunNoop(c *cobra.Command, args []string) {
	// TODO: ew... better way to disable in cobra?
}

// preRun sets up global tasks used for most commands, some use PreRunNoop
// to remove this default behaviour.
func preRun(c *cobra.Command, args []string) error {
	err := Prepare(c, args)
	if err != nil {
		return err
	}

	return Project.Open()
}

// Prepare handles the global CLI flags and shared functionality without
// the assumption that a Project has already been initialized.
//
// Precedence is currently:
//
//  - flags such as --profile
//  - env vars such as AWS_PROFILE
//  - files such as ~/.aws/config
//
func Prepare(c *cobra.Command, args []string) error {
	if l, err := log.ParseLevel(logLevel); err == nil {
		log.SetLevel(l)
	}

	// config defaults
	Config = aws.NewConfig()

	// profile from flag, env, "default"
	if profile == "" {
		profile = os.Getenv("AWS_PROFILE")
		if profile == "" {
			profile = "default"
		}
	}

	// the default SharedCredentialsProvider checks the env
	os.Setenv("AWS_PROFILE", profile)

	// region from flag, env, file
	if region == "" {
		region = os.Getenv("AWS_REGION")
		if region == "" {
			region, _ = utils.GetRegion(profile)
		}
	}

	if region != "" {
		Config = Config.WithRegion(region)
	}

	Session = session.New(Config)

	Project = &project.Project{
		Environment: environment,
		Log:         log.Log,
		Path:        ".",
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

	return nil
}
