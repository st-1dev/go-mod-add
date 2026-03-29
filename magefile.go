//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// --------------------------------------------------------------------------------------------------------------------
// BUILD section.

// Build compiles the application.
func Build() error {
	args := []string{"build", `-ldflags=-s -w`, "./..."}
	return sh.RunWithV(map[string]string{"CGO_ENABLED": "0"}, "go", args...)
}

// Format formats Go source code.
func Format() error {
	return sh.RunV("gofmt", "-l", "-w", ".")
}
