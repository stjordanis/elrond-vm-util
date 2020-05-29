package testcontroller

import ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"

// ScenarioExecutor describes a component that can run a VM scenario.
type ScenarioExecutor interface {
	// ExecuteScenario executes the scenario and checks if it passed. Failure is signaled by returning an error.
	// The FileResolver helps with resolving external steps.
	// TODO: group into a "execution context" param.
	ExecuteScenario(*ij.Scenario, ij.FileResolver) error
}

// ScenarioRunner is a component that can run json scenarios, using a provided executor.
type ScenarioRunner struct {
	Executor ScenarioExecutor
	Parser   ij.Parser
}

// NewScenarioRunner creates new ScenarioRunner instance.
func NewScenarioRunner(executor ScenarioExecutor, fileResolver ij.FileResolver) *ScenarioRunner {
	return &ScenarioRunner{
		Executor: executor,
		Parser: ij.Parser{
			FileResolver: fileResolver,
		},
	}
}
