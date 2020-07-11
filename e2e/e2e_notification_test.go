package e2e

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ebc-2in2crc/pa/cmd"
)

func testE2ENotificationCreate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "notification create --id=notification-id --name=notification-name --target=quantity " +
		"--condition=> --threshold=1 --remind-by=23 --channel-id=channel-id --graph-id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Notification create got: %+v\nwant: nil", err)
	}
}

func testE2ENotificationGetAll(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "notification get --graph-id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Notification get all got: %+v\nwant: nil", err)
	}
}

func testE2ENotificationUpdate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "notification update --id=notification-id --name=notification-name --target=quantity " +
		"--condition=> --threshold=1 --remind-by=23 --channel-id=channel-id --graph-id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Notification update got: %+v\nwant: nil", err)
	}
}

func testE2ENotificationDelete(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "notification delete --id=notification-id --graph-id=graph-id"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Notification delete got: %+v\nwant: nil", err)
	}
}
