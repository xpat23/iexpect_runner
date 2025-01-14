package internal

import (
	"encoding/json"
	"os/exec"
)

type TestMethods struct {
	path string
}

func NewTestMethods(path string) *TestMethods {
	return &TestMethods{path: path}
}

func (t *TestMethods) All() ([]TestMethod, error) {
	cmd := exec.Command("php", "get_tests", t.path)

	tests, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var methods []TestMethod

	jsonMappingErr := json.Unmarshal(tests, &methods)

	if jsonMappingErr != nil {
		return nil, jsonMappingErr
	}

	return methods, nil
}
