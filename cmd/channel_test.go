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

type pixelaChannelMock struct {
	result   pixela.Result
	err      error
	channels []pixela.ChannelDefinition
}

func (p *pixelaChannelMock) Create(input *pixela.ChannelCreateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaChannelMock) GetAll() (*pixela.ChannelDefinitions, error) {
	result := &pixela.ChannelDefinitions{
		Channels: p.channels,
		Result:   p.result,
	}
	return result, p.err
}

func (p *pixelaChannelMock) Update(input *pixela.ChannelUpdateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaChannelMock) Delete(input *pixela.ChannelDeleteInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func TestChannelCreateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.ChannelCreateInput
	}{
		{
			commandline: "channel create --id=channel-id --name=channel-name --type=slack " +
				"--slack-username=user --slack-channel-name=channel --slack-url=url",
			expected: pixela.ChannelCreateInput{
				ID:   pixela.String("channel-id"),
				Name: pixela.String("channel-name"),
				Type: pixela.String(pixela.ChannelTypeSlack),
				SlackDetail: &pixela.SlackDetail{
					UserName:    pixela.String("user"),
					ChannelName: pixela.String("channel"),
					URL:         pixela.String("url"),
				},
			},
		},
		{
			commandline: "channel create",
			expected:    pixela.ChannelCreateInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createChannelCreateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "ID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Name), pixela.StringValue(input.Name), "Name")
		assert.EqualValues(t, pixela.StringValue(p.expected.Type), pixela.StringValue(input.Type), "Type")

		if p.expected.SlackDetail == nil {
			continue
		}
		expectedDetail := p.expected.SlackDetail
		gotDetail := input.SlackDetail
		assert.EqualValues(t, pixela.StringValue(expectedDetail.UserName), pixela.StringValue(gotDetail.UserName), "UserName")
		assert.EqualValues(t, pixela.StringValue(expectedDetail.ChannelName), pixela.StringValue(gotDetail.ChannelName), "ChannelName")
		assert.EqualValues(t, pixela.StringValue(expectedDetail.URL), pixela.StringValue(gotDetail.URL), "URL")
	}
}

func TestChannelCreate(t *testing.T) {
	defer func() { pixelaClient.channel = nil }()
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
				Message:   "Same id channel is exist.",
				IsSuccess: false,
			},
			occur:    nil,
			expected: `{"message":"Same id channel is exist.","isSuccess":false}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `channel create failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.channel = &pixelaChannelMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdChannelCreate()
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

func _TestChannelGetAllInput(t *testing.T) {
}

func TestChannelGetAll(t *testing.T) {
	defer func() { pixelaClient.channel = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		channels []pixela.ChannelDefinition
		expected string
	}{
		{
			Result: pixela.Result{
				Message:   "Success.",
				IsSuccess: true,
			},
			occur: nil,
			channels: []pixela.ChannelDefinition{
				{
					ID:   "ch-id",
					Name: "ch-name",
					Type: "slack",
					Detail: pixela.SlackDetail{
						UserName:    pixela.String("slack-user"),
						ChannelName: pixela.String("slack-ch"),
						URL:         pixela.String("slack-url"),
					},
				},
			},
			expected: `{"channels":[{"id":"ch-id","name":"ch-name","type":"slack","detail":{"url":"slack-url","userName":"slack-user","channelName":"slack-ch"}}]}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:   "User does not exist.",
				IsSuccess: false,
			},
			occur:    nil,
			expected: `{"message":"User does not exist.","isSuccess":false}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `channel get all failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.channel = &pixelaChannelMock{
			result:   v.Result,
			err:      v.occur,
			channels: v.channels,
		}
		c := NewCmdChannelGetAll()
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

func TestChannelUpdateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.ChannelUpdateInput
	}{
		{
			commandline: "channel update --id=channel-id --name=channel-name --type=slack " +
				"--slack-username=user --slack-channel-name=channel --slack-url=url",
			expected: pixela.ChannelUpdateInput{
				ID:   pixela.String("channel-id"),
				Name: pixela.String("channel-name"),
				Type: pixela.String(pixela.ChannelTypeSlack),
				SlackDetail: &pixela.SlackDetail{
					UserName:    pixela.String("user"),
					ChannelName: pixela.String("channel"),
					URL:         pixela.String("url"),
				},
			},
		},
		{
			commandline: "channel update",
			expected:    pixela.ChannelUpdateInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createChannelUpdateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "ID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Name), pixela.StringValue(input.Name), "Name")
		assert.EqualValues(t, pixela.StringValue(p.expected.Type), pixela.StringValue(input.Type), "Type")

		if p.expected.SlackDetail == nil {
			continue
		}
		expectedDetail := p.expected.SlackDetail
		gotDetail := input.SlackDetail
		assert.EqualValues(t, pixela.StringValue(expectedDetail.UserName), pixela.StringValue(gotDetail.UserName), "UserName")
		assert.EqualValues(t, pixela.StringValue(expectedDetail.ChannelName), pixela.StringValue(gotDetail.ChannelName), "ChannelName")
		assert.EqualValues(t, pixela.StringValue(expectedDetail.URL), pixela.StringValue(gotDetail.URL), "URL")
	}
}

func TestChannelUpdate(t *testing.T) {
	defer func() { pixelaClient.channel = nil }()
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
				Message:   "Specified channel is not exist.",
				IsSuccess: false,
			},
			occur:    nil,
			expected: `{"message":"Specified channel is not exist.","isSuccess":false}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `channel update failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.channel = &pixelaChannelMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdChannelUpdate()
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

func TestChannelDeleteInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.ChannelDeleteInput
	}{
		{
			commandline: "channel delete --id=channel-id",
			expected: pixela.ChannelDeleteInput{
				ID: pixela.String("channel-id"),
			},
		},
		{
			commandline: "channel delete",
			expected:    pixela.ChannelDeleteInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createChannelDeleteInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "ID")
	}
}

func TestChannelDelete(t *testing.T) {
	defer func() { pixelaClient.channel = nil }()
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
				Message:   "Specified channel is not exist.",
				IsSuccess: false,
			},
			occur:    nil,
			expected: `{"message":"Specified channel is not exist.","isSuccess":false}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `channel delete failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.channel = &pixelaChannelMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdChannelDelete()
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
