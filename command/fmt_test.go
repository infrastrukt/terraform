package command

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestFmt_errorReporting(t *testing.T) {
	tempDir, err := fmtFixtureWriteDir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(tempDir)

	ui := new(cli.MockUi)
	c := &FmtCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	dummy_file := filepath.Join(tempDir, "doesnotexist")
	args := []string{dummy_file}
	if code := c.Run(args); code != 2 {
		t.Fatalf("wrong exit code. errors: \n%s", ui.ErrorWriter.String())
	}

	expected := fmt.Sprintf("Error running fmt: stat %s: no such file or directory", dummy_file)
	if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
		t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
	}
}

func TestFmt_tooManyArgs(t *testing.T) {
	ui := new(cli.MockUi)
	c := &FmtCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{"bad"}
	if code := c.Run(args); code != 1 {
		t.Fatalf("wrong exit code. errors: \n%s", ui.ErrorWriter.String())
	}

	expected := "The fmt command expects no arguments."
	if actual := ui.ErrorWriter.String(); !strings.Contains(actual, expected) {
		t.Fatalf("expected:\n%s\n\nto include: %q", actual, expected)
	}
}

func TestFmt_workingDirectory(t *testing.T) {
	tempDir, err := fmtFixtureWriteDir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(tempDir)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Chdir(cwd)

	ui := new(cli.MockUi)
	c := &FmtCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", ui.ErrorWriter.String())
	}

	expected := fmtFixture.golden
	if actual := ui.OutputWriter.Bytes(); !bytes.Equal(actual, expected) {
		t.Fatalf("got: %q\nexpected: %q", actual, expected)
	}
}

var fmtFixture = struct {
	filename      string
	input, golden []byte
}{
	"main.tf",
	[]byte(`  foo  =  "bar"
`),
	[]byte(`foo = "bar"
`),
}

func fmtFixtureWriteDir() (string, error) {
	dir, err := ioutil.TempDir("", "tf")
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filepath.Join(dir, fmtFixture.filename), fmtFixture.input, 0644)
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}

	return dir, nil
}
