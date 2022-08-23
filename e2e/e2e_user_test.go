package e2e

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ebc-2in2crc/pa/cmd"
)

func testE2EUserCreate(t *testing.T) {
	assert.NoError(t, os.Setenv("PA_TOKEN", os.Getenv("PA_FIRST_TOKEN")))

	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "user create --agree-terms-of-service --not-minor"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("User create got: %+v\nwant: nil", err)
	}
}

func testE2EUserUpdate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	newToken := os.Getenv("PA_SECOND_TOKEN")
	commandline := "user update --new-token=" + newToken
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("User update got: %+v\nwant: nil", err)
	}

	assert.NoError(t, os.Setenv("PA_TOKEN", newToken))
}

func testE2EUserDelete(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(io.Discard)
	commandline := "user delete --delete-me"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("User delete got: %+v\nwant: nil", err)
	}
}
