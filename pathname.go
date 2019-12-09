package pathname

import (
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

// Version is the package version.
const Version = "1.0.5"

// Basename wraps filepath.Base(path str).
func Basename(path string) string {
	return filepath.Base(ExpandPath(path))
}

// Dirname wraps filepath.Dir(path str).
func Dirname(path string) string {
	return filepath.Dir(ExpandPath(path))
}

// DoesExist returns true if the specified path exists on disk, false
// otherwise.
func DoesExist(path string) bool {
	if _, err := os.Stat(ExpandPath(path)); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		panic(err)
	}
}

// ExpandPath will expand the specified path accounting for ~ or ~user
// shortcuts.
func ExpandPath(path string) string {
	var usr, e = user.Current()
	if e != nil {
		return path
	}

	var home = usr.HomeDir

	// If just ~
	if path == "~" {
		return home
	}

	// If path starting with ~/
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}

	// If ~user shortcut
	var r = regexp.MustCompile(`^~([^/]+)`)
	var matches = r.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		usr, e = user.Lookup(match[1])

		// If user's home directory is found
		if e == nil {
			home = usr.HomeDir
			return r.ReplaceAllString(path, home)
		}
	}

	// Otherwise just return path
	return path
}
