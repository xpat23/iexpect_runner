package internal

import "sync"

type TestingResults struct {
	testMethods *TestMethods
}

func NewTestingResults(testMethods *TestMethods) *TestingResults {
	return &TestingResults{testMethods: testMethods}
}

func (t *TestingResults) All() chan []TestingResult {

	results := make(chan []TestingResult)
	methods, err := t.testMethods.All()

	if err != nil {
		panic("error getting test methods" + err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(len(methods))

	for _, method := range methods {
		go func(m TestMethod) {
			defer wg.Done()
			results <- m.Run()
		}(method)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
