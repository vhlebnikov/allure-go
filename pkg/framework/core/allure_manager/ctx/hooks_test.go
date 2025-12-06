package ctx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/core/constants"
)

func TestNewAfterAllCtx(t *testing.T) {
	ctx := NewAfterAllCtx(allure.NewContainer())
	require.NotNil(t, ctx)
}

func TestNewAfterEachCtx(t *testing.T) {
	ctx := NewAfterEachCtx(allure.NewContainer())
	require.NotNil(t, ctx)
}

func TestNewBeforeEachCtx(t *testing.T) {
	ctx := NewBeforeEachCtx(allure.NewContainer())
	require.NotNil(t, ctx)
}

func TestNewBeforeAllCtx(t *testing.T) {
	ctx := NewBeforeAllCtx(allure.NewContainer())
	require.NotNil(t, ctx)
}

func TestHooksCtx_GetName(t *testing.T) {
	th := hooksCtx{name: "test"}
	require.Equal(t, "test", th.GetName())
}

func TestHooksCtx_AddStep(t *testing.T) {
	testStep := allure.NewSimpleStep("test")
	beforeEach := hooksCtx{name: constants.BeforeEachContextName, container: &allure.Container{}}
	beforeEach.AddStep(testStep)
	require.NotEmpty(t, beforeEach.container.Befores)
	require.Len(t, beforeEach.container.Befores, 1)
	require.Equal(t, testStep, beforeEach.container.Befores[0])

	beforeAll := hooksCtx{name: constants.BeforeAllContextName, container: &allure.Container{}}
	beforeAll.AddStep(testStep)
	require.NotEmpty(t, beforeAll.container.Befores)
	require.Len(t, beforeAll.container.Befores, 1)
	require.Equal(t, testStep, beforeAll.container.Befores[0])

	afterEach := hooksCtx{name: constants.AfterEachContextName, container: &allure.Container{}}
	afterEach.AddStep(testStep)
	require.NotEmpty(t, afterEach.container.Afters)
	require.Len(t, afterEach.container.Afters, 1)
	require.Equal(t, testStep, afterEach.container.Afters[0])

	afterAll := hooksCtx{name: constants.AfterAllContextName, container: &allure.Container{}}
	afterAll.AddStep(testStep)
	require.NotEmpty(t, afterAll.container.Afters)
	require.Len(t, afterAll.container.Afters, 1)
	require.Equal(t, testStep, afterAll.container.Afters[0])
}

func TestHooksCtx_AddAttachment(t *testing.T) {
	attach := allure.NewAttachment("testAttach", allure.Text, []byte("test"))
	beforeAll := hooksCtx{name: constants.BeforeAllContextName, container: &allure.Container{}}
	beforeAll.AddAttachments(attach)
	require.NotEmpty(t, beforeAll.container.Befores)
	require.Len(t, beforeAll.container.Befores, 1)
	require.NotEmpty(t, beforeAll.container.Befores[0].Attachments)
	require.Len(t, beforeAll.container.Befores[0].Attachments, 1)
	require.Equal(t, attach, beforeAll.container.Befores[0].Attachments[0])

	beforeEach := hooksCtx{name: constants.BeforeEachContextName, container: &allure.Container{}}
	beforeEach.AddAttachments(attach)
	require.NotEmpty(t, beforeEach.container.Befores)
	require.Len(t, beforeEach.container.Befores, 1)
	require.NotEmpty(t, beforeEach.container.Befores[0].Attachments)
	require.Len(t, beforeEach.container.Befores[0].Attachments, 1)
	require.Equal(t, attach, beforeEach.container.Befores[0].Attachments[0])

	afterAll := hooksCtx{name: constants.AfterAllContextName, container: &allure.Container{}}
	afterAll.AddAttachments(attach)
	require.NotEmpty(t, afterAll.container.Afters)
	require.Len(t, afterAll.container.Afters, 1)
	require.NotEmpty(t, afterAll.container.Afters[0].Attachments)
	require.Len(t, afterAll.container.Afters[0].Attachments, 1)
	require.Equal(t, attach, afterAll.container.Afters[0].Attachments[0])

	afterEach := hooksCtx{name: constants.AfterEachContextName, container: &allure.Container{}}
	afterEach.AddAttachments(attach)
	require.NotEmpty(t, afterEach.container.Afters)
	require.Len(t, afterEach.container.Afters, 1)
	require.NotEmpty(t, afterEach.container.Afters[0].Attachments)
	require.Len(t, afterEach.container.Afters[0].Attachments, 1)
	require.Equal(t, attach, afterEach.container.Afters[0].Attachments[0])
}

// Tests for GetTestResult functionality
func TestNewAfterEachCtxWithResult(t *testing.T) {
	container := allure.NewContainer()
	result := allure.NewResult("TestName", "TestFullName")
	result.Status = allure.Passed

	ctx := NewAfterEachCtxWithResult(container, result)
	require.NotNil(t, ctx)

	// Verify we can get the result back
	retrievedResult := ctx.GetTestResult()
	require.NotNil(t, retrievedResult)
	require.Equal(t, result, retrievedResult)
	require.Equal(t, allure.Passed, retrievedResult.Status)
	require.Equal(t, "TestName", retrievedResult.Name)
}

func TestHooksCtx_GetTestResult_AfterEach(t *testing.T) {
	container := allure.NewContainer()
	testResult := allure.NewResult("TestName", "TestFullName")
	testResult.Status = allure.Failed
	testResult.SetStatusMessage("Test failed")

	ctx := NewAfterEachCtxWithResult(container, testResult)

	result := ctx.GetTestResult()
	require.NotNil(t, result)
	require.Equal(t, allure.Failed, result.Status)
	require.Equal(t, "TestName", result.Name)
	require.Equal(t, "Test failed", result.GetStatusMessage())
}

func TestHooksCtx_GetTestResult_NilForOtherContexts(t *testing.T) {
	container := allure.NewContainer()

	// BeforeEach should return nil
	beforeEachCtx := NewBeforeEachCtx(container)
	require.Nil(t, beforeEachCtx.GetTestResult())

	// BeforeAll should return nil
	beforeAllCtx := NewBeforeAllCtx(container)
	require.Nil(t, beforeAllCtx.GetTestResult())

	// AfterEach without result should return nil
	afterEachCtx := NewAfterEachCtx(container)
	require.Nil(t, afterEachCtx.GetTestResult())
}

// Tests for parametrized test scenarios
func TestHooksCtx_GetTestResult_ParametrizedTestScenario(t *testing.T) {
	// Simulate a parametrized test where the parent test is Passed
	// even though subtests may have failed
	container := allure.NewContainer()
	parentResult := allure.NewResult("TestParametrized", "FullTestParametrized")
	parentResult.Status = allure.Passed // Parent can be Passed even if subtests fail

	ctx := NewAfterEachCtxWithResult(container, parentResult)

	result := ctx.GetTestResult()
	require.NotNil(t, result)
	require.Equal(t, allure.Passed, result.Status)
	require.Equal(t, "TestParametrized", result.Name)
}

// Test for nested test scenario
func TestHooksCtx_GetTestResult_NestedTestScenario(t *testing.T) {
	// Simulate nested tests where parent test status may differ from subtests
	container := allure.NewContainer()

	// Parent test can be Passed even if nested subtests failed
	parentResult := allure.NewResult("TestNested", "FullTestNested")
	parentResult.Status = allure.Passed

	ctx := NewAfterEachCtxWithResult(container, parentResult)

	result := ctx.GetTestResult()
	require.NotNil(t, result)
	require.Equal(t, allure.Passed, result.Status)
}

// Test for BeforeEach failure scenario
func TestHooksCtx_GetTestResult_BeforeEachFailure(t *testing.T) {
	container := allure.NewContainer()

	// When BeforeEach fails, test doesn't run but result is created with Failed status
	result := allure.NewResult("TestWillNotRun", "FullTestWillNotRun")
	result.Status = allure.Failed
	result.SetStatusMessage("TestWillNotRun/BeforeEach setup was failed")

	ctx := NewAfterEachCtxWithResult(container, result)

	retrievedResult := ctx.GetTestResult()
	require.NotNil(t, retrievedResult)
	require.Equal(t, allure.Failed, retrievedResult.Status)
	require.Contains(t, retrievedResult.GetStatusMessage(), "setup was failed")
}

// Test concurrent access safety
func TestHooksCtx_GetTestResult_ConcurrentAccess(t *testing.T) {
	container := allure.NewContainer()
	result := allure.NewResult("TestConcurrent", "FullTestConcurrent")
	result.Status = allure.Passed

	ctx := NewAfterEachCtxWithResult(container, result)

	// Multiple goroutines reading the result
	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			r := ctx.GetTestResult()
			require.NotNil(t, r)
			require.Equal(t, allure.Passed, r.Status)
		}()
	}

	wg.Wait()
}
