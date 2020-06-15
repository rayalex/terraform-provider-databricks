package databricks

import (
	"errors"
	"fmt"
	"testing"

	"github.com/databrickslabs/databricks-terraform/client/model"
	"github.com/stretchr/testify/assert"
)

var executeMock func(clusterID, langauge, commandStr string) (model.Command, error)

type commandExecutorMock struct{}

func (a commandExecutorMock) Execute(clusterID, langauge, commandStr string) (model.Command, error) {
	return executeMock(clusterID, langauge, commandStr)
}

func TestAzureBlobMountReadRetrievesMountInformation(t *testing.T) {
	const cn = "mycontainer"
	const sacc = "mystorage"
	const dir = "mydirectory"

	testCases := []struct {
		ExpectedResult string
		ExpectedError  error
		CommandResult  *model.CommandResults
	}{
		{
			ExpectedResult: fmt.Sprintf("abfss://%s@%s.dfs.core.windows.net/%s", cn, sacc, dir),
			CommandResult: &model.CommandResults{
				ResultType: "text",
				Data:       fmt.Sprintf("abfss://%s@%s.dfs.core.windows.net/%s", cn, sacc, dir),
			},
		},
		{
			ExpectedError: errors.New("unable to find mount point"),
			CommandResult: &model.CommandResults{
				ResultType: "text",
				Data:       "",
			},
		},
		{
			ExpectedError: fmt.Errorf("does not match uri with storage account and container values %s@%s != abfss://x@y.dfs.core.windows.net/z!", cn, sacc),
			CommandResult: &model.CommandResults{
				ResultType: "text",
				Data:       "abfss://x@y.dfs.core.windows.net/z",
			},
		},
		{
			ExpectedError: errors.New("out of wibble error"),
			CommandResult: &model.CommandResults{
				ResultType: "error",
				Summary:    "out of wibble error",
			},
		},
	}

	for _, tc := range testCases {
		executeMock = func(clusterID, langauge, commandStr string) (model.Command, error) {
			return model.Command{
				Results: tc.CommandResult,
			}, nil
		}
		executorMock := commandExecutorMock{}

		blobMount := NewAzureBlobMount(cn, sacc, dir, "mount", "", "", "")

		result, err := blobMount.Read(executorMock, "wibble")

		assert.Equal(t, tc.ExpectedResult, result)
		assert.Equal(t, tc.ExpectedError, err)
	}
}

func TestAzureADLSGen1MountReadRetrievesMountInformation(t *testing.T) {
	const sacc = "mystorage"
	const dir = "mydirectory"

	testCases := []struct {
		ExpectedResult string
		ExpectedError  error
		CommandResult  *model.CommandResults
	}{
		{
			ExpectedResult: fmt.Sprintf("adl://%s.azuredatalakestore.net/%s", sacc, dir),
			CommandResult: &model.CommandResults{
				ResultType: "text",
				Data:       fmt.Sprintf("adl://%s.azuredatalakestore.net/%s", sacc, dir),
			},
		},
		{
			ExpectedError: errors.New("unable to find mount point"),
			CommandResult: &model.CommandResults{
				ResultType: "text",
				Data:       "",
			},
		},
		{
			ExpectedError: fmt.Errorf("does not match uri with storage account and container values %s@%s != adl://x.azuredatalakestore.net/z!", sacc, dir),
			CommandResult: &model.CommandResults{
				ResultType: "text",
				Data:       "adl://x.azuredatalakestore.net/z",
			},
		},
		{
			ExpectedError: errors.New("out of wibble error"),
			CommandResult: &model.CommandResults{
				ResultType: "error",
				Summary:    "out of wibble error",
			},
		},
	}

	for _, tc := range testCases {
		executeMock = func(clusterID, langauge, commandStr string) (model.Command, error) {
			return model.Command{
				Results: tc.CommandResult,
			}, nil
		}
		executorMock := commandExecutorMock{}

		adlsGen2Mount := NewAzureADLSGen1Mount(sacc, dir, "mount", "", "", "", "", "")

		result, err := adlsGen2Mount.Read(executorMock, "wibble")

		assert.Equal(t, tc.ExpectedResult, result)
		assert.Equal(t, tc.ExpectedError, err)
	}
}

func TestAzureADLSGen2MountReadRetrievesMountInformation(t *testing.T) {
	const cn = "mycontainer"
	const sacc = "mystorage"
	const dir = "mydirectory"

	testCases := []struct {
		ExpectedResult string
		ExpectedError  error
		CommandResult  *model.CommandResults
	}{
		{
			ExpectedResult: fmt.Sprintf("abfss://%s@%s.dfs.core.windows.net/%s", cn, sacc, dir),
			CommandResult: &model.CommandResults{
				ResultType: "text",
				Data:       fmt.Sprintf("abfss://%s@%s.dfs.core.windows.net/%s", cn, sacc, dir),
			},
		},
		{
			ExpectedError: errors.New("unable to find mount point"),
			CommandResult: &model.CommandResults{
				ResultType: "text",
				Data:       "",
			},
		},
		{
			ExpectedError: fmt.Errorf("does not match uri with storage account and container values %s@%s != abfss://x@y.dfs.core.windows.net/z!", cn, sacc),
			CommandResult: &model.CommandResults{
				ResultType: "text",
				Data:       "abfss://x@y.dfs.core.windows.net/z",
			},
		},
		{
			ExpectedError: errors.New("out of wibble error"),
			CommandResult: &model.CommandResults{
				ResultType: "error",
				Summary:    "out of wibble error",
			},
		},
	}

	for _, tc := range testCases {
		executeMock = func(clusterID, langauge, commandStr string) (model.Command, error) {
			return model.Command{
				Results: tc.CommandResult,
			}, nil
		}
		executorMock := commandExecutorMock{}

		adlsGen2Mount := NewAzureADLSGen2Mount(cn, sacc, dir, "mount", "", "", "", "", true)

		result, err := adlsGen2Mount.Read(executorMock, "wibble")

		assert.Equal(t, tc.ExpectedResult, result)
		assert.Equal(t, tc.ExpectedError, err)
	}
}

func TestProcessAzureWasbAbfssUrisCorrectlySplitsURI(t *testing.T) {
	testCases := []struct {
		URI                string
		ExpectedContainer  string
		ExpectedStorageAcc string
		ExpectedDirectory  string
	}{
		{
			URI:                "abfss://wibble@mystorage.dfs.core.windows.net/wobble",
			ExpectedContainer:  "wibble",
			ExpectedStorageAcc: "mystorage",
			ExpectedDirectory:  "/wobble",
		},
		{
			URI:                "abfss://wibble@mystorage.dfs.core.windows.net",
			ExpectedContainer:  "wibble",
			ExpectedStorageAcc: "mystorage",
			ExpectedDirectory:  "",
		},
	}

	for _, tc := range testCases {
		container, storageAcc, dir, err := ProcessAzureWasbAbfssUris(tc.URI)
		assert.Equal(t, tc.ExpectedContainer, container)
		assert.Equal(t, tc.ExpectedStorageAcc, storageAcc)
		assert.Equal(t, tc.ExpectedDirectory, dir)
		assert.Nil(t, err)
	}
}
