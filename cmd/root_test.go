package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCmdRootFlags(t *testing.T) {
	params := []struct {
		envs             map[string]string
		commandline      string
		expectedUserName string
		expectedToken    string
	}{
		{
			envs:             map[string]string{"PA_USERNAME": "pa-user", "PA_TOKEN": "pa-token"},
			commandline:      "",
			expectedUserName: "pa-user",
			expectedToken:    "pa-token",
		},
		{
			envs:             map[string]string{"PA_USERNAME": "pa-user", "PA_TOKEN": "pa-token"},
			commandline:      "--username=papa-user --token=papa-token",
			expectedUserName: "papa-user",
			expectedToken:    "papa-token",
		},
	}

	for _, p := range params {
		setOSEnv(p.envs)
		cmd := NewCmdRoot()
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.SetOut(ioutil.Discard)
		cmd.Execute()

		if getUsername() != p.expectedUserName {
			t.Errorf("expected username: %s, but got %s", p.expectedUserName, getUsername())
		}
		if getToken() != p.expectedToken {
			t.Errorf("expected token: %s, but got %s", p.expectedToken, getToken())
		}
	}
}

func setOSEnv(m map[string]string) {
	for k, v := range m {
		os.Setenv(k, v)
	}
}
