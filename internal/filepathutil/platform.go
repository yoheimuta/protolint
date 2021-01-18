package filepathutil

import (
	"os"
	"strings"
)

// OSPathSeparator is just a variable that is equal to os.PathSeparator.
// It's set by outside for mainly test usage.
var OSPathSeparator = os.PathSeparator

// IsSameUnixPath compares an unix path to a cross platform path.
//
// The interpretation of the former depends on the platform it runs.
func IsSameUnixPath(unixPath, crossPlatformPath string) bool {
	unixPathSeparator := '/'
	if OSPathSeparator == unixPathSeparator {
		return unixPath == crossPlatformPath
	}
	return strings.Replace(
		unixPath,
		string(unixPathSeparator),
		string(OSPathSeparator),
		-1,
	) == crossPlatformPath
}
