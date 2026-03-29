package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

func main() {
	exitCode, err := run(os.Args[0], os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode)
}

func run(program string, args []string) (exitCode int, err error) {
	if len(args) == 0 {
		return 1, fmt.Errorf("usage: %s <module@version> [<module@version> ...]", program)
	}

	deps := make([]module.Version, 0, len(args))
	for _, arg := range args {
		path, version, found := strings.Cut(arg, "@")
		if !found {
			return 1, fmt.Errorf("invalid dependency format %s", arg)
		}
		deps = append(deps, module.Version{Path: path, Version: version})
	}
	slices.Reverse(deps)

	// For each dependency: remove from go.mod
	for _, dep := range deps {
		// Remove the entry from go.mod
		err = removeFromGoMod(dep.Path)
		if err != nil {
			return 1, fmt.Errorf("failed to remove %s from go.mod: %w", dep.Path, err)
		}
	}

	// For each dependency: run go get -u
	for _, dep := range deps {
		// go get -u module@version
		err = runCommand("go", "get", "-u", dep.String())
		if err != nil {
			return 1, fmt.Errorf("go get failed for %s: %w", dep.String(), err)
		}
	}

	// go mod tidy
	err = runCommand("go", "mod", "tidy")
	if err != nil {
		return 1, fmt.Errorf("go mod tidy failed: %w", err)
	}

	// go mod vendor — only if the vendor directory exists
	if dirExists("vendor") {
		err = runCommand("go", "mod", "vendor")
		if err != nil {
			return 1, fmt.Errorf("go mod vendor failed: %w", err)
		}
	}

	return 0, nil
}

// removeFromGoMod removes the module's require entry from go.mod.
func removeFromGoMod(modulePath string) (err error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("could not read go.mod: %w", err)
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return fmt.Errorf("could not parse go.mod: %w", err)
	}

	err = f.DropRequire(modulePath)
	if err != nil {
		return fmt.Errorf("could not drop require %s: %w", modulePath, err)
	}

	f.Cleanup()

	out, err := f.Format()
	if err != nil {
		return fmt.Errorf("could not format go.mod: %w", err)
	}

	err = os.WriteFile("go.mod", out, 0644)
	if err != nil {
		return fmt.Errorf("could not write go.mod: %w", err)
	}

	return nil
}

// runCommand executes an external command with stdout/stderr forwarded to the console.
func runCommand(name string, args ...string) (err error) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// dirExists checks whether the given directory exists.
func dirExists(path string) (exists bool) {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
