package sets

func StringSlice(s Set) []string {
	slice := make([]string, 0)
	for _, item := range s.List() {
		v, ok := item.(string)
		if !ok {
			continue
		}
		slice = append(slice, v)
	}
	return slice
}

func IntSlice(s Set) []int {
	slice := make([]int, 0)
	for _, item := range s.List() {
		v, ok := item.(int)
		if !ok {
			continue
		}
		slice = append(slice, v)
	}
	return slice
}
