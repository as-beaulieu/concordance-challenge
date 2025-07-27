package filters

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"concordance-challenge/filtered/models"
)

// TestApplyFilters ensures each map in FilterParameters
// sends your word to the right FilterCategory.
func TestApplyFilters(t *testing.T) {
	filters := models.FilterParameters{
		Actions:    map[string]bool{"run": true},
		Frameworks: map[string]bool{"react": true},
		Languages:  map[string]bool{"go": true},
		Leadership: map[string]bool{"lead": true},
		Tools:      map[string]bool{"docker": true},
		Trivial:    map[string]bool{"the": true},
	}

	cases := []struct {
		word string
		want models.FilterCategory
	}{
		{"run", models.FilterActions},
		{"react", models.FilterFrameworks},
		{"go", models.FilterLanguages},
		{"lead", models.FilterLeaderships},
		{"docker", models.FilterTools},
		{"the", models.FilterTrivial},
		{"xyz", models.FilterCategory("")}, // default
	}

	for _, tc := range cases {
		got := ApplyFilters(tc.word, filters)
		if got != tc.want {
			t.Errorf("ApplyFilters(%q) = %q; want %q", tc.word, got, tc.want)
		}
	}
}

// TestFilterWords_Routing writes a single line into FileProcess
// and asserts each token winds up on the correct output channel.
func TestFilterWords_Routing(t *testing.T) {
	// 1) build your process channels
	fileCh := make(chan string)
	p := models.Processes{
		FileProcess:       fileCh,
		ActionsProcess:    make(chan string),
		FrameworksProcess: make(chan string),
		LanguagesProcess:  make(chan string),
		LeadershipProcess: make(chan string),
		OthersProcess:     make(chan string),
		ToolsProcess:      make(chan string),
		TrivialProcess:    make(chan string),
	}
	filters := models.FilterParameters{
		Actions:    map[string]bool{"run": true},
		Frameworks: map[string]bool{"react": true},
		Languages:  map[string]bool{"go": true},
		Leadership: map[string]bool{"lead": true},
		Tools:      map[string]bool{"docker": true},
		Trivial:    map[string]bool{"the": true},
	}

	// 2) kick off the router
	go FilterWords(p, filters)

	// 3) send one mixed‐case line, then close
	fileCh <- "Run React GO LEAD Docker the unknown"
	close(fileCh)

	// 4) helper to drain a channel into a slice
	drain := func(ch chan string) []string {
		var got []string
		for w := range ch {
			got = append(got, w)
		}
		return got
	}

	tests := []struct {
		name string
		ch   chan string
		want []string
	}{
		{"Actions", p.ActionsProcess, []string{"run"}},
		{"Frameworks", p.FrameworksProcess, []string{"react"}},
		{"Languages", p.LanguagesProcess, []string{"go"}},
		{"Leadership", p.LeadershipProcess, []string{"lead"}},
		{"Tools", p.ToolsProcess, []string{"docker"}},
		{"Trivial", p.TrivialProcess, []string{"the"}},
		{"Others", p.OthersProcess, []string{"unknown"}},
	}

	for _, tc := range tests {
		if got := drain(tc.ch); !equal(got, tc.want) {
			t.Errorf("%s channel = %v; want %v", tc.name, got, tc.want)
		}
	}
}

// equal is a tiny helper for comparing string slices.
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestReadFilters writes a temporary filters.json and
// validates that ReadFilters() populates all maps correctly.
func TestReadFilters(t *testing.T) {
	// snapshot cwd, switch into a temp dir
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig)

	// craft a FiltersInput and write filters.json
	input := models.FiltersInput{
		Actions:    []string{"run"},
		Frameworks: []string{"react"},
		Languages:  []string{"go"},
		Leadership: []string{"lead"},
		Tools:      []string{"docker"},
		Trivial:    []string{"the"},
	}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "filters.json"), data, 0644); err != nil {
		t.Fatal(err)
	}

	// call ReadFilters and inspect the result
	got, err := ReadFilters()
	if err != nil {
		t.Fatalf("ReadFilters error: %v", err)
	}

	checkMap := func(name string, m map[string]bool, want []string) {
		for _, k := range want {
			if !m[k] {
				t.Errorf("%s missing key %q", name, k)
			}
		}
		if len(m) != len(want) {
			t.Errorf("%s has extra keys %v; want precisely %v", name, m, want)
		}
	}

	checkMap("Actions", got.Actions, input.Actions)
	checkMap("Frameworks", got.Frameworks, input.Frameworks)
	checkMap("Languages", got.Languages, input.Languages)
	checkMap("Leadership", got.Leadership, input.Leadership)
	checkMap("Tools", got.Tools, input.Tools)
	checkMap("Trivial", got.Trivial, input.Trivial)
}
