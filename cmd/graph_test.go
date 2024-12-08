package cmd

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/stretchr/testify/assert"
)

type pixelaGraphMock struct {
	result      pixela.Result
	err         error
	svg         string
	stats       pixela.Stats
	pixel       pixela.PixelWithBody
	pixels      pixela.Pixels
	definitions pixela.GraphDefinitions
	definition  pixela.GraphDefinition
}

func (p *pixelaGraphMock) Create(input *pixela.GraphCreateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaGraphMock) GetAll() (*pixela.GraphDefinitions, error) {
	return &p.definitions, p.err
}

func (p *pixelaGraphMock) Get(input *pixela.GraphGetInput) (*pixela.GraphDefinition, error) {
	return &p.definition, p.err
}

func (p *pixelaGraphMock) GetSVG(input *pixela.GraphGetSVGInput) (string, error) {
	return p.svg, p.err
}

func (p *pixelaGraphMock) URL(input *pixela.GraphURLInput) string {
	return "https://pixe.la/v1/users/pa/graphs/graph-id.html"
}

func (p *pixelaGraphMock) GraphsURL() string {
	return "https://pixe.la/v1/users/pa/graphs.html"
}

func (p *pixelaGraphMock) Stats(input *pixela.GraphStatsInput) (*pixela.Stats, error) {
	return &p.stats, p.err
}

func (p *pixelaGraphMock) Update(input *pixela.GraphUpdateInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaGraphMock) GetPixelDates(input *pixela.GraphGetPixelDatesInput) (*pixela.Pixels, error) {
	return &p.pixels, p.err
}

func (p *pixelaGraphMock) Delete(input *pixela.GraphDeleteInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaGraphMock) Stopwatch(input *pixela.GraphStopwatchInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaGraphMock) Add(input *pixela.GraphAddInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaGraphMock) Subtract(input *pixela.GraphSubtractInput) (*pixela.Result, error) {
	return &p.result, p.err
}

func (p *pixelaGraphMock) GetLatestPixel(input *pixela.GraphGetLatestPixelInput) (*pixela.GraphPixel, error) {
	result := &pixela.GraphPixel{
		Date:         p.pixel.Date,
		Quantity:     p.pixel.Quantity,
		OptionalData: p.pixel.OptionalData,
		Result:       p.result,
	}
	return result, p.err
}

func TestGraphCreateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphCreateInput
	}{
		{
			commandline: "graph create --id=graph-id --name=graph-name --unit=times --type=int --color=sora" +
				" --timezone=Asia/Tokyo --self-sufficient=none --secret --publish-optional-data --start-on-monday",
			expected: pixela.GraphCreateInput{
				ID:                  pixela.String("graph-id"),
				Name:                pixela.String("graph-name"),
				Unit:                pixela.String("times"),
				Type:                pixela.String("int"),
				Color:               pixela.String("sora"),
				TimeZone:            pixela.String("Asia/Tokyo"),
				SelfSufficient:      pixela.String("none"),
				IsSecret:            pixela.Bool(true),
				PublishOptionalData: pixela.Bool(true),
				StartOnMonday:       pixela.Bool(true),
			},
		},
		{
			commandline: "graph create",
			expected:    pixela.GraphCreateInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphCreateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Name), pixela.StringValue(input.Name), "GraphName")
		assert.EqualValues(t, pixela.StringValue(p.expected.Unit), pixela.StringValue(input.Unit), "Unit")
		assert.EqualValues(t, pixela.StringValue(p.expected.Type), pixela.StringValue(input.Type), "Type")
		assert.EqualValues(t, pixela.StringValue(p.expected.Color), pixela.StringValue(input.Color), "Color")
		assert.EqualValues(t, pixela.StringValue(p.expected.TimeZone), pixela.StringValue(input.TimeZone), "TimeZone")
		assert.EqualValues(t, pixela.StringValue(p.expected.SelfSufficient), pixela.StringValue(input.SelfSufficient), "SelfSufficient")
		assert.EqualValues(t, pixela.BoolValue(p.expected.IsSecret), pixela.BoolValue(input.IsSecret), "IsSecret")
		assert.EqualValues(t, pixela.BoolValue(p.expected.PublishOptionalData), pixela.BoolValue(input.PublishOptionalData), "PublishOptionalData")
		assert.EqualValues(t, pixela.BoolValue(p.expected.StartOnMonday), pixela.BoolValue(input.StartOnMonday), "StartOnMonday")
	}
}

func TestGraphCreate(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
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
				Message:    "This graphID already exist.",
				IsSuccess:  false,
				StatusCode: http.StatusBadRequest,
			},
			occur:    nil,
			expected: `{"message":"This graphID already exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `graph create failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdGraphCreate()
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

func TestGraphGetAll(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
	params := []struct {
		definitions pixela.GraphDefinitions
		occur       error
		expected    string
	}{
		{
			definitions: pixela.GraphDefinitions{
				Result: pixela.Result{
					Message:    "Success.",
					IsSuccess:  true,
					StatusCode: http.StatusOK,
				},
				Graphs: []pixela.GraphDefinition{
					{
						ID:                  "graph-id",
						Name:                "graph-name",
						Unit:                "count",
						Type:                "int",
						Color:               "ichou",
						TimeZone:            "Asia/Tokyo",
						PurgeCacheURLs:      nil,
						SelfSufficient:      "increment",
						IsSecret:            true,
						PublishOptionalData: true,
					},
				},
			},
			occur: nil,
			expected: `{"graphs":[{"id":"graph-id","name":"graph-name","unit":"count","type":"int","color":"ichou",` +
				`"timezone":"Asia/Tokyo","purgeCacheURLs":null,"selfSufficient":"increment","isSecret":true,` +
				`"publishOptionalData":true}]}` + "\n",
		},
		{
			definitions: pixela.GraphDefinitions{
				Result: pixela.Result{
					Message:    "User does not exists.",
					IsSuccess:  false,
					StatusCode: http.StatusBadRequest,
				},
			},
			occur:    nil,
			expected: `{"message":"User does not exists.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			occur:    errors.New("some error occur"),
			expected: `graph get all failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			definitions: v.definitions,
			err:         v.occur,
		}
		c := NewCmdGraphGetAll()
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

func TestGraphGetInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphGetInput
	}{
		{
			commandline: "graph get --id=graph-id",
			expected:    pixela.GraphGetInput{ID: pixela.String("graph-id")},
		},
		{
			commandline: "graph get",
			expected:    pixela.GraphGetInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphGetInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
	}
}

func TestGraphGet(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
	params := []struct {
		definition pixela.GraphDefinition
		occur      error
		expected   string
	}{
		{
			definition: pixela.GraphDefinition{
				Result: pixela.Result{
					Message:    "Success.",
					IsSuccess:  true,
					StatusCode: http.StatusOK,
				},
				ID:                  "graph-id",
				Name:                "graph-name",
				Unit:                "count",
				Type:                "int",
				Color:               "ichou",
				TimeZone:            "Asia/Tokyo",
				PurgeCacheURLs:      nil,
				SelfSufficient:      "increment",
				IsSecret:            true,
				PublishOptionalData: true,
			},
			occur: nil,
			expected: `{"id":"graph-id","name":"graph-name","unit":"count","type":"int","color":"ichou",` +
				`"timezone":"Asia/Tokyo","purgeCacheURLs":null,"selfSufficient":"increment","isSecret":true,` +
				`"publishOptionalData":true}` + "\n",
		},
		{
			definition: pixela.GraphDefinition{
				Result: pixela.Result{
					Message:    "User does not exists.",
					IsSuccess:  false,
					StatusCode: http.StatusBadRequest,
				},
			},
			occur:    nil,
			expected: `{"message":"User does not exists.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			occur:    errors.New("some error occur"),
			expected: `graph get failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			definition: v.definition,
			err:        v.occur,
		}
		c := NewCmdGraphGet()
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

func TestGraphGetSVGInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphGetSVGInput
	}{
		{
			commandline: "graph svg --id=graph-id --date=20200101 --mode=badge --appearance=dark",
			expected: pixela.GraphGetSVGInput{
				ID:         pixela.String("graph-id"),
				Date:       pixela.String("20200101"),
				Mode:       pixela.String("badge"),
				Appearance: pixela.String("dark"),
			},
		},
		{
			commandline: "graph svg",
			expected:    pixela.GraphGetSVGInput{},
		},
	}

	for _, p := range params {
		pixelaClient.graph = &pixelaGraphMock{
			err: nil,
			svg: "SVG",
		}
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphGetSVGInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Date), pixela.StringValue(input.Date), "Date")
		assert.EqualValues(t, pixela.StringValue(p.expected.Mode), pixela.StringValue(input.Mode), "Mode")
		assert.EqualValues(t, pixela.StringValue(p.expected.Appearance), pixela.StringValue(input.Appearance), "Appearance")
	}
}

func TestGraphGetSVG(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
	params := []struct {
		occur    error
		expected string
		svg      string
	}{
		{
			occur:    nil,
			expected: "SVG\n",
			svg:      "SVG",
		},
		{
			occur:    errors.New(`some error occur: "{"message": "error"}`),
			expected: `{"message": "error"}` + "\n",
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			err: v.occur,
			svg: v.svg,
		}
		c := NewCmdGraphGetSVG()
		buffer := bytes.NewBuffer([]byte{})
		c.SetOut(buffer)

		_ = c.RunE(c, []string{})

		assert.Equal(t, v.expected, buffer.String())
	}
}

func TestGraphURLInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphURLInput
	}{
		{
			commandline: "graph detail --id=graph-id --mode=simple",
			expected: pixela.GraphURLInput{
				ID:   pixela.String("graph-id"),
				Mode: pixela.String("simple"),
			},
		},
		{
			commandline: "graph detail",
			expected:    pixela.GraphURLInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphURLInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Mode), pixela.StringValue(input.Mode), "Mode")
	}
}

func TestGraphURL(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
	pixelaClient.graph = &pixelaGraphMock{}
	c := NewCmdGraphURL()
	buffer := bytes.NewBuffer([]byte{})
	c.SetOut(buffer)

	err := c.RunE(c, []string{})

	assert.NoError(t, err)
	expected := pixelaClient.graph.URL(nil) + "\n"
	assert.Equal(t, expected, buffer.String())
}

func TestGraphStatsInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphStatsInput
	}{
		{
			commandline: "graph stats --id=graph-id",
			expected: pixela.GraphStatsInput{
				ID: pixela.String("graph-id"),
			},
		},
		{
			commandline: "graph stats",
			expected:    pixela.GraphStatsInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphStatsInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
	}
}

func TestGraphStats(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
	params := []struct {
		stats    pixela.Stats
		occur    error
		expected string
	}{
		{
			stats: pixela.Stats{
				Result: pixela.Result{
					Message:    "Success.",
					IsSuccess:  true,
					StatusCode: http.StatusOK,
				},
				TotalPixelsCount:  1,
				MaxQuantity:       2,
				MaxDate:           "2023-09-01",
				MinQuantity:       3,
				MinDate:           "2023-09-02",
				TotalQuantity:     4,
				AvgQuantity:       5,
				TodaysQuantity:    6,
				YesterdayQuantity: 66,
			},
			occur: nil,
			expected: `{"totalPixelsCount":1,"maxQuantity":2,"maxDate":"2023-09-01","minQuantity":3,"minDate":"2023-09-02","totalQuantity":4,` +
				`"avgQuantity":5,"todaysQuantity":6,"yesterdayQuantity":66}` + "\n",
		},
		{
			stats: pixela.Stats{
				Result: pixela.Result{
					Message:    "Specified graphID not exist.",
					IsSuccess:  false,
					StatusCode: http.StatusBadRequest,
				},
			},
			occur:    nil,
			expected: `{"message":"Specified graphID not exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			stats:    pixela.Stats{},
			occur:    errors.New("some error occur"),
			expected: `graph stats failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			stats: v.stats,
			err:   v.occur,
		}
		c := NewCmdGraphStats()
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

func TestGraphUpdateInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphUpdateInput
	}{
		{
			commandline: "graph update --id=graph-id --name=graph-name --unit=times --color=sora" +
				" --timezone=Asia/Tokyo --purge-cache-urls=url --self-sufficient=none --secret --publish-optional-data --start-on-monday",
			expected: pixela.GraphUpdateInput{
				ID:                  pixela.String("graph-id"),
				Name:                pixela.String("graph-name"),
				Unit:                pixela.String("times"),
				Color:               pixela.String("sora"),
				TimeZone:            pixela.String("Asia/Tokyo"),
				PurgeCacheURLs:      []string{"url"},
				SelfSufficient:      pixela.String("none"),
				IsSecret:            pixela.Bool(true),
				PublishOptionalData: pixela.Bool(true),
				StartOnMonday:       pixela.Bool(true),
			},
		},
		{
			commandline: "graph update --id=graph-id --publish --hide-optional-data",
			expected: pixela.GraphUpdateInput{
				ID:                  pixela.String("graph-id"),
				PurgeCacheURLs:      []string{},
				IsSecret:            pixela.Bool(false),
				PublishOptionalData: pixela.Bool(false),
			},
		},
		{
			commandline: "graph update",
			expected: pixela.GraphUpdateInput{
				PurgeCacheURLs:      []string{},
				IsSecret:            nil,
				PublishOptionalData: nil,
			},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphUpdateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Name), pixela.StringValue(input.Name), "GraphName")
		assert.EqualValues(t, pixela.StringValue(p.expected.Unit), pixela.StringValue(input.Unit), "Unit")
		assert.EqualValues(t, pixela.StringValue(p.expected.Color), pixela.StringValue(input.Color), "Color")
		assert.EqualValues(t, pixela.StringValue(p.expected.TimeZone), pixela.StringValue(input.TimeZone), "TimeZone")
		assert.EqualValues(t, p.expected.PurgeCacheURLs, input.PurgeCacheURLs, "PurgeCacheURLs")
		assert.EqualValues(t, pixela.StringValue(p.expected.SelfSufficient), pixela.StringValue(input.SelfSufficient), "SelfSufficient")
		if p.expected.IsSecret == nil {
			assert.Nil(t, input.IsSecret, "IsSecret")
		} else {
			assert.NotNil(t, input.IsSecret, "IsSecret")
			assert.EqualValues(t, pixela.BoolValue(p.expected.IsSecret), pixela.BoolValue(input.IsSecret), "IsSecret")
		}
		if p.expected.PublishOptionalData == nil {
			assert.Nil(t, input.PublishOptionalData, "PublishOptionalData")
		} else {
			assert.NotNil(t, input.PublishOptionalData, "PublishOptionalData")
			assert.EqualValues(t, pixela.BoolValue(p.expected.PublishOptionalData), pixela.BoolValue(input.PublishOptionalData), "PublishOptionalData")
		}
		if p.expected.StartOnMonday == nil {
			assert.Nil(t, input.StartOnMonday, "StartOnMonday")
		} else {
			assert.NotNil(t, input.StartOnMonday, "StartOnMonday")
			assert.EqualValues(t, pixela.BoolValue(p.expected.StartOnMonday), pixela.BoolValue(input.StartOnMonday), "StartOnMonday")
		}
	}
}

func TestGraphUpdate(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
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
			expected: `graph update failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdGraphUpdate()
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

func TestGraphDeleteInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphDeleteInput
	}{
		{
			commandline: "graph delete --id=graph-id",
			expected: pixela.GraphDeleteInput{
				ID: pixela.String("graph-id"),
			},
		},
		{
			commandline: "graph delete",
			expected:    pixela.GraphDeleteInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphDeleteInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
	}
}

func TestGraphDelete(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
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
			expected: `graph delete failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdGraphDelete()
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

func TestGraphPixelsInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphGetPixelDatesInput
	}{
		{
			commandline: "graph pixels --id=graph-id --from=20200101 --to=20200130 --with-body",
			expected: pixela.GraphGetPixelDatesInput{
				ID:       pixela.String("graph-id"),
				From:     pixela.String("20200101"),
				To:       pixela.String("20200130"),
				WithBody: pixela.Bool(true),
			},
		},
		{
			commandline: "graph pixels",
			expected:    pixela.GraphGetPixelDatesInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphGetPixelDatesInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.From), pixela.StringValue(input.From), "From")
		assert.EqualValues(t, pixela.StringValue(p.expected.To), pixela.StringValue(input.To), "To")
		assert.EqualValues(t, pixela.BoolValue(p.expected.WithBody), pixela.BoolValue(input.WithBody), "WithBody")
	}
}

func TestGraphPixels(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
	params := []struct {
		withBody bool
		pixels   pixela.Pixels
		occur    error
		expected string
	}{
		{
			pixels: pixela.Pixels{
				Result: pixela.Result{
					Message:    "Success.",
					IsSuccess:  true,
					StatusCode: http.StatusOK,
				},
				Pixels: []string{"20200101"},
			},
			occur:    nil,
			expected: `{"pixels":["20200101"]}` + "\n",
		},
		{
			withBody: true,
			pixels: pixela.Pixels{
				Result: pixela.Result{
					Message:    "Success.",
					IsSuccess:  true,
					StatusCode: http.StatusOK,
				},
				Pixels: []pixela.PixelWithBody{
					{
						Date:         "20200101",
						Quantity:     "5",
						OptionalData: "OD",
					},
				},
			},
			occur:    nil,
			expected: `{"pixels":[{"date":"20200101","quantity":"5","optionalData":"OD"}]}` + "\n",
		},
		{
			pixels: pixela.Pixels{
				Result: pixela.Result{
					Message:    "Specified graphID not exist.",
					IsSuccess:  false,
					StatusCode: http.StatusBadRequest,
				},
			},
			occur:    nil,
			expected: `{"message":"Specified graphID not exist.","isSuccess":false,"isRejected":false,"statusCode":400}` + "\n",
		},
		{
			occur:    errors.New("some error occur"),
			expected: `graph get pixel dates failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			pixels: v.pixels,
			err:    v.occur,
		}
		c := NewCmdGraphGetPixelDates()
		buffer := bytes.NewBuffer([]byte{})
		c.SetOut(buffer)
		assert.NoError(t, c.Flags().Set("with-body", strconv.FormatBool(v.withBody)))

		err := c.RunE(c, []string{})

		if v.occur == nil {
			assert.Equal(t, v.expected, buffer.String())
		} else {
			assert.Contains(t, err.Error(), v.expected)
		}
	}
}

func TestGraphStopwatchInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphStopwatchInput
	}{
		{
			commandline: "graph stopwatch --id=graph-id",
			expected: pixela.GraphStopwatchInput{
				ID: pixela.String("graph-id"),
			},
		},
		{
			commandline: "graph stopwatch",
			expected:    pixela.GraphStopwatchInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphStopwatchInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
	}
}

func TestGraphStopwatch(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
	params := []struct {
		Result   pixela.Result
		occur    error
		expected string
	}{
		{
			Result: pixela.Result{
				Message:    "Stopwatch start successful.",
				IsSuccess:  true,
				StatusCode: http.StatusOK,
			},
			occur:    nil,
			expected: `{"message":"Stopwatch start successful.","isSuccess":true,"isRejected":false,"statusCode":200}` + "\n",
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
			expected: `graph stopwatch failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdGraphStopwatch()
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

func TestGraphAddInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphAddInput
	}{
		{
			commandline: "graph add --id=graph-id --quantity=1",
			expected: pixela.GraphAddInput{
				ID:       pixela.String("graph-id"),
				Quantity: pixela.String("1"),
			},
		},
		{
			commandline: "graph add",
			expected:    pixela.GraphAddInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphAddInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Quantity), pixela.StringValue(input.Quantity), "Quantity")
	}
}

func TestGraphAdd(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
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
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `graph add failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdGraphAdd()
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

func TestGraphSubtractInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphSubtractInput
	}{
		{
			commandline: "graph subtract --id=graph-id --quantity=1",
			expected: pixela.GraphSubtractInput{
				ID:       pixela.String("graph-id"),
				Quantity: pixela.String("1"),
			},
		},
		{
			commandline: "graph subtract",
			expected:    pixela.GraphSubtractInput{},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphSubtractInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
		assert.EqualValues(t, pixela.StringValue(p.expected.Quantity), pixela.StringValue(input.Quantity), "Quantity")
	}
}

func TestGraphSubtract(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
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
			Result:   pixela.Result{},
			occur:    errors.New("some error occur"),
			expected: `graph subtract failed:`,
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			result: v.Result,
			err:    v.occur,
		}
		c := NewCmdGraphSubtract()
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

func TestGraphGetLatestPixelInput(t *testing.T) {
	params := []struct {
		commandline string
		expected    pixela.GraphGetLatestPixelInput
	}{
		{
			commandline: "graph get-latest-pixel --id=graph-id",
			expected: pixela.GraphGetLatestPixelInput{
				ID: pixela.String("graph-id"),
			},
		},
	}

	for _, p := range params {
		cmd := NewCmdRoot()
		cmd.SetOut(io.Discard)
		args := strings.Split(p.commandline, " ")
		cmd.SetArgs(args)
		_ = cmd.Execute()

		input := createGraphCreateInput()

		assert.EqualValues(t, pixela.StringValue(p.expected.ID), pixela.StringValue(input.ID), "GraphID")
	}
}

func TestGraphGetLatestPixel(t *testing.T) {
	defer func() { pixelaClient.graph = nil }()
	params := []struct {
		result   pixela.Result
		pixel    pixela.PixelWithBody
		occur    error
		expected string
	}{
		{
			result: pixela.Result{
				Message:    "Success.",
				IsSuccess:  true,
				StatusCode: http.StatusOK,
			},
			pixel: pixela.PixelWithBody{
				Date:         "20241205",
				Quantity:     "5",
				OptionalData: "OD",
			},
			occur:    nil,
			expected: `{"date":"20241205","quantity":"5","optionalData":"OD"}` + "\n",
		},
	}

	for _, v := range params {
		pixelaClient.graph = &pixelaGraphMock{
			result: v.result,
			pixel:  v.pixel,
			err:    v.occur,
		}
		c := NewCmdGraphGetLatestPixel()
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
