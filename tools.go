//go:build tools
// +build tools

package tools

import (
	_ "github.com/gordonklaus/ineffassign"
	_ "github.com/kisielk/errcheck"
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "honnef.co/go/tools/simple"
	_ "honnef.co/go/tools/unused"
)
