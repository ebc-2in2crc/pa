package e2e

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ebc-2in2crc/pa/cmd"
)

func testE2EChannelCreate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "channel create --id=channel-id --name=channel-name --type=slack " +
		"--slack-username=user --slack-channel-name=channel " +
		"--slack-url=https://hooks.slack.com/services/foo/bar/baz/xxx"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Channel create got: %+v\nwant: nil", err)
	}
}

func testE2EChannelGetAll(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "channel get"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Channel get all got: %+v\nwant: nil", err)
	}
}

func testE2EChannelUpdate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "channel update --id=channel-id --name=channel-name --type=slack " +
		"--slack-username=user --slack-channel-name=channel " +
		"--slack-url=https://hooks.slack.com/services/foo/bar/baz/xxx"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Channel update got: %+v\nwant: nil", err)
	}
}

func testE2EChannelDelete(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "channel delete --id=channel-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Channel delete got: %+v\nwant: nil", err)
	}
}
