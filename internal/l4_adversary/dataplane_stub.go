package l4_adversary

// DataPlane abstracts actual offensive execution (agents, infra).
// In PoC, this is a stub that simulates execution.
type DataPlane interface {
	Execute(task *Task) (interface{}, error)
}

// StubDataPlane simulates task execution for PoC.
type StubDataPlane struct{}

// Execute simulates execution and returns a stub result.
func (StubDataPlane) Execute(task *Task) (interface{}, error) {
	return map[string]string{
		"status":     "simulated",
		"task_id":    task.ID,
		"technique":  task.TechniqueID,
		"simulated":  "true",
	}, nil
}
