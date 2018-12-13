package function_test

import (
	"errors"
	"testing"

	_ "github.com/apex/apex/plugins/golang"
	_ "github.com/apex/apex/plugins/hooks"
	_ "github.com/apex/apex/plugins/inference"
	_ "github.com/apex/apex/plugins/nodejs"
	_ "github.com/apex/apex/plugins/python"
	_ "github.com/apex/apex/plugins/ruby"
	_ "github.com/apex/apex/plugins/shim"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/apex/apex/function"
	"github.com/apex/apex/mock"
	"github.com/apex/apex/utils"
)

func init() {
	log.SetHandler(discard.New())
}

func TestFunction_Open_requireConfigValues(t *testing.T) {
	fn := &function.Function{
		Path: "_fixtures/invalidRuntime",
		Log:  log.Log,
	}
	runtimeErr := fn.Open("")
	assert.Contains(t, runtimeErr.Error(), "Runtime: zero value")

	fn = &function.Function{
		Path: "_fixtures/invalidMemory",
		Log:  log.Log,
	}
	memoryErr := fn.Open("")
	assert.Contains(t, memoryErr.Error(), "Memory: zero value")

	fn = &function.Function{
		Path: "_fixtures/invalidTimeout",
		Log:  log.Log,
	}
	timeoutErr := fn.Open("")
	assert.Contains(t, timeoutErr.Error(), "Timeout: zero value")

	fn = &function.Function{
		Path: "_fixtures/invalidRole",
		Log:  log.Log,
	}
	roleErr := fn.Open("")
	assert.Contains(t, roleErr.Error(), "Role: zero value")

}

func TestFunction_Open_detectRuntime(t *testing.T) {
	fn := &function.Function{
		Config: function.Config{
			Memory:  128,
			Timeout: 3,
			Role:    "iamrole",
		},
		Path: "_fixtures/nodejsDefaultFile",
		Name: "foo",
		Log:  log.Log,
	}

	assert.Nil(t, fn.Open(""))
}

func TestFunction_Create_edgeFunction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)
	zip := []byte("abcdef")
	fnName := "testedgefn"
	desc := "description"
	handler := "handle"
	kmskeyarn := "kmskeyarn"
	runtime := "node"
	memorysize := int64(128)
	timeout := int64(3)
	role := "testrole"
	updatedVersion := "1"
	fnAlias := "current"
	retainedVersions := 1

	serviceMock.EXPECT().CreateFunction(&lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ZipFile: zip,
		},
		FunctionName: &fnName,
		Description:  &desc,
		Environment: &lambda.Environment{
			Variables: make(map[string]*string),
		}, // Edge does not support environment
		Publish:    aws.Bool(true),
		Handler:    &handler,
		KMSKeyArn:  &kmskeyarn,
		MemorySize: &memorysize,
		Role:       &role,
		Runtime:    &runtime,
		Timeout:    &timeout,
		VpcConfig: &lambda.VpcConfig{
			SecurityGroupIds: []*string{},
			SubnetIds:        []*string{},
		},
	}).Return(&lambda.FunctionConfiguration{
		Version: &updatedVersion,
	}, nil)
	serviceMock.EXPECT().CreateAlias(&lambda.CreateAliasInput{
		FunctionName:    &fnName,
		FunctionVersion: &updatedVersion,
		Name:            &fnAlias,
	}).Return(&lambda.AliasConfiguration{}, nil)

	fn := &function.Function{
		FunctionName: fnName,
		Service:      serviceMock,
		Log:          log.Log,
		Alias:        fnAlias,
		Config: function.Config{
			Description:      desc,
			Runtime:          runtime,
			Memory:           memorysize,
			Timeout:          timeout,
			Role:             role,
			Handler:          handler,
			KMSKeyArn:        kmskeyarn,
			Edge:             true,
			RetainedVersions: &retainedVersions,
			Environment:      map[string]string{"FOO": "foovalue"},
		},
	}

	err := fn.Create(zip)

	assert.Nil(t, err)
}

func TestFunction_Delete_success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: aws.String("testfn"),
	}).Return(&lambda.DeleteFunctionOutput{}, nil)

	fn := &function.Function{
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

	fn := &function.Function{
		FunctionName: "testfn",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.Delete()

	assert.EqualError(t, err, "API err")
}

func TestFunction_DeployCode_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)
	code := []byte("something")
	codeSha256 := utils.Sha256([]byte("something else"))
	functionOutput := &lambda.GetFunctionOutput{
		Configuration: &lambda.FunctionConfiguration{
			CodeSha256: &codeSha256,
		},
	}
	fnName := "testfn"
	updatedVersion := "1"
	fnAlias := "current"
	retainedVersions := 1

	serviceMock.EXPECT().UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: &fnName,
		Publish:      aws.Bool(true),
		ZipFile:      code,
	}).Return(&lambda.FunctionConfiguration{
		Version: &updatedVersion,
	}, nil)
	serviceMock.EXPECT().CreateAlias(&lambda.CreateAliasInput{
		FunctionName:    &fnName,
		FunctionVersion: &updatedVersion,
		Name:            &fnAlias,
	}).Return(&lambda.AliasConfiguration{}, nil)
	serviceMock.EXPECT().ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
		FunctionName: &fnName,
	}).Return(&lambda.ListVersionsByFunctionOutput{
		Versions: []*lambda.FunctionConfiguration{
			&lambda.FunctionConfiguration{},
		},
	}, nil)

	fn := &function.Function{
		FunctionName: fnName,
		Service:      serviceMock,
		Log:          log.Log,
		Alias:        fnAlias,
		Config: function.Config{
			RetainedVersions: &retainedVersions,
			Environment:      make(map[string]string),
		},
	}

	err := fn.DeployCode(code, functionOutput)

	assert.Nil(t, err)
}

func TestFunction_DeployCode_UnchangedUpdatesAlias(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)
	code := []byte("something")
	codeSha256 := utils.Sha256(code)
	currentFnVersion := "$LATEST"
	functionOutput := &lambda.GetFunctionOutput{
		Configuration: &lambda.FunctionConfiguration{
			Version:    &currentFnVersion,
			CodeSha256: &codeSha256,
		},
	}
	fnName := "testfn"
	fnAlias := "current"
	retainedVersions := 1

	versions := []string{"1", "2", "3"}

	serviceMock.EXPECT().ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
		FunctionName: &fnName,
	}).Return(&lambda.ListVersionsByFunctionOutput{
		Versions: []*lambda.FunctionConfiguration{
			&lambda.FunctionConfiguration{},
			&lambda.FunctionConfiguration{
				Version: &versions[0],
			},
			&lambda.FunctionConfiguration{
				Version: &versions[1],
			},
			&lambda.FunctionConfiguration{
				Version: &versions[2],
			},
		},
	}, nil)
	serviceMock.EXPECT().CreateAlias(&lambda.CreateAliasInput{
		FunctionName:    &fnName,
		FunctionVersion: &versions[2],
		Name:            &fnAlias,
	}).Return(nil, awserr.New("ResourceConflictException", "message", errors.New("message")))
	serviceMock.EXPECT().UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    &fnName,
		FunctionVersion: &versions[2],
		Name:            &fnAlias,
	}).Return(&lambda.AliasConfiguration{}, nil)

	fn := &function.Function{
		FunctionName: fnName,
		Service:      serviceMock,
		Log:          log.Log,
		Alias:        fnAlias,
		Config: function.Config{
			RetainedVersions: &retainedVersions,
			Environment:      make(map[string]string),
		},
	}

	err := fn.DeployCode(code, functionOutput)

	assert.Nil(t, err)
}

func TestFunction_Rollback_GetAlias_failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().GetAlias(&lambda.GetAliasInput{
		FunctionName: aws.String("testfn"),
		Name:         aws.String("current"),
	}).Return(nil, errors.New("API err"))

	fn := &function.Function{
		FunctionName: "testfn",
		Alias:        "current",
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

	fn := &function.Function{
		FunctionName: "testfn",
		Alias:        "current",
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

	fn := &function.Function{
		FunctionName: "testfn",
		Alias:        "current",
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

	fn := &function.Function{
		FunctionName: "testfn",
		Alias:        "current",
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

	fn := &function.Function{
		FunctionName: "testfn",
		Alias:        "current",
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

	fn := &function.Function{
		FunctionName: "testfn",
		Alias:        "current",
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
		Name:         aws.String(function.CurrentAlias),
	}).Return(nil, errors.New("API err"))

	fn := &function.Function{
		FunctionName: "testfn",
		Alias:        "current",
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

	fn := &function.Function{
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
		Name:            aws.String(function.CurrentAlias),
		FunctionVersion: &specifiedVersion,
	})

	fn := &function.Function{
		FunctionName: "testfn",
		Alias:        "current",
		Service:      serviceMock,
		Log:          log.Log,
	}
	err := fn.RollbackVersion("3")

	assert.Nil(t, err)
}
