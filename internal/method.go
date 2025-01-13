package internal

type Method struct {
	Class      string   `json:"class"`
	Method     string   `json:"method"`
	Attributes []string `json:"attributes"`
	ReturnType string   `json:"returnType"`
	File       string   `json:"file"`
}
