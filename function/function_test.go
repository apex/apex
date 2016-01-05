package function_test

import (
	"errors"
	"testing"

	"github.com/apex/apex/function"
	"github.com/apex/apex/mock"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFunction_Delete_success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	serviceMock.EXPECT().DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: aws.String("testfn"),
	}).Return(&lambda.DeleteFunctionOutput{}, nil)

	fn := &function.Function{
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
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
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
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

	fn := &function.Function{
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
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
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
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
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
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
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
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
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
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
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
	}
	err := fn.Rollback()

	assert.EqualError(t, err, "API err")
}

func TestFunction_Rollback_sameVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	currentVersion := aws.String("2")
	serviceMock.EXPECT().GetAlias(gomock.Any()).Return(&lambda.AliasConfiguration{FunctionVersion: currentVersion}, nil)

	fn := &function.Function{
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
	}
	err := fn.Rollback("2")

	assert.EqualError(t, err, "Specified version currently deployed.")
}

func TestFunction_Rollback_specifiedVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	deployedVersions := []*lambda.FunctionConfiguration{
		&lambda.FunctionConfiguration{Version: aws.String("$LATEST")},
		&lambda.FunctionConfiguration{Version: aws.String("1")},
		&lambda.FunctionConfiguration{Version: aws.String("2")},
		&lambda.FunctionConfiguration{Version: aws.String("3")},
	}
	currentVersion := aws.String("3")
	afterRollbackVersion := aws.String("1")
	serviceMock.EXPECT().GetAlias(gomock.Any()).Return(&lambda.AliasConfiguration{FunctionVersion: currentVersion}, nil)
	serviceMock.EXPECT().ListVersionsByFunction(gomock.Any()).Return(&lambda.ListVersionsByFunctionOutput{Versions: deployedVersions}, nil)
	serviceMock.EXPECT().UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    aws.String("testfn"),
		Name:            aws.String("current"),
		FunctionVersion: afterRollbackVersion,
	})

	fn := &function.Function{
		Config:  function.Config{Name: "testfn"},
		Service: serviceMock,
	}
	err := fn.Rollback("1")

	assert.Nil(t, err)
}
