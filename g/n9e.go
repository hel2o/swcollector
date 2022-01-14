package g

func Array_include(array_a []string, array_b []string) bool { //b include a
	for _, v := range array_a {
		if In_array(v, array_b) {
			continue
		} else {
			return false
		}
	}
	return true
}

func In_array(a string, array []string) bool {
	for _, v := range array {
		if a == v {
			return true
		}
	}
	return false
}
