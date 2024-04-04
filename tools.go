//go:build tools
// +build tools

package repairmanagementbackend

import (
	_ "github.com/go-faster/errors"
	_ "go.uber.org/zap"
	_ "golang.org/x/exp/constraints"
	_ "golang.org/x/exp/maps"
	_ "golang.org/x/exp/slices"
	_ "golang.org/x/tools/imports"
)
