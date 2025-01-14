package internal

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type TestMethod struct {
	Class      string   `json:"class"`
	Method     string   `json:"method"`
	Attributes []string `json:"attributes"`
	ReturnType string   `json:"returnType"`
	File       string   `json:"file"`
}

func (tm *TestMethod) Run() []TestingResult {

	path := fmt.Sprintf("%s::%s", tm.File, tm.Method)

	cmd := exec.Command("php", "run_tests", path)

	stdOut, err := cmd.Output()

	if err != nil {
		panic("error running test" + err.Error())
	}

	var results []TestingResult

	jsonMappingErr := json.Unmarshal(stdOut, &results)

	if jsonMappingErr != nil {
		return []TestingResult{{Exception: jsonMappingErr.Error()}}
	}

	return results
}
