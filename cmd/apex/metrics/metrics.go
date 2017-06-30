// Package metrics outputs metrics for a function.
package metrics

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/dustin/go-humanize"
	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/colors"
	"github.com/apex/apex/cost"
	"github.com/apex/apex/metrics"
)

// duration of results.
var duration time.Duration

// example output.
const example = `
    Print the last 24 hours of metrics for all functions
    $ apex metrics

    Print the last 24 hours of metrics for a function
    $ apex metrics foo

    Print metrics for a function with a specified start time, e.g. the last 3 days
    $ apex metrics foo --since 72h`

// Command config.
var Command = &cobra.Command{
	Use:     "metrics [<name>...] [<duration>]",
	Short:   "Output function metrics",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.DurationVarP(&duration, "since", "s", 24*time.Hour, "Start time of the results")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Project.LoadFunctions(args...); err != nil {
		return err
	}

	service := lambda.New(root.Session)

	config := metrics.Config{
		Service:   cloudwatch.New(root.Session),
		StartDate: time.Now().UTC().Add(-duration),
		EndDate:   time.Now().UTC(),
	}

	m := metrics.Metrics{
		Config: config,
	}

	for _, fn := range root.Project.Functions {
		m.FunctionNames = append(m.FunctionNames, fn.FunctionName)
	}

	aggregated := m.Collect()

	fmt.Println()
	for _, fn := range root.Project.Functions {
		m := aggregated[fn.FunctionName]

		conf, err := service.GetFunctionConfiguration(&lambda.GetFunctionConfigurationInput{FunctionName: &fn.FunctionName})
		if err != nil {
			return err
		}

		memory := int(*conf.MemorySize)
		costTotal := humanize.FormatFloat("", cost.Cost(m.Invocations, m.Duration, memory))
		costDuration := humanize.FormatFloat("", cost.DurationCost(m.Duration, memory))
		costInvocations := humanize.FormatFloat("", cost.RequestCost(m.Invocations))

		fmt.Printf("  \033[%dm%s\033[0m\n", colors.Blue, fn.Name)
		fmt.Printf("    total cost: $%s\n", costTotal)
		fmt.Printf("    invocations: %s ($%s)\n", humanize.Comma(int64(m.Invocations)), costInvocations)
		fmt.Printf("    duration: %s ($%s)\n", time.Millisecond*time.Duration(m.Duration), costDuration)
		fmt.Printf("    throttles: %v\n", m.Throttles)
		fmt.Printf("    errors: %s\n", humanize.Comma(int64(m.Errors)))
		fmt.Printf("    memory: %d\n", memory)
		fmt.Println()
	}

	return nil
}
