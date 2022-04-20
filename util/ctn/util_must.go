package ctn

// Generic version of must function.
func Must[T any](data T, err error) (res T) {
	if err != nil {
		panic(err)
	}
	return data
}
