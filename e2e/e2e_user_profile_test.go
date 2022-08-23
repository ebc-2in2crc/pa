package e2e

import (
	"io"
	"strings"
	"testing"

	"github.com/ebc-2in2crc/pa/cmd"
)

func testE2EUserProfileUpdate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "profile update --display-name=display-name"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("User update got: %+v\nwant: nil", err)
	}
}
