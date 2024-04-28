package projectpath

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
)

func Root() string {
	return filepath.Join(filepath.Dir(b), "../../")
}
