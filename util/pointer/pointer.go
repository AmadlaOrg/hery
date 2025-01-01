package pointer

// ToPtr is a helper function to create a pointer to a string or any other data type
func ToPtr[T any](v T) *T {
	return &v
}
