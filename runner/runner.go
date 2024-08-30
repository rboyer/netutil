package runner

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

var (
	debugExec bool
)

func SetDebug(v bool) {
	debugExec = v
}

func ExecSimpleOutput(
	name string,
	args []string,
) (string, error) {
	return ExecOutput(name, args, nil, nil, "")
}

func ExecSimple(name string, args []string) error {
	_, err := Exec(name, args, nil, nil, nil, "")
	return err
}

func ExecSimpleExitCode(name string, args []string) (int, error) {
	code, err := Exec(name, args, nil, nil, nil, "")
	if err != nil && code != 0 {
		return code, nil
	}
	return code, err
}

func ExecOutput(
	name string,
	args []string,
	stderr io.Writer,
	stdin io.Reader,
	dir string,
) (string, error) {
	buf := &bytes.Buffer{}

	_, err := Exec(name, args, buf, stderr, stdin, dir)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func Exec(
	name string,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
	stdin io.Reader,
	dir string,
) (int, error) {
	var errWriter bytes.Buffer

	if stdout == nil {
		stdout = io.Discard
	}

	if debugExec {
		if len(args) > 0 {
			log.Printf("EXEC: %s %s", name, strings.Join(args, " "))
		} else {
			log.Printf("EXEC: %s", name)
		}
	}

	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = stdout
	cmd.Stderr = &errWriter
	if stderr != nil {
		cmd.Stderr = io.MultiWriter(stderr, cmd.Stderr)
	}
	cmd.Stdin = stdin
	err := cmd.Run()
	if err != nil {
		var code int
		if ee := (&exec.ExitError{}); errors.As(err, &ee) {
			code = ee.ExitCode()
		}

		return code, fmt.Errorf("run[%s]: <<%s>> %w", name, errWriter.String(), err)
	}

	return 0, nil
}
