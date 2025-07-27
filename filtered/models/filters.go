package models

var (
	FilterActions     FilterCategory = "actions"
	FilterFrameworks  FilterCategory = "frameworks"
	FilterLanguages   FilterCategory = "language"
	FilterLeaderships FilterCategory = "leadership"
	FilterTools       FilterCategory = "tools"
	FilterTrivial     FilterCategory = "trivial"
)

type FilterCategory string

type FiltersInput struct {
	Actions    []string `json:"actions"`
	Frameworks []string `json:"frameworks"`
	Languages  []string `json:"languages"`
	Leadership []string `json:"leadership"`
	Tools      []string `json:"tools"`
	Trivial    []string `json:"trivial"`
}

type FilterParameters struct {
	Actions    map[string]bool
	Frameworks map[string]bool
	Languages  map[string]bool
	Leadership map[string]bool
	Tools      map[string]bool
	Trivial    map[string]bool
}

type Processes struct {
	FileProcess       chan string
	ActionsProcess    chan string
	FrameworksProcess chan string
	LanguagesProcess  chan string
	LeadershipProcess chan string
	OthersProcess     chan string
	ToolsProcess      chan string
	TrivialProcess    chan string
}
