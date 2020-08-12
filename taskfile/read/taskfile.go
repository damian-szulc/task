package read

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-task/task/v2/taskfile"

	"gopkg.in/yaml.v2"
)

var (
	// ErrIncludedTaskfilesCantHaveIncludes is returned when a included Taskfile contains includes
	ErrIncludedTaskfilesCantHaveIncludes = errors.New("task: Included Taskfiles can't have includes. Please, move the include to the main Taskfile")
)

// Taskfile reads a Taskfile for a given directory
func Taskfile(dir string, entrypoint string) (*taskfile.Taskfile, error) {
	path := filepath.Join(dir, entrypoint)
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf(`task: No Taskfile found on "%s". Use "task --init" to create a new one`, path)
	}
	t, err := readTaskfile(path)
	if err != nil {
		return nil, err
	}

	for namespace, path := range t.Includes {
		path = filepath.Join(dir, path)
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			path = filepath.Join(path, "Taskfile.yml")
		}
		includedTaskfile, err := readTaskfile(path)
		if err != nil {
			return nil, err
		}
		if len(includedTaskfile.Includes) > 0 {
			return nil, ErrIncludedTaskfilesCantHaveIncludes
		}
		if err = taskfile.Merge(t, includedTaskfile, namespace); err != nil {
			return nil, err
		}
	}

	path = filepath.Join(dir, fmt.Sprintf("Taskfile_%s.yml", runtime.GOOS))
	if _, err = os.Stat(path); err == nil {
		osTaskfile, err := readTaskfile(path)
		if err != nil {
			return nil, err
		}
		if err = taskfile.Merge(t, osTaskfile); err != nil {
			return nil, err
		}
	}

	for name, task := range t.Tasks {
		task.Task = name
	}

	return t, nil
}

func readTaskfile(file string) (*taskfile.Taskfile, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var t taskfile.Taskfile
	return &t, yaml.NewDecoder(f).Decode(&t)
}
