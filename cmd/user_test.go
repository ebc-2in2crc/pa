package cmd

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	pixela "github.com/ebc-2in2crc/pixela4go"

	"github.com/stretchr/testify/assert"
)

type pixelaUserMock struct {
	result pixela.Result
	err    error
}

func (p *pixelaUserMock) Create(input *pixela.UserCreateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaUserMock) Update(input *pixela.UserUpdateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaUserMock) Delete() (*pixela.Result, error) {
	return &p.result, p.err
}

func TestUserCreateInput(t *testing.T) {
	params := []struct {
		envs        map[string]string
		commandline string
		expected    pixela.UserCreateInput
	}{
		{
			envs: map[string]string{
				"PA_AGREE_TERMS_OF_SERVICE": "true",
				"PA_NOT_MINOR":              "true",
				"PA_THANKS_CODE":            "xxx"},
			commandline: "user create",
			expected: pixela.UserCreateInput{
				AgreeTermsOfService: pixela.Bool(true),
				NotMinor:            pixela.Bool(true),
				ThanksCode:          pixela.String("xxx"),
			},
		},
		{
			envs: map[string]string{
				"PA_AGREE_TERMS_OF_SERVICE": "false",
				"PA_NOT_MINOR":              "false",
				"PA_THANKS_CODE":            "xxx"},
			commandline: "user create --agree-terms-of-service --not-minor --thanks-code=xxx2",
			expected: pixela.UserCreateInput{
				AgreeTermsOfService: pixela.Bool(true),
				NotMinor:            pixela.Bool(true),
				ThanksCode:          pixela.String("xxx2"),
			},
		},
		{
			envs: map[string]string{
				"PA_AGREE_TERMS_OF_SERVICE": "false",
				"PA_NOT_MINOR":              "false",
				"PA_THANKS_CODE":            ""},
			commandline: "user create",
			expected:    pixela.UserCreateInput{},
		},
	}

	for _, p := range params {
		setOSEnv(p.envs)
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createUserCreateInput()

		assert.EqualValues(t, p.expected.AgreeTermsOfService, input.AgreeTermsOfService, "AgreeTermsOfService")
		assert.EqualValues(t, p.expected.NotMinor, input.NotMinor, "NotMinor")
		assert.EqualValues(t, pixela.StringValue(p.expected.ThanksCode), pixela.StringValue(input.ThanksCode), "ThanksCode")
	}
}

func TestUserCreate(t *testing.T) {
	defer func() { pixelaClient.user = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		expected string
	}{
		{
			Result: pixela.Result{
				Message:    "Success.",
				IsSuccess:  true,
				StatusCode: http.StatusOK,
			},
			occur:    nil,
			expected: `{"message":"Success.","isSuccess":true,"isRejected":false,"statusCode":200}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:    "This user already exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"This user already exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `user create failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.user = &pixelaUserMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdUserCreate()
		buffer := bytes.NewBuffer([]byte{})
		c.SetOut(buffer)

		err := c.RunE(c, []string{})

		if v.occur == nil {
			assert.Equal(t, v.expected, buffer.String())
		} else {
			assert.Contains(t, err.Error(), v.expected)
		}
	}
}

func TestUserUpdateInput(t *testing.T) {
	params := []struct {
		envs        map[string]string
		commandline string
		expected    pixela.UserUpdateInput
	}{
		{
			envs: map[string]string{
				"PA_NEW_TOKEN":   "newToken",
				"PA_THANKS_CODE": "xxx"},
			commandline: "user update",
			expected: pixela.UserUpdateInput{
				NewToken:   pixela.String("newToken"),
				ThanksCode: pixela.String("xxx"),
			},
		},
		{
			envs: map[string]string{
				"PA_NEW_TOKEN":   "newToken",
				"PA_THANKS_CODE": "xxx"},
			commandline: "user update --new-token=newToken2 --thanks-code=xxx2",
			expected: pixela.UserUpdateInput{
				NewToken:   pixela.String("newToken2"),
				ThanksCode: pixela.String("xxx2"),
			},
		},
		{
			envs: map[string]string{
				"PA_NEW_TOKEN":   "",
				"PA_THANKS_CODE": ""},
			commandline: "user update",
			expected:    pixela.UserUpdateInput{},
		},
	}

	for _, p := range params {
		setOSEnv(p.envs)
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createUserUpdateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.NewToken), pixela.StringValue(input.NewToken), "NewToken")
		assert.EqualValues(t, pixela.StringValue(p.expected.ThanksCode), pixela.StringValue(input.ThanksCode), "ThanksCode")
	}
}

func TestUserUpdate(t *testing.T) {
	defer func() { pixelaClient.user = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		expected string
	}{
		{
			Result: pixela.Result{
				Message:    "Success.",
				IsSuccess:  true,
				StatusCode: http.StatusOK,
			},
			occur:    nil,
			expected: `{"message":"Success.","isSuccess":true,"isRejected":false,"statusCode":200}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:    "User foo does not exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"User foo does not exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `user update failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.user = &pixelaUserMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdUserUpdate()
		buffer := bytes.NewBuffer([]byte{})
		c.SetOut(buffer)

		err := c.RunE(c, []string{})

		if v.occur == nil {
			assert.Equal(t, v.expected, buffer.String())
		} else {
			assert.Contains(t, err.Error(), v.expected)
		}
	}
}

func TestUserDelete(t *testing.T) {
	defer func() { pixelaClient.user = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		expected string
	}{
		{
			Result: pixela.Result{
				Message:    "Success.",
				IsSuccess:  true,
				StatusCode: http.StatusOK,
			},
			occur:    nil,
			expected: `{"message":"Success.","isSuccess":true,"isRejected":false,"statusCode":200}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:    "User foo does not exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"User foo does not exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `user delete failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.user = &pixelaUserMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdUserDelete()
		buffer := bytes.NewBuffer([]byte{})
		c.SetOut(buffer)
		assert.NoError(t, c.Flags().Set("delete-me", "true"))

		err := c.RunE(c, []string{})

		if v.occur == nil {
			assert.Equal(t, v.expected, buffer.String())
		} else {
			assert.Contains(t, err.Error(), v.expected)
		}
	}
}
