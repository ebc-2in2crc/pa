package e2e

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ebc-2in2crc/pa/cmd"
	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/pkg/errors"
)

func testE2EWebhookCreate(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "webhook create --graph-id=graph-id --type=increment"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Webhook create got: %+v\nwant: nil", err)
	}
}

func testE2EWebhookGetAll(t *testing.T) {
	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "webhook get"
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err := cmd.Execute()

	if err != nil {
		t.Errorf("Webhook get got: %+v\nwant: nil", err)
	}
}

func testE2EWebhookInvoke(t *testing.T) {
	hash, err := getWebhookHash()
	if err != nil {
		t.Error("testE2EWebhookInvoke failed", err)
	}

	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "webhook invoke --hash=" + hash
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err = cmd.Execute()

	if err != nil {
		t.Errorf("Webhook invoke got: %+v\nwant: nil", err)
	}
}

func getWebhookHash() (string, error) {
	user := os.Getenv("PA_USERNAME")
	token := os.Getenv("PA_SECOND_TOKEN")
	client := pixela.New(user, token)
	input := &pixela.WebhookCreateInput{
		GraphID: pixela.String("graph-id"),
		Type:    pixela.String(pixela.WebhookTypeIncrement),
	}
	result, err := client.Webhook().Create(input)
	if err != nil {
		e := errors.Wrapf(err, "Webhook get all failed")
		return "", e
	}
	return result.WebhookHash, nil
}

func testE2EWebhookDelete(t *testing.T) {
	hash, err := getWebhookHash()
	if err != nil {
		t.Error("testE2EWebhookDelete failed", err)
	}

	cmd := cmd.NewCmdRoot()
	cmd.SetOut(ioutil.Discard)
	commandline := "webhook delete --hash=" + hash
	args := strings.Split(commandline, " ")
	cmd.SetArgs(args)

	err = cmd.Execute()

	if err != nil {
		t.Errorf("Webhook delete got: %+v\nwant: nil", err)
	}
}
