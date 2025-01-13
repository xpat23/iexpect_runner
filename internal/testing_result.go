package internal

type TestingResult struct {
	TestLocation       string              `json:"testLocation"`
	Exception          string              `json:"exception"`
	Response           any                 `json:"response"`
	RequestUrl         string              `json:"requestUrl"`
	ExpectationResults []ExpectationResult `json:"expectationResults"`
}
