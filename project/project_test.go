package project_test

import (
	"testing"

	"github.com/apex/apex/mock/service"
	"github.com/apex/apex/project"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	_ "github.com/apex/apex/plugins/golang"
	_ "github.com/apex/apex/plugins/hooks"
	_ "github.com/apex/apex/plugins/inference"
	_ "github.com/apex/apex/plugins/nodejs"
	_ "github.com/apex/apex/plugins/python"
	_ "github.com/apex/apex/plugins/shim"
)

func init() {
	log.SetHandler(discard.New())
}

func TestProject_Open_requireConfigValues(t *testing.T) {
	p := &project.Project{
		Path: "_fixtures/invalidName",
		Log:  log.Log,
	}
	nameErr := p.Open()

	assert.Contains(t, nameErr.Error(), "Name: zero value")
}

func TestProject_LoadFunctions_loadAll(t *testing.T) {
	p := &project.Project{
		Path:            "_fixtures/twoFunctions",
		Log:             log.Log,
		ServiceProvider: mock_service.NewMockProvideriface(nil),
	}

	assert.NoError(t, p.Open(), "open")
	assert.NoError(t, p.LoadFunctions(), "load")

	assert.Equal(t, 2, len(p.Functions))
	assert.Equal(t, "bar", p.Functions[0].Name)
	assert.Equal(t, "foo", p.Functions[1].Name)
}

func TestProject_LoadFunctions_loadSpecified(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockProvider := mock_service.NewMockProvideriface(mockCtrl)
	mockProvider.EXPECT().NewService(nil)

	p := &project.Project{
		Path:            "_fixtures/twoFunctions",
		Log:             log.Log,
		ServiceProvider: mockProvider,
	}

	assert.NoError(t, p.Open(), "open")
	assert.NoError(t, p.LoadFunctions("foo"), "load")

	assert.Equal(t, 1, len(p.Functions))
	assert.Equal(t, "foo", p.Functions[0].Name)
}

func TestProject_LoadFunctions_onlyExisting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockProvider := mock_service.NewMockProvideriface(mockCtrl)
	mockProvider.EXPECT().NewService(nil)

	p := &project.Project{
		Path:            "_fixtures/twoFunctions",
		Log:             log.Log,
		ServiceProvider: mockProvider,
	}

	assert.NoError(t, p.Open(), "open")
	assert.NoError(t, p.LoadFunctions("foo", "something"), "load")

	assert.Equal(t, 1, len(p.Functions))
	assert.Equal(t, "foo", p.Functions[0].Name)
}

func TestProject_LoadFunctions_noFunctionLoaded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockProvider := mock_service.NewMockProvideriface(mockCtrl)
	mockProvider.EXPECT().NewService(nil).MaxTimes(0)

	p := &project.Project{
		Path:            "_fixtures/twoFunctions",
		Log:             log.Log,
		ServiceProvider: mockProvider,
	}

	p.Open()
	err := p.LoadFunctions("something")

	assert.EqualError(t, err, "no function loaded")
}

func TestProject_LoadFunctionByPath_mergeEnvWithFunctionEnv(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockProvider := mock_service.NewMockProvideriface(mockCtrl)
	mockProvider.EXPECT().NewService(nil)

	p := &project.Project{
		Path:            "_fixtures/envMerge",
		Log:             log.Log,
		ServiceProvider: mockProvider,
	}

	assert.NoError(t, p.Open(), "open")
	assert.NoError(t, p.LoadFunctions("foo"), "load")

	assert.Equal(t, map[string]string{"PROJECT_ENV": "projectEnv", "FUNCTION_ENV": "functionEnv", "APEX_FUNCTION_NAME": "foo", "LAMBDA_FUNCTION_NAME": "envMerge_foo"}, p.Functions[0].Environment)
}

func TestProject_LoadFunctionByPath_overrideVpcWithFunctionVpc(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockProvider := mock_service.NewMockProvideriface(mockCtrl)
	mockProvider.EXPECT().NewService(nil).MaxTimes(2)

	p := &project.Project{
		Path:            "_fixtures/vpcOverride",
		Log:             log.Log,
		ServiceProvider: mockProvider,
	}

	p.Open()

	assert.Equal(t, "sg-default", p.VPC.SecurityGroups[0])

	bar, _ := p.LoadFunction("bar")
	assert.Equal(t, "sg-override", bar.VPC.SecurityGroups[0])
	assert.Equal(t, "sg-default", p.VPC.SecurityGroups[0])

	foo, _ := p.LoadFunction("foo")
	assert.Equal(t, "sg-default", foo.VPC.SecurityGroups[0])
	assert.Equal(t, "sg-default", p.VPC.SecurityGroups[0])
}

func TestProject_LoadFunctionByPath_edgeFunction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockProvider := mock_service.NewMockProvideriface(mockCtrl)
	mockProvider.EXPECT().NewService(aws.NewConfig().WithRegion("us-east-1"))

	p := &project.Project{
		Path:            "_fixtures/edgeFunction",
		Log:             log.Log,
		ServiceProvider: mockProvider,
	}

	assert.NoError(t, p.Open(), "open")
	assert.NoError(t, p.LoadFunctions("foo"), "load")
}

func TestProject_LoadFunctionByPath_multiRegionFunction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockProvider := mock_service.NewMockProvideriface(mockCtrl)
	gomock.InOrder(
		mockProvider.EXPECT().NewService(aws.NewConfig().WithRegion("eu-west-2")),
		mockProvider.EXPECT().NewService(nil),
		mockProvider.EXPECT().NewService(aws.NewConfig().WithRegion("ap-northeast-1")),
	)

	p := &project.Project{
		Path:            "_fixtures/multiRegionFunctions",
		Log:             log.Log,
		ServiceProvider: mockProvider,
	}

	assert.NoError(t, p.Open(), "open")
	assert.NoError(t, p.LoadFunctions(), "load")

	assert.Equal(t, 3, len(p.Functions))
	assert.Equal(t, "bar", p.Functions[0].Name)
	assert.Equal(t, "baz", p.Functions[1].Name)
	assert.Equal(t, "foo", p.Functions[2].Name)
}
