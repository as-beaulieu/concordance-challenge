package sorting

import (
	"reflect"
	"testing"

	"concordance-challenge/filtered/models"
)

func TestMapToSortedSlice(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]int
		want []models.WordCount
	}{
		{
			name: "empty map",
			in:   map[string]int{},
			want: []models.WordCount{},
		},
		{
			name: "single element",
			in:   map[string]int{"alpha": 1},
			want: []models.WordCount{
				{Word: "alpha", Count: 1},
			},
		},
		{
			name: "distinct counts",
			in: map[string]int{
				"low":    1,
				"medium": 5,
				"high":   10,
			},
			want: []models.WordCount{
				{Word: "high", Count: 10},
				{Word: "medium", Count: 5},
				{Word: "low", Count: 1},
			},
		},
		{
			name: "tie on count, lex order",
			in: map[string]int{
				"apple":  3,
				"banana": 3,
				"cherry": 2,
				"date":   2,
			},
			want: []models.WordCount{
				// 3s descending, tie broken alphabetically
				{Word: "apple", Count: 3},
				{Word: "banana", Count: 3},
				// then 2s, tie broken alphabetically
				{Word: "cherry", Count: 2},
				{Word: "date", Count: 2},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := MapToSortedSlice(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("MapToSortedSlice(%v) =\n  got:  %#v\n  want: %#v",
					tc.in, got, tc.want)
			}
		})
	}
}
