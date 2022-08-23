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
		expectedRetry    int
	}{
		{
			envs:             map[string]string{"PA_USERNAME": "pa-user", "PA_TOKEN": "pa-token", "PA_RETRY": "5"},
			commandline:      "",
			expectedUserName: "pa-user",
			expectedToken:    "pa-token",
			expectedRetry:    5,
		},
		{
			envs:             map[string]string{"PA_USERNAME": "pa-user", "PA_TOKEN": "pa-token", "PA_RETRY": "5"},
			commandline:      "--username=papa-user --token=papa-token --retry=10",
			expectedUserName: "papa-user",
			expectedToken:    "papa-token",
			expectedRetry:    10,
		},
	}

	for _, p := range params {
		setOSEnv(p.envs)
		cmd := NewCmdRoot()
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.SetOut(ioutil.Discard)
		_ = cmd.Execute()

		if getUsername() != p.expectedUserName {
			t.Errorf("expected username: %s, but got %s", p.expectedUserName, getUsername())
		}
		if getToken() != p.expectedToken {
			t.Errorf("expected token: %s, but got %s", p.expectedToken, getToken())
		}
		if getRetry() != p.expectedRetry {
			t.Errorf("expected retry: %d, but got %d", p.expectedRetry, getRetry())
		}
	}
}

func setOSEnv(m map[string]string) {
	for k, v := range m {
		_ = os.Setenv(k, v)
	}
}
