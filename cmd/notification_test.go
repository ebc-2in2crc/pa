package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/stretchr/testify/assert"
)

type pixelaNotificationMock struct {
	result        pixela.Result
	err           error
	notifications []pixela.NotificationDefinition
}

func (p *pixelaNotificationMock) Create(input *pixela.NotificationCreateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaNotificationMock) GetAll(input *pixela.NotificationGetAllInput) (*pixela.NotificationDefinitions, error) {
	result := &pixela.NotificationDefinitions{
		Notifications: p.notifications,
		Result:        p.result,
	}
	return result, p.err
}

func (p *pixelaNotificationMock) Update(input *pixela.NotificationUpdateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaNotificationMock) Delete(input *pixela.NotificationDeleteInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func TestNotificationCreateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.NotificationCreateInput
	}{
		{
			commandline: "notification create --id=notification-id --name=notification-name --target=quantity " +
				"--condition=> --threshold=10 --remind-by=23 --channel-id=channel-id --graph-id=graph-id",
			expected: pixela.NotificationCreateInput{
				ID:        pixela.String("notification-id"),
				Name:      pixela.String("notification-name"),
				Target:    pixela.String("quantity"),
				Condition: pixela.String(">"),
				Threshold: pixela.String("10"),
				RemindBy:  pixela.String("23"),
				ChannelID: pixela.String("channel-id"),
				GraphID:   pixela.String("graph-id"),
			},
		},
		{
			commandline: "notification create",
			expected:    pixela.NotificationCreateInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createNotificationCreateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "ID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Name), pixela.StringValue(input.Name), "Name")
		assert.EqualValues(t, pixela.StringValue(p.expected.Target), pixela.StringValue(input.Target), "Target")
		assert.EqualValues(t, pixela.StringValue(p.expected.Condition), pixela.StringValue(input.Condition), "Condition")
		assert.EqualValues(t, pixela.StringValue(p.expected.Threshold), pixela.StringValue(input.Threshold), "Threshold")
		assert.EqualValues(t, pixela.StringValue(p.expected.RemindBy), pixela.StringValue(input.RemindBy), "RemindBy")
		assert.EqualValues(t, pixela.StringValue(p.expected.ChannelID), pixela.StringValue(input.ChannelID), "ChannelID")
		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
	}
}

func TestNotificationCreate(t *testing.T) {
	defer func() { pixelaClient.notification = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		expected string
	}{
		{
			Result: pixela.Result{
				Message:   "Success.",
				IsSuccess: true,
			},
			occur:    nil,
			expected: `{"message":"Success.","isSuccess":true}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:   "Same id notification is exist.",
				IsSuccess: false,
			},
			occur:    nil,
			expected: `{"message":"Same id notification is exist.","isSuccess":false}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `notification create failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.notification = &pixelaNotificationMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdNotificationCreate()
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

func TestNotificationGetAllInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.NotificationGetAllInput
	}{
		{
			commandline: "notification get --graph-id=graph-id",
			expected: pixela.NotificationGetAllInput{
				GraphID: pixela.String("graph-id"),
			},
		},
		{
			commandline: "notification get",
			expected:    pixela.NotificationGetAllInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createNotificationGetAllInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
	}
}

func TestNotificationGetAll(t *testing.T) {
	defer func() { pixelaClient.notification = nil }()
	params := []struct {
		Result        pixela.Result
		occur         error
		notifications []pixela.NotificationDefinition
		expected      string
	}{
		{
			Result: pixela.Result{
				Message:   "Success.",
				IsSuccess: true,
			},
			occur: nil,
			notifications: []pixela.NotificationDefinition{
				{
					ID:        "notification-id",
					Name:      "notification-name",
					Target:    "quantity",
					Condition: ">",
					Threshold: "1",
					RemindBy:  "23",
					ChannelID: "channel-id",
				},
			},
			expected: `{"notifications":[{"id":"notification-id","name":"notification-name","target":"quantity","condition":">","threshold":"1","remindBy":"23","channelID":"channel-id"}]}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:   "Specified graphID not exist.",
				IsSuccess: false,
			},
			occur:    nil,
			expected: `{"message":"Specified graphID not exist.","isSuccess":false}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `notification get all failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.notification = &pixelaNotificationMock{
			result:        v.Result,
			err:           v.occur,
			notifications: v.notifications,
		}
		c := NewCmdNotificationGetAll()
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

func TestNotificationUpdateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.NotificationUpdateInput
	}{
		{
			commandline: "notification update --id=notification-id --name=notification-name --target=quantity " +
				"--condition=> --threshold=10 --remind-by=23 --channel-id=channel-id --graph-id=graph-id",
			expected: pixela.NotificationUpdateInput{
				ID:        pixela.String("notification-id"),
				Name:      pixela.String("notification-name"),
				Target:    pixela.String("quantity"),
				Condition: pixela.String(">"),
				Threshold: pixela.String("10"),
				RemindBy:  pixela.String("23"),
				ChannelID: pixela.String("channel-id"),
				GraphID:   pixela.String("graph-id"),
			},
		},
		{
			commandline: "notification update",
			expected:    pixela.NotificationUpdateInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createNotificationUpdateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "ID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Name), pixela.StringValue(input.Name), "Name")
		assert.EqualValues(t, pixela.StringValue(p.expected.Target), pixela.StringValue(input.Target), "Target")
		assert.EqualValues(t, pixela.StringValue(p.expected.Condition), pixela.StringValue(input.Condition), "Condition")
		assert.EqualValues(t, pixela.StringValue(p.expected.Threshold), pixela.StringValue(input.Threshold), "Threshold")
		assert.EqualValues(t, pixela.StringValue(p.expected.RemindBy), pixela.StringValue(input.RemindBy), "RemindBy")
		assert.EqualValues(t, pixela.StringValue(p.expected.ChannelID), pixela.StringValue(input.ChannelID), "ChannelID")
		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
	}
}

func TestNotificationUpdate(t *testing.T) {
	defer func() { pixelaClient.notification = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		expected string
	}{
		{
			Result: pixela.Result{
				Message:   "Success.",
				IsSuccess: true,
			},
			occur:    nil,
			expected: `{"message":"Success.","isSuccess":true}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:   "Specified graphID not exist.",
				IsSuccess: false,
			},
			occur:    nil,
			expected: `{"message":"Specified graphID not exist.","isSuccess":false}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `notification update failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.notification = &pixelaNotificationMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdNotificationUpdate()
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

func TestNotificationDeleteInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.NotificationDeleteInput
	}{
		{
			commandline: "notification delete --id=notification-id --graph-id=graph-id",
			expected: pixela.NotificationDeleteInput{
				ID:      pixela.String("notification-id"),
				GraphID: pixela.String("graph-id"),
			},
		},
		{
			commandline: "notification delete",
			expected:    pixela.NotificationDeleteInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createNotificationDeleteInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "ID")
		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
	}
}

func TestNotificationDelete(t *testing.T) {
	defer func() { pixelaClient.notification = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		expected string
	}{
		{
			Result: pixela.Result{
				Message:   "Success.",
				IsSuccess: true,
			},
			occur:    nil,
			expected: `{"message":"Success.","isSuccess":true}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:   "Specified id notification is not exist.",
				IsSuccess: false,
			},
			occur:    nil,
			expected: `{"message":"Specified id notification is not exist.","isSuccess":false}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `notification delete failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.notification = &pixelaNotificationMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdNotificationDelete()
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
