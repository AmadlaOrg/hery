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

func TestExistInStringArr(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		strArr   []string
		expected bool
	}{
		{
			name:     "String exists in array",
			str:      "hello",
			strArr:   []string{"hello", "world"},
			expected: true,
		},
		{
			name:     "String does not exist in array",
			str:      "goodbye",
			strArr:   []string{"hello", "world"},
			expected: false,
		},
		{
			name:     "Empty string in array",
			str:      "",
			strArr:   []string{"", "world"},
			expected: true,
		},
		{
			name:     "Empty string not in array",
			str:      "",
			strArr:   []string{"hello", "world"},
			expected: false,
		},
		{
			name:     "Array is empty",
			str:      "hello",
			strArr:   []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExistInStringArr(tt.str, tt.strArr); got != tt.expected {
				t.Errorf("ExistInStringArr() = %v, want %v", got, tt.expected)
			}
		})
	}
}
