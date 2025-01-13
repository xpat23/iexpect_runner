package internal

type ExpectationResult struct {
	Satisfied bool   `json:"satisfied"`
	Label     string `json:"label"`
	Message   string `json:"message"`
}
