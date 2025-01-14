package internal

import (
	"fmt"
	"os"
	"time"
)

type TestingOutput struct {
	testingResults *TestingResults
}

func NewTestingOutput(testingResults *TestingResults) *TestingOutput {
	return &TestingOutput{testingResults: testingResults}
}

func (t *TestingOutput) Print() {

	summary := Summary{}
	startTime := time.Now()
	items := t.testingResults.All()
	testLocation := ""

	for item := range items {
		for _, result := range item {
			if result.TestLocation != testLocation {
				testLocation = result.TestLocation
				fmt.Printf("\033[1;37;35m ...... \033[0m \033[0;36m %s \033[0m\n", result.TestLocation)
			}

			if result.Exception != "" {
				fmt.Printf(
					"\033[1;37;31m Failed \033[0m   \033[0;91m    %s \033[0m - \033[0;31m %s\033[0m \n",
					"Uncaught Exception",
					result.Exception,
				)
				summary.Failed++
				continue
			}

			fmt.Printf("\033[1;37;35m ...... \033[0m \033[0;95m   %s \033[0m\n", result.RequestUrl)

			for _, expectation := range result.ExpectationResults {
				if expectation.Satisfied {
					summary.Passed++
					fmt.Printf(
						"\033[1;37;92m Passed \033[0m   \033[0;96m    %s \033[0m\n",
						expectation.Label,
					)
				} else {
					summary.Failed++
					fmt.Printf(
						"\033[1;37;31m Failed \033[0m   \033[0;91m    %s \033[0m - \033[0;31m %s\033[0m \n",
						expectation.Label,
						expectation.Message,
					)
				}
			}
		}
	}

	fmt.Printf("\033[1;37;35m ...... \033[0m \033[0;36m %s \033[0m\n", "Summary")
	fmt.Printf("\033[1;37;35m ...... \033[0m \033[0;36m %s \033[0m\n", fmt.Sprintf("Total: %d", summary.Passed+summary.Failed))
	fmt.Printf("\033[1;37;35m ...... \033[0m \033[0;36m %s \033[0m\n", fmt.Sprintf("Passed: %d", summary.Passed))
	fmt.Printf("\033[1;37;35m ...... \033[0m \033[0;36m %s \033[0m\n", fmt.Sprintf("Failed: %d", summary.Failed))
	fmt.Printf(
		"\033[1;37;35m ...... \033[0m \033[0;36m %s \033[0m\n",
		fmt.Sprintf("Time: %.3fs", time.Since(startTime).Seconds()),
	)

	if summary.Failed > 0 {
		fmt.Println("\033[1;37;31m ...... \033[0m \033[1;37;31m Tests failed \033[0m")
		os.Exit(1)
	}
}
