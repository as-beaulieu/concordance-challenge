package models

type WordCount struct {
	Word  string
	Count int
}

type FinalWordCount struct {
	Actions    []WordCount `json:"actions"`
	Frameworks []WordCount `json:"frameworks"`
	Langauges  []WordCount `json:"langauges"`
	Leadership []WordCount `json:"leadership"`
	Others     []WordCount `json:"others"`
	Tools      []WordCount `json:"tools"`
	Trivial    []WordCount `json:"trivial"`
}
