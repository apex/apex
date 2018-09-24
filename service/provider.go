// Package service implements the Lambda service provider.
package service

import (
	"github.com/apex/apex/dryrun"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

// Provideriface is service factory
type Provideriface interface {
	NewService(cfg *aws.Config) lambdaiface.LambdaAPI
}

// Provider implements interface
type Provider struct {
	Session *session.Session
	DryRun  bool
}

// NewProvider with session and dry run
func NewProvider(session *session.Session, dryRun bool) Provideriface {
	return &Provider{
		Session: session,
		DryRun:  dryRun,
	}
}

// NewService returns Lambda service with AWS config
func (p *Provider) NewService(cfg *aws.Config) lambdaiface.LambdaAPI {
	if p.DryRun {
		return dryrun.New(p.Session)
	} else if cfg != nil {
		return lambda.New(p.Session, cfg)
	} else {
		return lambda.New(p.Session)
	}
}
