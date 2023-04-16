package e2e

import (
	"io"
	"strings"
	"testing"

	"github.com/ebc-2in2crc/pa/cmd"
)

func testE2EGraphCreate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph create --id=graph-id --name=graph-name --unit=times --type=int --color=sora" +
		" --timezone=Asia/Tokyo --self-sufficient=none"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph create got: %+v\nwant: nil", err)
	}
}

func testE2EGraphGetAll(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph get-all"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph get all got: %+v\nwant: nil", err)
	}
}

func testE2EGraphGetSVG(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph svg --id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph get svg got: %+v\nwant: nil", err)
	}
}

func testE2EGraphStats(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph stats --id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph stats got: %+v\nwant: nil", err)
	}
}

func testE2EGraphUpdate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph update --id=graph-id --name=graph-name --unit=times --color=sora" +
		" --timezone=Asia/Tokyo --self-sufficient=none"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph update got: %+v\nwant: nil", err)
	}
}

func testE2EGraphAdd(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph add --id=graph-id --quantity=1"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph update got: %+v\nwant: nil", err)
	}
}

func testE2EGraphSubtract(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph subtract --id=graph-id --quantity=1"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph subtract got: %+v\nwant: nil", err)
	}
}

func testE2EGraphGetPixelDates(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph pixels --id=graph-id --from=20200101 --to=20200130"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph pixels got: %+v\nwant: nil", err)
	}
}

func testE2EGraphStopwatch(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph stopwatch --id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph stopwatch got: %+v\nwant: nil", err)
	}
}

func testE2EGraphDelete(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "graph delete --id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("graph delete got: %+v\nwant: nil", err)
	}
}
