package string

import "testing"

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		sortedList []string
		target     string
		expected   int
	}{
		{[]string{"apple", "banana", "cherry", "date", "fig", "grape"}, "cherry", 2},
		{[]string{"apple", "banana", "cherry", "date", "fig", "grape"}, "banana", 1},
		{[]string{"apple", "banana", "cherry", "date", "fig", "grape"}, "grape", 5},
		{[]string{"apple", "banana", "cherry", "date", "fig", "grape"}, "apple", 0},
		{[]string{"apple", "banana", "cherry", "date", "fig", "grape"}, "fig", 4},
		{[]string{"apple", "banana", "cherry", "date", "fig", "grape"}, "orange", -1},
	}

	for _, test := range tests {
		t.Run(test.target, func(t *testing.T) {
			result := BinarySearch(test.sortedList, test.target)
			if result != test.expected {
				t.Errorf("BinarySearch(%v, %s) = %d; expected %d", test.sortedList, test.target, result, test.expected)
			}
		})
	}
}
