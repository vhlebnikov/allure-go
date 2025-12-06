package ctx

import (
	"sync"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/core/constants"
	"github.com/stretchr/testify/require"
)

func TestNewTestCtx(t *testing.T) {
	ctx := NewTestCtx(&allure.Result{})
	require.NotNil(t, ctx)
}

func TestTestCtx_GetName(t *testing.T) {
	th := testCtx{name: "test"}
	require.Equal(t, "test", th.GetName())
}

func TestTestCtx_AddStep(t *testing.T) {
	testStep := allure.NewSimpleStep("test")
	test := testCtx{name: constants.TestContextName, result: &allure.Result{}}
	test.AddStep(testStep)
	require.NotEmpty(t, test.result.Steps)
	require.Len(t, test.result.Steps, 1)
	require.Equal(t, testStep, test.result.Steps[0])
}

func TestTestCtx_AddAttachment(t *testing.T) {
	attach := allure.NewAttachment("testAttach", allure.Text, []byte("test"))
	test := testCtx{name: constants.TestContextName, result: &allure.Result{}}
	test.AddAttachments(attach)
	require.NotEmpty(t, test.result.Attachments)
	require.Len(t, test.result.Attachments, 1)
	require.Equal(t, attach, test.result.Attachments[0])
}

// Tests for GetTestResult functionality in test context
func TestTestCtx_GetTestResult(t *testing.T) {
	result := allure.NewResult("TestName", "TestFullName")
	result.Status = allure.Passed

	testCtx := NewTestCtx(result)

	retrievedResult := testCtx.GetTestResult()
	require.NotNil(t, retrievedResult)
	require.Equal(t, result, retrievedResult)
	require.Equal(t, allure.Passed, retrievedResult.Status)
	require.Equal(t, "TestName", retrievedResult.Name)
}

func TestTestCtx_GetTestResult_WithDifferentStatuses(t *testing.T) {
	testCases := []struct {
		name   string
		status allure.Status
	}{
		{"Passed", allure.Passed},
		{"Failed", allure.Failed},
		{"Broken", allure.Broken},
		{"Skipped", allure.Skipped},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := allure.NewResult("Test", "FullTest")
			result.Status = tc.status

			testCtx := NewTestCtx(result)

			retrievedResult := testCtx.GetTestResult()
			require.NotNil(t, retrievedResult)
			require.Equal(t, tc.status, retrievedResult.Status)
		})
	}
}

func TestTestCtx_GetTestResult_NilResult(t *testing.T) {
	testCtx := testCtx{name: constants.TestContextName, result: nil}

	retrievedResult := testCtx.GetTestResult()
	require.Nil(t, retrievedResult)
}

func TestTestCtx_GetName_TestContext(t *testing.T) {
	result := allure.NewResult("TestName", "TestFullName")
	testCtx := NewTestCtx(result)

	name := testCtx.GetName()
	require.Equal(t, constants.TestContextName, name)
}

// Additional tests for parametrized and nested test scenarios
func TestTestCtx_GetTestResult_ParametrizedTest(t *testing.T) {
	// Simulate parametrized test result
	result := allure.NewResult("TestParametrized", "suite.TestParametrized")
	result.Status = allure.Passed

	testCtx := NewTestCtx(result)

	retrievedResult := testCtx.GetTestResult()
	require.NotNil(t, retrievedResult)
	require.Equal(t, allure.Passed, retrievedResult.Status)
	require.Equal(t, "TestParametrized", retrievedResult.Name)
}

func TestTestCtx_GetTestResult_NestedTest(t *testing.T) {
	// Simulate nested test result
	result := allure.NewResult("TestNested", "suite.TestNested")
	result.Status = allure.Passed

	testCtx := NewTestCtx(result)

	retrievedResult := testCtx.GetTestResult()
	require.NotNil(t, retrievedResult)
	require.Equal(t, allure.Passed, retrievedResult.Status)
}

func TestTestCtx_GetTestResult_WithStatusDetails(t *testing.T) {
	result := allure.NewResult("TestWithDetails", "suite.TestWithDetails")
	result.Status = allure.Failed
	result.SetStatusMessage("Assertion failed: expected 1, got 2")
	result.SetStatusTrace("trace stack here")

	testCtx := NewTestCtx(result)

	retrievedResult := testCtx.GetTestResult()
	require.NotNil(t, retrievedResult)
	require.Equal(t, allure.Failed, retrievedResult.Status)
	require.Equal(t, "Assertion failed: expected 1, got 2", retrievedResult.GetStatusMessage())
	require.Equal(t, "trace stack here", retrievedResult.GetStatusTrace())
}

// Test concurrent access to result
func TestTestCtx_GetTestResult_ConcurrentAccess(t *testing.T) {
	result := allure.NewResult("ConcurrentTest", "suite.ConcurrentTest")
	result.Status = allure.Passed

	testCtx := NewTestCtx(result)

	const numGoroutines = 20
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			r := testCtx.GetTestResult()
			require.NotNil(t, r)
			require.Equal(t, allure.Passed, r.Status)
		}()
	}

	wg.Wait()
}
