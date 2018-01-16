package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// Input of the function.
type Input struct {
	Name string
}

// greet function.
func greet(in *Input) (string, error) {
	return fmt.Sprintf("Hello %s", in.Name), nil
}

func main() {
	lambda.Start(greet)
}
