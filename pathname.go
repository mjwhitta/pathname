package pathname

import (
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

const Version = "1.0.0"

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
