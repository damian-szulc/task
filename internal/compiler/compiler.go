package compiler

import (
	"github.com/damian-szulc/task/v2/taskfile"
)

// Compiler handles compilation of a task before its execution.
// E.g. variable merger, template processing, etc.
type Compiler interface {
	GetVariables(t *taskfile.Task, call taskfile.Call) (taskfile.Vars, error)
	HandleDynamicVar(v taskfile.Var) (string, error)
}
