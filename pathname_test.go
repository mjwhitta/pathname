package pathname_test

import (
	"crypto/rand"
	"encoding/hex"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"testing"

	assert "github.com/stretchr/testify/require"
	"gitlab.com/mjwhitta/pathname"
)

// This is better than t.TempDir() b/c it's self-contained (so doesn't
// break GitLab runners) and it ensures the permissions are fixed
// before trying to delete during Cleanup().
func tempDir(t *testing.T) string {
	var b []byte = make([]byte, 32)
	var e error
	var out string

	if e = os.MkdirAll("testdata", 0o700); e != nil {
		t.Fatal("could not create testdata directory")
	} else if _, e = rand.Read(b); e != nil {
		t.Fatal("could not create random temp directory")
	}

	out = filepath.Join("testdata", "Test_"+hex.EncodeToString(b))

	t.Cleanup(
		func() {
			// Ensure temp dir can be deleted
			filepath.WalkDir(
				out,
				func(path string, d fs.DirEntry, e error) error {
					if e != nil {
						return e
					}

					if d.IsDir() {
						os.Chmod(path, 0o700)
					} else {
						os.Chmod(path, 0o600)
					}

					return nil
				},
			)

			os.RemoveAll(out)
		},
	)

	return out
}

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
	var tmp string = tempDir(t)

	if runtime.GOOS == "windows" {
		t.Skip("runtime OS not supported")
	}

	// Create directories
	e = os.MkdirAll(tmp, 0o700)
	assert.Nil(t, e)

	e = os.MkdirAll(filepath.Join(tmp, "noread"), 0o700)
	assert.Nil(t, e)

	// Create files
	_, e = os.Create(filepath.Join(tmp, "test"))
	assert.Nil(t, e)

	_, e = os.Create(filepath.Join(tmp, "noread", "test"))
	assert.Nil(t, e)

	// Adjust permissions
	e = os.Chmod(filepath.Join(tmp, "noread"), 0o200)
	assert.Nil(t, e)

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
					assert.NotNil(t, e)
					assert.Equal(t, false, exists)
				} else {
					assert.Nil(t, e)
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
	assert.Nil(t, e)
	assert.NotNil(t, nobody)

	usr, e = user.Current()
	assert.Nil(t, e)
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
