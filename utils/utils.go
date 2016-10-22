package utils

func IntArrayRemove(arr []int, i int) []int {
	arr[len(arr)-1] = arr[i]
	arr[i] = arr[len(arr)-1]

	return arr[:len(arr)-1]
}

func IntArrayContains(arr []int, i int) bool {
	for _, b := range arr {
		if b == i {
			return true
		}
	}
	return false
}

func IntArrayRemoveIfContains(arr []int, val int) ([]int, bool) {
	for i, b := range arr {
		if b == val {
			arr[len(arr)-1] = arr[i]
			arr[i] = arr[len(arr)-1]
			arr = arr[:len(arr)-1]
			return arr, true
		}
	}
	return []int{}, false
}
