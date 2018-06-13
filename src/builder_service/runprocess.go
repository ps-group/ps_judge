package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// ProcessLimits - process limits
type processLimits struct {
	NumberOfFiles int // max number of open files
	NumberOfProc  int // max number of processors
	NumberOfLocks int // max number of file locks
	TimeInSeconds int // max CPU time, in seconds
	AddessSpaceMB int // max address space size, in mebabytes
}

type solutionError struct {
	reason string
}

func (e *solutionError) Error() string {
	return "solution failed: " + e.reason
}

type testCase struct {
	inputPath    string
	outputPath   string
	expectedPath string
}

type processRunOptions struct {
	workdir    string
	stdinPath  string
	stdoutPath string
	limits     *processLimits
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

	process := exec.Command("prlimit", prlimitArgs...)
	process.Dir = options.workdir
	stdin, err := process.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := process.StdoutPipe()
	if err != nil {
		return err
	}
	err = process.Start()
	if err != nil {
		return err
	}
	stdinBytes, err := ioutil.ReadFile(options.stdinPath)
	if err != nil {
		return err
	}
	_, err = stdin.Write(stdinBytes)
	if err != nil {
		return err
	}
	err = stdin.Close()
	if err != nil {
		return err
	}
	err = process.Wait()
	if err != nil {
		return &solutionError{
			reason: err.Error(),
		}
	}
	outputBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	expectedBytes, err := ioutil.ReadFile(options.stdoutPath)
	if err != nil {
		return err
	}
	if string(outputBytes) != string(expectedBytes) {
		return &solutionError{
			reason: "output doesn't match expected",
		}
	}
	return nil
}

// compileSolution - compiles given source code file
func compileSolution(filepath string, language language, outputPath string) error {
	var cmd *exec.Cmd
	if language == languagePascal {
		cmd = exec.Command("gpc", filepath, outputPath, "-w", "--extended-syntax", "--implicit-result")
	} else if language == languageCpp {
		cmd = exec.Command("gcc", filepath, outputPath, "--std=c++17")
	} else {
		return errors.New("unknown language value passed")
	}
	return cmd.Run()
}

func checkSolution(executablePath string, cases []testCase, workdir string) []error {
	var errors []error
	for _, c := range cases {
		limits := newProcessLimits()
		options := processRunOptions{
			limits:     limits,
			workdir:    workdir,
			stdinPath:  c.inputPath,
			stdoutPath: c.outputPath,
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

func buildSolution(sourceCode string, language language, cases []testCase, workdir string) BuildResult {
	srcPath := filepath.Join(workdir, "solution"+getLanguageExt(language))
	exePath := filepath.Join(workdir, "solution")
	runWorkdir := filepath.Join(workdir, "run")
	err := os.MkdirAll(runWorkdir, 0)
	if err != nil {
		return BuildResult{
			internalError: err,
		}
	}

	err = ioutil.WriteFile(srcPath, []byte(sourceCode), 0)
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
