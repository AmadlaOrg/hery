package string

import "sort"

// BinarySearch checks if a string is in a sorted list of strings using binary search
// and returns the index of the target string. If the target is not found, it returns -1.
func BinarySearch(sortedList []string, target string) int {
	index := sort.Search(len(sortedList), func(i int) bool {
		return sortedList[i] >= target
	})
	if index < len(sortedList) && sortedList[index] == target {
		return index
	}
	return -1
}

// ExistInStringArr checks if a string exist in an array
func ExistInStringArr(str string, strArr []string) bool {
	for _, strItem := range strArr {
		if strItem == str {
			return true
		}
	}
	return false
}
