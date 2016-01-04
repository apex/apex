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

	assert.Error(t, err, "API err")
}
