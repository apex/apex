package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type input struct {
	Name    string
	Species string
}

func greet(in *input) (string, error) {
	return fmt.Sprintf("Hello %s, you are a %s", in.Name, in.Species), nil
}

func main() {
	lambda.Start(greet)
}
