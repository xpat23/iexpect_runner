package main

import (
	"encoding/json"
	"fmt"
	"iexpect_go/internal"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {

	path := os.Args[1]
	cmd := exec.Command("php", "get_tests", path)

	tests, err := cmd.Output()

	if err != nil {
		panic(err.Error())
	}

	var methods []internal.Method

	jsonMappingErr := json.Unmarshal(tests, &methods)

	if jsonMappingErr != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	summary := internal.Summary{}
	wg.Add(len(methods))
	startTime := time.Now()

	for _, method := range methods {
		go runCommand(
			fmt.Sprintf("%s::%s", method.File, method.Method),
			&wg,
			&mutex,
			&summary,
		)
	}

	wg.Wait()

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

func runCommand(path string, wg *sync.WaitGroup, mutex *sync.Mutex, summary *internal.Summary) {
	defer wg.Done()
	cmd := exec.Command("php", "run_tests", path)

	stdOut, err := cmd.Output()

	if err != nil {
		panic("error running test" + err.Error())
	}

	var results []internal.TestingResult

	jsonMappingErr := json.Unmarshal(stdOut, &results)

	if jsonMappingErr != nil {
		fmt.Println("Error unmarshalling JSON:", jsonMappingErr.Error())
		return
	}

	Print(results, mutex, summary)
}

func Print(results []internal.TestingResult, mutex *sync.Mutex, summary *internal.Summary) {
	testLocation := ""
	//startTime := time.Now()

	for _, result := range results {
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
			mutex.Lock()
			summary.Failed++
			mutex.Unlock()
			continue
		}

		fmt.Printf("\033[1;37;35m ...... \033[0m \033[0;95m   %s \033[0m\n", result.RequestUrl)

		for _, expectation := range result.ExpectationResults {
			if expectation.Satisfied {
				mutex.Lock()
				summary.Passed++
				mutex.Unlock()
				fmt.Printf(
					"\033[1;37;92m Passed \033[0m   \033[0;96m    %s \033[0m\n",
					expectation.Label,
				)
			} else {
				mutex.Lock()
				summary.Failed++
				mutex.Unlock()
				fmt.Printf(
					"\033[1;37;31m Failed \033[0m   \033[0;91m    %s \033[0m - \033[0;31m %s\033[0m \n",
					expectation.Label,
					expectation.Message,
				)
			}
		}
	}
}
