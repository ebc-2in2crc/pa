package e2e

import (
	"io"
	"strings"
	"testing"

	"github.com/ebc-2in2crc/pa/cmd"
)

func testE2EPixelCreate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "pixel create --graph-id=graph-id --date=20200101 --quantity=5"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("pixel create got: %+v\nwant: nil", err)
	}
}

func testE2EPixelIncrement(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "pixel increment --graph-id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("pixel increment got: %+v\nwant: nil", err)
	}
}

func testE2EPixelDecrement(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "pixel decrement --graph-id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("pixel decrement got: %+v\nwant: nil", err)
	}
}

func testE2EPixelGet(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "pixel get --graph-id=graph-id --date=20200101"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("pixel get got: %+v\nwant: nil", err)
	}
}

func testE2EPixelUpdate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "pixel update --graph-id=graph-id --date=20200101 --quantity=5"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("pixel update got: %+v\nwant: nil", err)
	}
}

func testE2EPixelDelete(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "pixel delete --graph-id=graph-id --date=20200101"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("pixel delete got: %+v\nwant: nil", err)
	}
}
