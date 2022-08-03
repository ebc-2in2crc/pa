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

type pixelaPixelMock struct {
	result   pixela.Result
	err      error
	quantity pixela.Quantity
}

func (p *pixelaPixelMock) Create(input *pixela.PixelCreateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaPixelMock) Increment(input *pixela.PixelIncrementInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaPixelMock) Decrement(input *pixela.PixelDecrementInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaPixelMock) Get(input *pixela.PixelGetInput) (*pixela.Quantity, error) {
	result := &pixela.Quantity{
		Quantity:     p.quantity.Quantity,
		OptionalData: p.quantity.OptionalData,
		Result:       p.result,
	}
	return result, p.err
}

func (p *pixelaPixelMock) Update(input *pixela.PixelUpdateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaPixelMock) Delete(input *pixela.PixelDeleteInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func TestPixelCreateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.PixelCreateInput
	}{
		{
			commandline: "pixel create --graph-id=graph-id --date=20200101 --quantity=5 --optional-data=OD",
			expected: pixela.PixelCreateInput{
				GraphID:      pixela.String("graph-id"),
				Date:         pixela.String("20200101"),
				Quantity:     pixela.String("5"),
				OptionalData: pixela.String("OD"),
			},
		},
		{
			commandline: "pixel create",
			expected:    pixela.PixelCreateInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createPixelCreateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Date), pixela.StringValue(input.Date), "Date")
		assert.EqualValues(t, pixela.StringValue(p.expected.Quantity), pixela.StringValue(input.Quantity), "Quantity")
		assert.EqualValues(t, pixela.StringValue(p.expected.OptionalData), pixela.StringValue(input.OptionalData), "OptionalData")
	}
}

func TestPixelCreate(t *testing.T) {
	defer func() { pixelaClient.pixel = nil }()
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
				Message:    "This date pixel already exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"This date pixel already exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `pixel create failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.pixel = &pixelaPixelMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdPixelCreate()
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

func TestPixelIncrementInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.PixelIncrementInput
	}{
		{
			commandline: "pixel increment --graph-id=graph-id",
			expected: pixela.PixelIncrementInput{
				GraphID: pixela.String("graph-id"),
			},
		},
		{
			commandline: "pixel increment",
			expected:    pixela.PixelIncrementInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createPixelIncrementInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
	}
}

func TestPixelIncrement(t *testing.T) {
	defer func() { pixelaClient.pixel = nil }()
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
				Message:    "Specified graphID not exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"Specified graphID not exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `pixel increment failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.pixel = &pixelaPixelMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdPixelIncrement()
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

func TestPixelDecrementInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.PixelDecrementInput
	}{
		{
			commandline: "pixel decrement --graph-id=graph-id",
			expected: pixela.PixelDecrementInput{
				GraphID: pixela.String("graph-id"),
			},
		},
		{
			commandline: "pixel decrement",
			expected:    pixela.PixelDecrementInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createPixelDecrementInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
	}
}

func TestPixelDecrement(t *testing.T) {
	defer func() { pixelaClient.pixel = nil }()
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
				Message:    "Specified graphID not exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"Specified graphID not exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `pixel decrement failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.pixel = &pixelaPixelMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdPixelDecrement()
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

func TestPixelGetInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.PixelGetInput
	}{
		{
			commandline: "pixel get --graph-id=graph-id --date=20200101",
			expected: pixela.PixelGetInput{
				GraphID: pixela.String("graph-id"),
				Date:    pixela.String("20200101"),
			},
		},
		{
			commandline: "pixel get",
			expected:    pixela.PixelGetInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createPixelGetInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Date), pixela.StringValue(input.Date), "Date")
	}
}

func TestPixelGet(t *testing.T) {
	defer func() { pixelaClient.pixel = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		expected string
		quantity pixela.Quantity
	}{
		{
			Result: pixela.Result{
				Message:    "Success.",
				IsSuccess:  true,
				StatusCode: http.StatusOK,
			},
			occur: nil,
			quantity: pixela.Quantity{
				Quantity:     "1",
				OptionalData: "OD",
			},
			expected: `{"quantity":"1","optionalData":"OD"}` + "\n",
		},
		{
			Result: pixela.Result{
				Message:    "Specified pixel not found.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"Specified pixel not found.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `pixel get failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.pixel = &pixelaPixelMock{
			result:   v.Result,
			err:      v.occur,
			quantity: v.quantity,
		}
		c := NewCmdPixelGet()
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

func TestPixelUpdateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.PixelUpdateInput
	}{
		{
			commandline: "pixel update --graph-id=graph-id --date=20200101 --quantity=5 --optional-data=OD",
			expected: pixela.PixelUpdateInput{
				GraphID:      pixela.String("graph-id"),
				Date:         pixela.String("20200101"),
				Quantity:     pixela.String("5"),
				OptionalData: pixela.String("OD"),
			},
		},
		{
			commandline: "pixel update",
			expected:    pixela.PixelUpdateInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createPixelUpdateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Date), pixela.StringValue(input.Date), "Date")
		assert.EqualValues(t, pixela.StringValue(p.expected.Quantity), pixela.StringValue(input.Quantity), "Quantity")
		assert.EqualValues(t, pixela.StringValue(p.expected.OptionalData), pixela.StringValue(input.OptionalData), "OptionalData")
	}
}

func TestPixelUpdate(t *testing.T) {
	defer func() { pixelaClient.pixel = nil }()
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
				Message:    "Specified graphID not found.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"Specified graphID not found.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `pixel update failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.pixel = &pixelaPixelMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdPixelUpdate()
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

func TestPixelDeleteInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.PixelDeleteInput
	}{
		{
			commandline: "pixel delete --graph-id=graph-id --date=20200101",
			expected: pixela.PixelDeleteInput{
				GraphID: pixela.String("graph-id"),
				Date:    pixela.String("20200101"),
			},
		},
		{
			commandline: "pixel update",
			expected:    pixela.PixelDeleteInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(ioutil.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		cmd.Execute()

		input := createPixelDeleteInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.GraphID), pixela.StringValue(input.GraphID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Date), pixela.StringValue(input.Date), "Date")
	}
}

func TestPixelDelete(t *testing.T) {
	defer func() { pixelaClient.pixel = nil }()
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
				Message:    "Specified pixel not found.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"Specified pixel not found.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `pixel delete failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.pixel = &pixelaPixelMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdPixelDelete()
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
