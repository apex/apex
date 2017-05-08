package infra

import (
	"fmt"
	"strings"

	"github.com/apex/apex/function"
	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// ConfigRelation is a structure to store a combination of a function and it's configuration
type ConfigRelation struct {
	Function      *function.Function
	Configuration *lambda.FunctionConfiguration
}

// functionVars returns the function variables as terraform -var arguments.
func (p *Proxy) functionVars() (args []string) {

	args = append(args, "-var")
	args = append(args, fmt.Sprintf("aws_region=%s", p.Region))

	args = append(args, "-var")
	args = append(args, fmt.Sprintf("apex_environment=%s", p.Environment))

	if p.Role != "" {
		args = append(args, "-var")
		args = append(args, fmt.Sprintf("apex_function_role=%s", p.Role))
	}

	// GetConfig is slow, so store the results and then use the configurations as needed
	var relations []ConfigRelation
	for _, fn := range p.Functions {
		config, err := fn.GetConfig()
		if err != nil {
			log.Debugf("can't fetch function config: %s", err.Error())
			continue
		}

		var relation ConfigRelation
		relation.Function = fn
		relation.Configuration = config.Configuration
		relations = append(relations, relation)
	}

	args = append(args, getFunctionArnVars(relations)...)
	args = append(args, getFunctionNameVars(relations)...)
	args = append(args, getFunctionArns(relations)...)
	args = append(args, getFunctionNames(relations)...)

	return args
}

// Generate a series of variables apex_function_FUNCTION that contains the arn of the functions
// This function is being phased out in favour of apex_function_arns
func getFunctionArnVars(relations []ConfigRelation) (args []string) {
	log.Debugf("Generating the tfvar apex_function_FUNCTION")
	for _, rel := range relations {
		args = append(args, "-var")
		args = append(args, fmt.Sprintf("apex_function_%s=%s", rel.Function.Name, *rel.Configuration.FunctionArn))
	}
	return args
}

// Generates a series of variable apex_function_FUNCTION_name that contains the full name of the functions
// This function is being phased out in favour of apex_function_names
func getFunctionNameVars(relations []ConfigRelation) (args []string) {
	log.Debugf("Generating the tfvar apex_function_FUNCTION_name")
	for _, rel := range relations {
		args = append(args, "-var")
		args = append(args, fmt.Sprintf("apex_function_%s_name=%s", rel.Function.Name, *rel.Configuration.FunctionArn))
	}
	return args
}

// Generates a map that has the function's name as a key and the arn of the function as a value
func getFunctionArns(relations []ConfigRelation) (args []string) {
	log.Debugf("Generating the tfvar apex_function_arns")
	var arns []string
	for _, rel := range relations {
		arns = append(arns, fmt.Sprintf("%s = %q", rel.Function.Name, *rel.Configuration.FunctionArn))
	}
	args = append(args, "-var")
	args = append(args, fmt.Sprintf("apex_function_arns={ %s }", strings.Join(arns, ", ")))
	return args
}

// Generates a map that has the function's name as a key and the full name of the function as a value
func getFunctionNames(relations []ConfigRelation) (args []string) {
	log.Debugf("Generating the tfvar apex_function_names")
	var names []string
	for _, rel := range relations {
		names = append(names, fmt.Sprintf("%s = %q", rel.Function.Name, rel.Function.FunctionName))
	}
	args = append(args, "-var")
	args = append(args, fmt.Sprintf("apex_function_names={ %s }", strings.Join(names, ", ")))
	return args
}
