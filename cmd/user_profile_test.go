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

type pixelaUserProfileMock struct {
	result pixela.Result
	err    error
}

func (p *pixelaUserProfileMock) Update(input *pixela.UserProfileUpdateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaUserProfileMock) URL() string {
	return "https://pixe.la/@pa"
}

func TestUserProfileUpdateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.UserProfileUpdateInput
	}{
		{
			commandline: "profile update --display-name=display-name" +
				" --gravatar-icon-email=gravatar-icon-email" +
				" --title=title" +
				" --timezone=Asia/Tokyo" +
				" --about-url=about-URL" +
				" --contribute-urls=url" +
				" --pinned-graph-id=pinned-graph-id",
			expected: pixela.UserProfileUpdateInput{
				DisplayName:       pixela.String("display-name"),
				GravatarIconEmail: pixela.String("gravatar-icon-email"),
				Title:             pixela.String("title"),
				Timezone:          pixela.String("Asia/Tokyo"),
				AboutURL:          pixela.String("about-URL"),
				ContributeURLs:    []string{"url"},
				PinnedGraphID:     pixela.String("pinned-graph-id"),
			},
		},
		{
			commandline: "profile update",
			expected: pixela.UserProfileUpdateInput{
				ContributeURLs: []string{},
			},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createUserProfileUpdateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.DisplayName), pixela.StringValue(input.DisplayName), "DisplayName")
		assert.EqualValues(t, pixela.StringValue(p.expected.DisplayName), pixela.StringValue(input.DisplayName), "DisplayName")
		assert.EqualValues(t, pixela.StringValue(p.expected.GravatarIconEmail), pixela.StringValue(input.GravatarIconEmail), "GravatarIconEmail")
		assert.EqualValues(t, pixela.StringValue(p.expected.Title), pixela.StringValue(input.Title), "Title")
		assert.EqualValues(t, pixela.StringValue(p.expected.Timezone), pixela.StringValue(input.Timezone), "Timezone")
		assert.EqualValues(t, pixela.StringValue(p.expected.AboutURL), pixela.StringValue(input.AboutURL), "AboutURL")
		assert.EqualValues(t, p.expected.ContributeURLs, input.ContributeURLs, "ContributeURLs")
		assert.EqualValues(t, pixela.StringValue(p.expected.PinnedGraphID), pixela.StringValue(input.PinnedGraphID), "PinnedGraphID")
	}
}

func TestUserProfileUpdate(t *testing.T) {
	defer func() { pixelaClient.profile = nil }()
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
				Message:   "Specified user is not found.",
				IsSuccess: false,
			},
			occur:    nil,
			expected: `{"message":"Specified user is not found.","isSuccess":false}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `user profile update failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.profile = &pixelaUserProfileMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdUserProfileUpdate()
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

func TestUserProfileURL(t *testing.T) {
	defer func() { pixelaClient.profile = nil }()
	pixelaClient.profile = &pixelaUserProfileMock{}
	c := NewCmdUserProfileURL()
	buffer := bytes.NewBuffer([]byte{})
	c.SetOut(buffer)

	err := c.RunE(c, []string{})

	assert.NoError(t, err)
	expected := pixelaClient.profile.URL() + "\n"
	assert.Equal(t, expected, buffer.String())
}
