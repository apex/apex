package vpc

// VPC represents part of function configuration used to define security groups and subnets for Lambda function.
type VPC struct {
	SecurityGroups []string `json:"securityGroups"`
	Subnets        []string `json:"subnets"`
}
