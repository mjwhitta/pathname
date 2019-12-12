package pathname

import (
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

// Version is the package version.
const Version = "1.0.6"

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
// shortcuts as well as ENV vars.
func ExpandPath(path string) string {
	var e error
	var r = regexp.MustCompile(`^~([^/]+)`)
	var usr *user.User

	// Expand ENV vars
	path = os.ExpandEnv(path)

	// If just ~
	if path == "~" {
		usr, e = user.Current()
		if e != nil {
			return path
		}
		return usr.HomeDir
	}

	// If path starting with ~/
	if strings.HasPrefix(path, "~/") {
		usr, e = user.Current()
		if e != nil {
			return path
		}
		return filepath.Join(usr.HomeDir, path[2:])
	}

	// If ~user shortcut
	for _, match := range r.FindAllStringSubmatch(path, -1) {
		usr, e = user.Lookup(match[1])

		// If user's home directory is found
		if e == nil {
			return r.ReplaceAllString(path, usr.HomeDir)
		}
	}

	// Otherwise just return path
	return path
}
