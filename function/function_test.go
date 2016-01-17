package function

import (
	"errors"
	"testing"

	_ "github.com/apex/apex/runtime/nodejs"

	"github.com/apex/apex/mock"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetHandler(discard.New())
}

func TestFunction_Open_requireConfigValues(t *testing.T) {
	fn := &Function{
		Path: "_fixtures/invalidRuntime",
		Log:  log.Log,
	}
	runtimeErr := fn.Open()

	fn = &Function{
		Path: "_fixtures/invalidMemory",
		Log:  log.Log,
	}
	memoryErr := fn.Open()

	fn = &Function{
		Path: "_fixtures/invalidTimeout",
		Log:  log.Log,
	}
	timeoutErr := fn.Open()

	fn = &Function{
		Path: "_fixtures/invalidRole",
		Log:  log.Log,
	}
	roleErr := fn.Open()

	assert.Contains(t, runtimeErr.Error(), "Runtime: zero value")
	assert.Contains(t, memoryErr.Error(), "Memory: zero value")
	assert.Contains(t, timeoutErr.Error(), "Timeout: zero value")
	assert.Contains(t, roleErr.Error(), "Role: zero value")
}

func TestFunction_Open_detectRuntime(t *testing.T) {
	fn := &Function{
		Config: Config{
			Memory:  128,
			Timeout: 3,
			Role:    "iamrole",
		},
		Path: "_fixtures/nodejsDefaultFile",
		Name: "foo",
		Log:  log.Log,
	}

	assert.Nil(t, fn.Open())
}

func TestFunction_Delete_success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: aws.String("testfn"),
	}).Return(&lambda.DeleteFunctionOutput{}, nil)

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.Delete()

	assert.Nil(t, err)
}

func TestFunction_Delete_failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().DeleteFunction(gomock.Any()).Return(nil, errors.New("API err"))

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.Delete()

	assert.EqualError(t, err, "API err")
}

func TestFunction_Rollback_GetAlias_failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().GetAlias(&lambda.GetAliasInput{
		FunctionName: aws.String("testfn"),
		Name:         aws.String("current"),
	}).Return(nil, errors.New("API err"))

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.Rollback()

	assert.EqualError(t, err, "API err")
}

func TestFunction_Rollback_ListVersions_failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().GetAlias(gomock.Any()).Return(&lambda.AliasConfiguration{FunctionVersion: aws.String("1")}, nil)
	serviceMock.EXPECT().ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String("testfn"),
	}).Return(nil, errors.New("API err"))

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.Rollback()

	assert.EqualError(t, err, "API err")
}

func TestFunction_Rollback_fewVersions(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().GetAlias(gomock.Any()).Return(&lambda.AliasConfiguration{FunctionVersion: aws.String("1")}, nil)
	serviceMock.EXPECT().ListVersionsByFunction(gomock.Any()).Return(&lambda.ListVersionsByFunctionOutput{
		Versions: []*lambda.FunctionConfiguration{&lambda.FunctionConfiguration{Version: aws.String("$LATEST")}},
	}, nil)

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.Rollback()

	assert.EqualError(t, err, "Can't rollback. Only one version deployed.")
}

func TestFunction_Rollback_previousVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	deployedVersions := []*lambda.FunctionConfiguration{
		&lambda.FunctionConfiguration{Version: aws.String("$LATEST")},
		&lambda.FunctionConfiguration{Version: aws.String("1")},
		&lambda.FunctionConfiguration{Version: aws.String("2")},
	}
	currentVersion := aws.String("2")
	afterRollbackVersion := aws.String("1")
	serviceMock.EXPECT().GetAlias(gomock.Any()).Return(&lambda.AliasConfiguration{FunctionVersion: currentVersion}, nil)
	serviceMock.EXPECT().ListVersionsByFunction(gomock.Any()).Return(&lambda.ListVersionsByFunctionOutput{Versions: deployedVersions}, nil)
	serviceMock.EXPECT().UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    aws.String("testfn"),
		Name:            aws.String("current"),
		FunctionVersion: afterRollbackVersion,
	})

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.Rollback()

	assert.Nil(t, err)
}

func TestFunction_Rollback_latestVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	deployedVersions := []*lambda.FunctionConfiguration{
		&lambda.FunctionConfiguration{Version: aws.String("$LATEST")},
		&lambda.FunctionConfiguration{Version: aws.String("1")},
		&lambda.FunctionConfiguration{Version: aws.String("2")},
	}
	currentVersion := aws.String("1")
	afterRollbackVersion := aws.String("2")
	serviceMock.EXPECT().GetAlias(gomock.Any()).Return(&lambda.AliasConfiguration{FunctionVersion: currentVersion}, nil)
	serviceMock.EXPECT().ListVersionsByFunction(gomock.Any()).Return(&lambda.ListVersionsByFunctionOutput{Versions: deployedVersions}, nil)
	serviceMock.EXPECT().UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    aws.String("testfn"),
		Name:            aws.String("current"),
		FunctionVersion: afterRollbackVersion,
	})

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.Rollback()

	assert.Nil(t, err)
}

func TestFunction_Rollback_UpdateAlias_failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().GetAlias(gomock.Any()).Return(&lambda.AliasConfiguration{FunctionVersion: aws.String("1")}, nil)
	serviceMock.EXPECT().ListVersionsByFunction(gomock.Any()).Return(&lambda.ListVersionsByFunctionOutput{
		Versions: []*lambda.FunctionConfiguration{
			&lambda.FunctionConfiguration{Version: aws.String("$LATEST")},
			&lambda.FunctionConfiguration{Version: aws.String("1")},
			&lambda.FunctionConfiguration{Version: aws.String("2")},
		},
	}, nil)
	serviceMock.EXPECT().UpdateAlias(gomock.Any()).Return(nil, errors.New("API err"))

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.Rollback()

	assert.EqualError(t, err, "API err")
}

func TestFunction_RollbackVersion_GetAlias_failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().GetAlias(&lambda.GetAliasInput{
		FunctionName: aws.String("testfn"),
		Name:         aws.String(CurrentAlias),
	}).Return(nil, errors.New("API err"))

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.RollbackVersion("1")

	assert.EqualError(t, err, "API err")
}

func TestFunction_RollbackVersion_sameVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	currentVersion := aws.String("2")
	serviceMock.EXPECT().GetAlias(gomock.Any()).Return(&lambda.AliasConfiguration{FunctionVersion: currentVersion}, nil)

	fn := &Function{
		Service: serviceMock,
		Log:     log.Log,
	}
	err := fn.RollbackVersion("2")

	assert.EqualError(t, err, "Specified version currently deployed.")
}

func TestFunction_RollbackVersion_success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	specifiedVersion := "3"
	currentVersion := "2"
	serviceMock.EXPECT().GetAlias(gomock.Any()).Return(&lambda.AliasConfiguration{FunctionVersion: &currentVersion}, nil)
	serviceMock.EXPECT().UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    aws.String("testfn"),
		Name:            aws.String(CurrentAlias),
		FunctionVersion: &specifiedVersion,
	})

	fn := &Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.RollbackVersion("3")

	assert.Nil(t, err)
}
