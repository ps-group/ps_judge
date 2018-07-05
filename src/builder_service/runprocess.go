package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

// ProcessLimits - process limits
type processLimits struct {
	NumberOfFiles int // max number of open files
	NumberOfProc  int // max number of processors
	NumberOfLocks int // max number of file locks
	TimeInSeconds int // max CPU time, in seconds
	AddessSpaceMB int // max address space size, in mebabytes
}

type TestCase struct {
	Input    string
	Expected string
}

type processRunOptions struct {
	workdir  string
	input    string
	expected string
	limits   *processLimits
}

// newProcessLimits - creates new ProcessLimits with default values
func newProcessLimits() *processLimits {
	limits := new(processLimits)
	limits.NumberOfFiles = 8
	limits.NumberOfProc = 1
	limits.NumberOfLocks = 8
	limits.TimeInSeconds = 2
	limits.AddessSpaceMB = 256
	return limits
}

// runLimitedProcess - calls command wrapped with prlimit to limit resources
// TODO: make type runResult which holds two errors: internal and runtime
func runLimitedProcess(options processRunOptions, cmd string, arg ...string) error {
	argNofile := fmt.Sprintf("--nofile=%d", options.limits.NumberOfFiles)
	argCPUTime := fmt.Sprintf("--cpu=%d", options.limits.TimeInSeconds)
	argCPUCount := fmt.Sprintf("--nproc=%d", options.limits.NumberOfProc)
	argFileLocks := fmt.Sprintf("--locks=%d", options.limits.NumberOfLocks)
	argFile := fmt.Sprintf("--nofile=%d", options.limits.NumberOfFiles)
	argMemory := fmt.Sprintf("--as=%d", options.limits.AddessSpaceMB*1024*1024)

	prlimitArgs := []string{argNofile, argCPUTime, argCPUCount, argFileLocks, argFile, argMemory}
	prlimitArgs = append(prlimitArgs, cmd)
	prlimitArgs = append(prlimitArgs, arg...)

	// process := exec.Command("prlimit", prlimitArgs...)
	process := exec.Command(cmd)
	var stdin, stdout, stderr bytes.Buffer
	_, err := stdin.WriteString(options.input)
	if err != nil {
		return errors.Wrap(err, "cannot write into stdin pipe")
	}

	process.Dir = options.workdir
	process.Stdin = &stdin
	process.Stdout = &stdout
	process.Stderr = &stderr

	err = process.Run()
	if err != nil {
		reason := fmt.Sprintf("run failed: %s", err.Error())
		errText := string(stderr.Bytes())
		if len(errText) > 0 {
			reason += "\n"
			reason += errText
		}
		outText := string(stdout.Bytes())
		if len(outText) > 0 {
			reason += "\n"
			reason += outText
		}
		return errors.New(reason)
	}

	output := string(stdout.Bytes())

	// TODO: allow fuzzy comparison (ignore extra whitespace at end)

	if output != options.expected {
		return errors.New(fmt.Sprintf(
			"output does not match expected:\n--OUTPUT--\n%s\n--EXPECTED--\n%s",
			output,
			options.expected))
	}

	return nil
}

// compileSolution - compiles given source code file
func compileSolution(filepath string, language language, outputPath string) error {
	var cmd *exec.Cmd
	if language == languagePascal {
		// GNU Pascal is outdated - we don't use it anymore.
		// cmd = exec.Command("gpc", filepath, outputPath, "-w", "--extended-syntax", "--implicit-result")
		cmd = exec.Command("fpc", "-Mtp", "-So", "-o"+outputPath, filepath)
	} else if language == languageCpp {
		cmd = exec.Command("gcc", filepath, "-o", outputPath, "--std=c++17")
	} else {
		return errors.New("unknown language value passed")
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		reason := fmt.Sprintf("compilation failed: %s", err.Error())
		errText := string(stderr.Bytes())
		if len(errText) > 0 {
			reason += "\n"
			reason += errText
		}
		outText := string(stdout.Bytes())
		if len(outText) > 0 {
			reason += "\n"
			reason += outText
		}
		return errors.New(reason)
	}
	return nil
}

func checkSolution(executablePath string, cases []TestCase, workdir string) []error {
	executablePath, _ = filepath.Abs(executablePath)

	var errors []error
	for _, c := range cases {
		limits := newProcessLimits()
		options := processRunOptions{
			limits:   limits,
			workdir:  workdir,
			input:    c.Input,
			expected: c.Expected,
		}
		err := runLimitedProcess(options, executablePath)
		errors = append(errors, err)
	}
	return errors
}

type BuildResult struct {
	internalError  error
	buildError     error
	testCaseErrors []error
}

func getLanguageExt(language language) string {
	switch language {
	case languageCpp:
		return ".cpp"
	case languagePascal:
		return ".pas"
	}
	return ".unknown"
}

func buildSolution(sourceCode string, language language, cases []TestCase, workdir string) BuildResult {
	srcPath := filepath.Join(workdir, "solution"+getLanguageExt(language))
	exePath := filepath.Join(workdir, "solution")
	runWorkdir := filepath.Join(workdir, "run")
	err := os.MkdirAll(runWorkdir, os.ModePerm)
	if err != nil {
		return BuildResult{
			internalError: err,
		}
	}

	err = ioutil.WriteFile(srcPath, []byte(sourceCode), os.ModePerm)
	if err != nil {
		return BuildResult{
			internalError: err,
		}
	}
	err = compileSolution(srcPath, language, exePath)
	if err != nil {
		return BuildResult{
			buildError: err,
		}
	}
	errs := checkSolution(exePath, cases, workdir)
	return BuildResult{
		testCaseErrors: errs,
	}
}
