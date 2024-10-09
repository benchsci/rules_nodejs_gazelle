package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-gazelle/testtools"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

const (
	testDataPath      = "tests/"
	gazelleBinaryName = "test_gazelle_bin"
)

func main() {
	// Check that we have exactly one argument
	if len(os.Args) != 2 {
		panic("expected exactly one argument")
	}

	// Get the folder name from the command line
	name := os.Args[1]

	// Wrap test with InternalTest so that it can be run by testing.Main
	theTest := testing.InternalTest{
		Name: fmt.Sprintf("test_%s", name),
		F:    func(t *testing.T) { RunTest(t, name) },
	}

	// Create a matchAll function that will match all tests
	matchAll := func(pat, str string) (bool, error) {
		return true, nil // Run all tests in list
	}

	// Run the test
	testing.Main(matchAll, []testing.InternalTest{theTest}, nil, nil)
}

func RunTest(t *testing.T, name string) {

	// Get path to gazelle binary
	gazellePath, ok := bazel.FindBinary("", gazelleBinaryName)
	if !ok {
		t.Errorf("could not find gazelle binary")
		t.FailNow()
		return
	}

	// Get all Bazel runfiles
	runfiles, err := bazel.ListRunfiles()
	if err != nil {
		t.Fatalf("bazel.ListRunfiles() error: %v", err)
	}

	var inputs []testtools.FileSpec
	var goldens []testtools.FileSpec

	// Get the path to the test directory
	testDirShortPath := testDataPath + name

	for _, f := range runfiles {

		// Exclude runfiles outside of test directory
		if !strings.HasPrefix(f.ShortPath, testDirShortPath) {
			continue
		}

		// Get relative file path from inside test directory
		testFilePath := strings.TrimPrefix(f.ShortPath, testDirShortPath)
		fmt.Println(testFilePath)
		info, err := os.Stat(f.Path)
		if err != nil {
			t.Fatalf("os.Stat(%q) error: %v", f.Path, err)
		}

		// Skip directories
		if info.IsDir() {
			continue
		}

		// Read file contents
		content, err := os.ReadFile(f.Path)
		if err != nil {
			t.Errorf("os.ReadFile(%q) error: %v", f.Path, err)
		}

		// Add file to inputs or goldens
		if strings.HasSuffix(testFilePath, ".in") {
			inputs = append(inputs, testtools.FileSpec{
				Path:    filepath.Join(name, strings.TrimSuffix(testFilePath, ".in")),
				Content: string(content),
			})
		} else if strings.HasSuffix(testFilePath, ".out") {
			goldens = append(goldens, testtools.FileSpec{
				Path:    filepath.Join(name, strings.TrimSuffix(testFilePath, ".out")),
				Content: string(content),
			})
		} else {
			inputs = append(inputs, testtools.FileSpec{
				Path:    filepath.Join(name, testFilePath),
				Content: string(content),
			})
			goldens = append(goldens, testtools.FileSpec{
				Path:    filepath.Join(name, testFilePath),
				Content: string(content),
			})
		}
	}

	// Create temporary directory with input files
	testdataDir, cleanup := testtools.CreateFiles(t, inputs)
	defer cleanup()
	defer func() {
		if t.Failed() {
			filepath.Walk(testdataDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				t.Logf("%q exists", strings.TrimPrefix(path, testdataDir))
				return nil
			})
		}
	}()

	// Run gazelle
	workspaceRoot := filepath.Join(testdataDir, name)
	args := []string{"-build_file_name=BUILD,BUILD.bazel"}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, gazellePath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = workspaceRoot
	if err := cmd.Run(); err != nil {
		var e *exec.ExitError
		if !errors.As(err, &e) {
			t.Fatal(err)
		}
	}

	// Check exit code and stdout/stderr
	procExitCode := cmd.ProcessState.ExitCode()
	if procExitCode != 0 {
		t.Errorf("expected gazelle exit code: 0\ngot: %d", procExitCode)
	}
	procStdout := stdout.String()
	if strings.TrimSpace(procStdout) != "" {
		t.Errorf("expected no gazelle stdout\ngot: %s", procStdout)
	}
	procStderr := stderr.String()
	if strings.TrimSpace(procStderr) != "" {
		t.Errorf("expected no gazelle stderr\ngot: %s", procStderr)
	}

	// Check output files
	testtools.CheckFiles(t, testdataDir, goldens)
}
