package stringsutil

// ContainsStringInSlice searches the haystack for the needle.
func ContainsStringInSlice(
	needle string,
	haystack []string,
) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}
