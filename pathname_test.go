//nolint:godoclint // These are tests
package pathname_test

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mjwhitta/pathname"
	assert "github.com/stretchr/testify/require"
)

func TestBasename(t *testing.T) {
	var path string = filepath.Join("path", "to", "file")

	assert.Equal(t, "file", pathname.Basename(path))
}

func TestDirname(t *testing.T) {
	var dir string = filepath.Join("path", "to")
	var path string = filepath.Join("path", "to", "file")

	assert.Equal(t, dir, pathname.Dirname(path))
}

func TestDoesExist(t *testing.T) {
	type testData struct {
		err      bool
		expected bool
		fn       string
	}

	var e error
	var tests map[string]testData
	var tmp string = t.TempDir()

	if runtime.GOOS == "windows" {
		t.Skip("runtime OS not supported")
	}

	// Create directories
	e = os.MkdirAll(tmp, 0o700)
	assert.NoError(t, e)

	e = os.MkdirAll(filepath.Join(tmp, "noread"), 0o700)
	assert.NoError(t, e)

	// Create files
	//nolint:gosec // G304 false-positive, no user input
	_, e = os.Create(filepath.Join(tmp, "test"))
	assert.NoError(t, e)

	//nolint:gosec // G304 false-positive, no user input
	_, e = os.Create(filepath.Join(tmp, "noread", "test"))
	assert.NoError(t, e)

	// Adjust permissions
	defer func() {
		//nolint:gosec // G302 false-positive, this is a directory
		_ = os.Chmod(filepath.Join(tmp, "noread"), 0o700)
	}()

	e = os.Chmod(filepath.Join(tmp, "noread"), 0o200)
	assert.NoError(t, e)

	tests = map[string]testData{
		"Exists": {
			err:      false,
			expected: true,
			fn:       filepath.Join(tmp, "test"),
		},
		"NoExists": {
			err:      false,
			expected: false,
			fn:       filepath.Join(tmp, "noexist"),
		},
		"ExistsNoRead": {
			err:      true,
			expected: false,
			fn:       filepath.Join(tmp, "noread", "test"),
		},
		"NoExistsNoRead": {
			err:      true,
			expected: false,
			fn:       filepath.Join(tmp, "noread", "noexist"),
		},
	}

	for test, data := range tests {
		t.Run(
			test,
			func(t *testing.T) {
				var e error
				var exists bool

				exists, e = pathname.DoesExist(data.fn)
				if data.err {
					assert.Error(t, e)
					assert.False(t, exists)
				} else {
					assert.NoError(t, e)
					assert.Equal(t, data.expected, exists)
				}
			},
		)
	}
}

func TestExpandPath(t *testing.T) {
	type testData struct {
		in   string
		out  string
		skip bool
	}

	var e error
	var nobody *user.User
	var tests map[string]testData
	var usr *user.User

	nobody, e = user.Lookup("nobody")
	assert.NoError(t, e)
	assert.NotNil(t, nobody)

	usr, e = user.Current()
	assert.NoError(t, e)
	assert.NotNil(t, usr)

	tests = map[string]testData{
		"BadUserDesktop": {in: "~asdf/Desktop", out: "~asdf/Desktop"},
		"BadUserHome":    {in: "~asdf", out: "~asdf"},
		"BinNobody": {
			in:   "~nobody/bin",
			out:  filepath.Join(nobody.HomeDir, "bin"),
			skip: true,
		},
		"DesktopCurrentUser": {
			in:  "~/Desktop",
			out: filepath.Join(usr.HomeDir, "Desktop"),
		},
		"HomeCurrentUser": {in: "~", out: usr.HomeDir},
		"HomeNobody": {
			in:   "~nobody",
			out:  nobody.HomeDir,
			skip: true,
		},
	}

	for test, data := range tests {
		t.Run(
			test,
			func(t *testing.T) {
				var tmp string

				if data.skip && (runtime.GOOS == "windows") {
					t.Skip("runtime OS not supported")
				}

				tmp = pathname.ExpandPath(data.in)
				assert.Equal(t, data.out, tmp)
			},
		)
	}
}
