package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/stretchr/testify/assert"
)

type pixelaWebhookMock struct {
	result      pixela.Result
	err         error
	WebhookHash string
	Webhooks    []pixela.WebhookDefinition
}

func (p *pixelaWebhookMock) Create(input *pixela.WebhookCreateInput) (*pixela.WebhookCreateResult, error) {
	result := &pixela.WebhookCreateResult{
		WebhookHash: p.WebhookHash,
		Result:      p.result,
	}
	return result, p.err
}

func (p *pixelaWebhookMock) GetAll() (*pixela.WebhookDefinitions, error) {
	result := &pixela.WebhookDefinitions{
		Webhooks: p.Webhooks,
		Result:   p.result,
	}
	return result, p.err
}

func (p *pixelaWebhookMock) Invoke(input *pixela.WebhookInvokeInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaWebhookMock) Delete(input *pixela.WebhookDeleteInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func TestWebhookCreateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.WebhookCreateInput
	}{
		{
			commandline: "webhook create --graph-id=graph-id --type=increment",
			expected: pixela.WebhookCreateInput{
				GraphID: pixela.String("graph-id"),
				Type:    pixela.String("increment"),
			},
		},
		{
			commandline: "webhook create",
			expected:    pixela.WebhookCreateInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createWebhookCreateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Type), pixela.StringValue(input.Type), "Type")
	}
}

func TestWebhookCreate(t *testing.T) {
	defer func() { pixelaClient.webhook = nil }()
	params := []struct {
		Result      pixela.Result
		occur       error
		webhookHash string
		expected    string
	}{
		{
			Result: pixela.Result{
				Message:    "Success.",
				IsSuccess:  true,
				StatusCode: http.StatusOK,
			},
			occur:       nil,
			webhookHash: "webhook hash",
			expected:    `{"webhookHash":"webhook hash","message":"Success.","isSuccess":true,"isRejected":false,"statusCode":200}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:    "It is necessary to specify graphID, and type.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"webhookHash":"","message":"It is necessary to specify graphID, and type.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `webhook create failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.webhook = &pixelaWebhookMock{
			result:      v.Result,
			err:         v.occur,
			WebhookHash: v.webhookHash,
		}
		c := NewCmdWebhookCreate()
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

func _TestWebhookGetInput(t *testing.T) {
}

func TestWebhookGetAll(t *testing.T) {
	defer func() { pixelaClient.webhook = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		webhooks []pixela.WebhookDefinition
		expected string
	}{
		{
			Result: pixela.Result{
				Message:    "Success.",
				IsSuccess:  true,
				StatusCode: http.StatusOK,
			},
			occur: nil,
			webhooks: []pixela.WebhookDefinition{
				{
					WebhookHash: "webhook hash",
					GraphID:     "graph id",
					Type:        "increment",
				},
			},
			expected: `{"webhooks":[{"webhookHash":"webhook hash","graphId":"graph id","type":"increment"}]}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:    "User does not exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"User does not exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `webhook get all failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.webhook = &pixelaWebhookMock{
			result:   v.Result,
			err:      v.occur,
			Webhooks: v.webhooks,
		}
		c := NewCmdWebhookGetAll()
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

func TestWebhookInvokeInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.WebhookInvokeInput
	}{
		{
			commandline: "webhook invoke --hash=webhookhash",
			expected: pixela.WebhookInvokeInput{
				WebhookHash: pixela.String("webhookhash"),
			},
		},
		{
			commandline: "webhook invoke",
			expected:    pixela.WebhookInvokeInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createWebhookInvokeInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.WebhookHash), pixela.StringValue(input.WebhookHash), "WebhookHash")
	}
}

func TestWebhookInvoke(t *testing.T) {
	defer func() { pixelaClient.webhook = nil }()
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
				Message:    "Specified user or webhook not exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"Specified user or webhook not exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `webhook invoke failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.webhook = &pixelaWebhookMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdWebhookInvoke()
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

func TestWebhookDeleteInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.WebhookDeleteInput
	}{
		{
			commandline: "webhook delete --hash=webhookhash",
			expected: pixela.WebhookDeleteInput{
				WebhookHash: pixela.String("webhookhash"),
			},
		},
		{
			commandline: "webhook delete",
			expected:    pixela.WebhookDeleteInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createWebhookDeleteInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.WebhookHash), pixela.StringValue(input.WebhookHash), "WebhookHash")
	}
}

func TestWebhookDelete(t *testing.T) {
	defer func() { pixelaClient.webhook = nil }()
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
				Message:    "Specified user or webhook not exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"Specified user or webhook not exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `webhook delete failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.webhook = &pixelaWebhookMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdWebhookDelete()
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
