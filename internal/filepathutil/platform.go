package filepathutil

import (
	"os"
	"strings"
)

const unixPathSeparator = '/'

// OSPathSeparator is just a variable that is equal to os.PathSeparator.
// It's set by outside for mainly test usage.
var OSPathSeparator = os.PathSeparator

// IsSameUnixPath compares an unix path to a cross platform path.
//
// The interpretation of the unix path depends on the platform it runs.
func IsSameUnixPath(unixPath, crossPlatformPath string) bool {
	if OSPathSeparator == unixPathSeparator {
		return unixPath == crossPlatformPath
	}
	return convertToOSPath(unixPath) == crossPlatformPath
}

// HasUnixPathPrefix checks whether a cross platform path has an unix path as its prefix.
//
// The interpretation of the unix path depends on the platform it runs.
func HasUnixPathPrefix(crossPlatformPath, unixPath string) bool {
	if OSPathSeparator == unixPathSeparator {
		return strings.HasPrefix(crossPlatformPath, unixPath)
	}
	return strings.HasPrefix(
		crossPlatformPath,
		convertToOSPath(unixPath),
	)
}

func convertToOSPath(unixPath string) string {
	return strings.Replace(
		unixPath,
		string(unixPathSeparator),
		osPathSeparator(),
		-1,
	)
}

func osPathSeparator() string {
	return string(OSPathSeparator)
}
