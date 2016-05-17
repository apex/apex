package function

import (
	"testing"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/apex/apex/mock"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func init() {
	log.SetHandler(discard.New())
}

var cannedVersions = []*lambda.FunctionConfiguration{
	&lambda.FunctionConfiguration{Version: aws.String("$LATEST")},
	&lambda.FunctionConfiguration{Version: aws.String("1")},
	&lambda.FunctionConfiguration{Version: aws.String("2")},
	&lambda.FunctionConfiguration{Version: aws.String("3")},
	&lambda.FunctionConfiguration{Version: aws.String("4")},
	&lambda.FunctionConfiguration{Version: aws.String("5")},
	&lambda.FunctionConfiguration{Version: aws.String("6")},
	&lambda.FunctionConfiguration{Version: aws.String("7")},
	&lambda.FunctionConfiguration{Version: aws.String("8")},
	&lambda.FunctionConfiguration{Version: aws.String("9")},
	&lambda.FunctionConfiguration{Version: aws.String("10")},
	&lambda.FunctionConfiguration{Version: aws.String("11")},
}

func TestFunction_versionsToCleanup_all(t *testing.T) {
	retainedVersions := 0
	versions, err := versionsToDelete(t, &retainedVersions)

	assert.Len(t, versions, 11)
	assert.Nil(t, err)
}

func TestFunction_versionsToCleanup_default(t *testing.T) {
	var retainedVersions *int
	assert.Nil(t, retainedVersions)

	versions, err := versionsToDelete(t, retainedVersions)

	assert.Len(t, versions, 1)
	assert.Nil(t, err)
}

func TestFunction_versionsToCleanup_two(t *testing.T) {
	retainedVersions := 2
	versions, err := versionsToDelete(t, &retainedVersions)

	assert.Len(t, versions, 9)
	assert.Nil(t, err)
}

func TestFunction_versionsToCleanup_none(t *testing.T) {
	retainedVersions := 11
	versions, err := versionsToDelete(t, &retainedVersions)

	assert.Len(t, versions, 0)
	assert.Nil(t, err)
}

func versionsToDelete(t *testing.T, retainedVersions *int) ([]*lambda.FunctionConfiguration, error) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().ListVersionsByFunction(gomock.Any()).Return(&lambda.ListVersionsByFunctionOutput{Versions: cannedVersions}, nil)

	fn := &Function{
		Config: Config{
			Memory:  128,
			Timeout: 3,
			Role:    "iamrole",
		},
		Path: "_fixtures/nodejsDefaultFile",
		Name: "test",
		Log:  log.Log,
		Service: serviceMock,
	}
	if &retainedVersions != nil {
		fn.RetainedVersions = retainedVersions
	}
	err := fn.Open()
	if err != nil {
		return nil, err
	}

	return fn.versionsToCleanup()

}

