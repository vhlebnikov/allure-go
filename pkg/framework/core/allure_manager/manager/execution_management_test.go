package manager

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/core/constants"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/require"
)

type testMetaMockExecM struct {
	result    *allure.Result
	container *allure.Container
	be        func(t provider.T)
	ae        func(t provider.T)
}

func (m *testMetaMockExecM) GetResult() *allure.Result {
	return m.result
}

func (m *testMetaMockExecM) SetResult(result *allure.Result) {
	m.result = result
}

func (m *testMetaMockExecM) GetContainer() *allure.Container {
	return m.container
}

func (m *testMetaMockExecM) SetBeforeEach(hook func(t provider.T)) {
	m.be = hook
}

func (m *testMetaMockExecM) GetBeforeEach() func(t provider.T) {
	return m.be
}

func (m *testMetaMockExecM) SetAfterEach(hook func(t provider.T)) {
	m.ae = hook
}

func (m *testMetaMockExecM) GetAfterEach() func(t provider.T) {
	return m.ae
}

type suiteMetaMockExecM struct {
	name      string
	container *allure.Container
	hook      func(t provider.T)
}

func (m *suiteMetaMockExecM) GetPackageName() string {
	return m.name
}

func (m *suiteMetaMockExecM) GetRunner() string {
	return m.name
}

func (m *suiteMetaMockExecM) GetSuiteName() string {
	return m.name
}

func (m *suiteMetaMockExecM) GetParentSuite() string {
	return ""
}

func (m *suiteMetaMockExecM) GetSuiteFullName() string {
	return m.name
}

func (m *suiteMetaMockExecM) GetContainer() *allure.Container {
	return m.container
}

func (m *suiteMetaMockExecM) SetBeforeAll(hook func(provider.T)) {
	m.hook = hook
}

func (m *suiteMetaMockExecM) SetAfterAll(hook func(provider.T)) {
	m.hook = hook
}

func (m *suiteMetaMockExecM) GetBeforeAll() func(provider.T) {
	return m.hook
}

func (m *suiteMetaMockExecM) GetAfterAll() func(provider.T) {
	return m.hook
}

func TestAllureManager_AfterAllContext(t *testing.T) {
	manager := allureManager{suiteMeta: &suiteMetaMockExecM{container: allure.NewContainer()}}
	manager.AfterAllContext()
	require.NotNil(t, manager.executionContext)
	require.Equal(t, constants.AfterAllContextName, manager.executionContext.GetName())
}

func TestAllureManager_BeforeAllContext(t *testing.T) {
	manager := allureManager{suiteMeta: &suiteMetaMockExecM{container: allure.NewContainer()}}
	manager.BeforeAllContext()
	require.NotNil(t, manager.executionContext)
	require.Equal(t, constants.BeforeAllContextName, manager.executionContext.GetName())
}

func TestAllureManager_BeforeEachContext(t *testing.T) {
	manager := allureManager{testMeta: &testMetaMockExecM{container: allure.NewContainer()}}
	manager.BeforeEachContext()
	require.NotNil(t, manager.executionContext)
	require.Equal(t, constants.BeforeEachContextName, manager.executionContext.GetName())
}

func TestAllureManager_AfterEachContext(t *testing.T) {
	manager := allureManager{testMeta: &testMetaMockExecM{container: allure.NewContainer()}}
	manager.AfterEachContext()
	require.NotNil(t, manager.executionContext)
	require.Equal(t, constants.AfterEachContextName, manager.executionContext.GetName())
}

func TestAllureManager_TestContext(t *testing.T) {
	manager := allureManager{testMeta: &testMetaMockExecM{result: &allure.Result{}}}
	manager.TestContext()
	require.NotNil(t, manager.executionContext)
	require.Equal(t, constants.TestContextName, manager.executionContext.GetName())
}

// Tests for context creation with test results
func TestAllureManager_AfterEachContextWithResult(t *testing.T) {
	result := allure.NewResult("TestName", "FullTestName")
	result.Status = allure.Passed

	manager := allureManager{
		testMeta: &testMetaMockExecM{
			result:    result,
			container: allure.NewContainer(),
		},
	}

	// Create AfterEach context (should have access to test result)
	manager.AfterEachContext()

	require.NotNil(t, manager.executionContext)
	require.Equal(t, constants.AfterEachContextName, manager.executionContext.GetName())

	// Verify we can get the test result from the context
	retrievedResult := manager.executionContext.GetTestResult()
	require.NotNil(t, retrievedResult)
	require.Equal(t, allure.Passed, retrievedResult.Status)
	require.Equal(t, "TestName", retrievedResult.Name)
}

func TestAllureManager_AfterEachContextWithDifferentStatuses(t *testing.T) {
	tests := []struct {
		name              string
		testName          string
		fullTestName      string
		status            allure.Status
		statusMessage     string
		expectedStatus    allure.Status
		messageContains   string
		exactMessageMatch bool
	}{
		{
			name:              "Failed test result",
			testName:          "FailedTest",
			fullTestName:      "FullFailedTest",
			status:            allure.Failed,
			statusMessage:     "Test assertion failed",
			expectedStatus:    allure.Failed,
			messageContains:   "Test assertion failed",
			exactMessageMatch: true,
		},
		{
			name:              "BeforeEach failure scenario",
			testName:          "TestWillNotRun",
			fullTestName:      "FullTestWillNotRun",
			status:            allure.Failed,
			statusMessage:     "TestName/BeforeEach setup was failed",
			expectedStatus:    allure.Failed,
			messageContains:   "setup was failed",
			exactMessageMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := allure.NewResult(tt.testName, tt.fullTestName)
			result.Status = tt.status
			result.SetStatusMessage(tt.statusMessage)

			manager := allureManager{
				testMeta: &testMetaMockExecM{
					result:    result,
					container: allure.NewContainer(),
				},
			}

			manager.AfterEachContext()

			retrievedResult := manager.executionContext.GetTestResult()
			require.NotNil(t, retrievedResult)
			require.Equal(t, tt.expectedStatus, retrievedResult.Status)

			if tt.exactMessageMatch {
				require.Equal(t, tt.messageContains, retrievedResult.GetStatusMessage())
			} else {
				require.Contains(t, retrievedResult.GetStatusMessage(), tt.messageContains)
			}
		})
	}
}

// Test parametrized test scenario
func TestAllureManager_AfterEachContext_ParametrizedTestScenario(t *testing.T) {
	// Parent test result (parametrized test)
	parentResult := allure.NewResult("TestParametrized", "FullTestParametrized")
	parentResult.Status = allure.Passed // Parent can be Passed even if subtests fail

	manager := allureManager{
		testMeta: &testMetaMockExecM{
			result:    parentResult,
			container: allure.NewContainer(),
		},
	}

	manager.AfterEachContext()

	retrievedResult := manager.executionContext.GetTestResult()
	require.NotNil(t, retrievedResult)
	require.Equal(t, allure.Passed, retrievedResult.Status)
	require.Equal(t, "TestParametrized", retrievedResult.Name)
}
