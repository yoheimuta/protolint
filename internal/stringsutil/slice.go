package stringsutil

import "github.com/yoheimuta/protolint/internal/filepathutil"

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

// ContainsCrossPlatformPathInSlice searches the unix path haystack for the cross-platform needle.
func ContainsCrossPlatformPathInSlice(
	needle string,
	unixPathHaystack []string,
) bool {
	for _, h := range unixPathHaystack {
		if filepathutil.IsSameUnixPath(h, needle) {
			return true
		}
	}
	return false
}
