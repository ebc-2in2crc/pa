package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/spf13/cobra"
)

var graphOptions = &struct {
	ID                  string
	Name                string
	Unit                string
	Type                string
	Color               string
	TimeZone            string
	PurgeCacheURLs      []string
	SelfSufficient      string
	IsSecret            bool
	IsPublish           bool
	PublishOptionalData bool
	HideOptionalData    bool
	Date                string
	Mode                string
	Appearance          string
	DeleteMe            bool
	From                string
	To                  string
	WithBody            bool
}{}

// NewCmdGraph creates a graph command.
func NewCmdGraph() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Graph",
		Args:  cobra.NoArgs,
		RunE:  showHelp,
	}

	cmd.AddCommand(NewCmdGraphCreate())
	cmd.AddCommand(NewCmdGraphGetAll())
	cmd.AddCommand(NewCmdGraphGet())
	cmd.AddCommand(NewCmdGraphGetSVG())
	cmd.AddCommand(NewCmdGraphURL())
	cmd.AddCommand(NewCmdGraphStats())
	cmd.AddCommand(NewCmdGraphUpdate())
	cmd.AddCommand(NewCmdGraphDelete())
	cmd.AddCommand(NewCmdGraphGetPixelDates())
	cmd.AddCommand(NewCmdGraphStopwatch())

	return cmd
}

// NewCmdGraphCreate creates a create graph command.
func NewCmdGraphCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Graph",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createGraphCreateInput()
			result, err := pixelaClient.Graph().Create(input)
			if err != nil {
				return fmt.Errorf("graph create failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal graph create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&graphOptions.ID, "id", "", "ID for identifying the pixelation graph")
	cmd.Flags().StringVar(&graphOptions.Name, "name", "", "The name of the pixelation graph")
	cmd.Flags().StringVar(&graphOptions.Unit, "unit", "", "A Unit of the quantity recorded in the pixelation graph")
	cmd.Flags().StringVar(&graphOptions.Type, "type", "", "The type of quantity to be handled in the graph")
	cmd.Flags().StringVar(&graphOptions.Color, "color", "", "Defines the display color of the pixel in the pixelation graph")
	cmd.Flags().StringVar(&graphOptions.TimeZone, "timezone", "", "The timezone for handling this graph")
	cmd.Flags().StringVar(&graphOptions.SelfSufficient, "self-sufficient", "", "See: https://docs.pixe.la/entry/post-graph")
	cmd.Flags().BoolVar(&graphOptions.IsSecret, "secret", false, "The Graph not displayed on the graph list page")
	cmd.Flags().BoolVar(&graphOptions.PublishOptionalData, "publish-optional-data", false, "Each pixel's optionalData will be added to the generated SVG data")

	return cmd
}

func createGraphCreateInput() *pixela.GraphCreateInput {
	return &pixela.GraphCreateInput{
		ID:                  getStringPtr(graphOptions.ID),
		Name:                getStringPtr(graphOptions.Name),
		Unit:                getStringPtr(graphOptions.Unit),
		Type:                getStringPtr(graphOptions.Type),
		Color:               getStringPtr(graphOptions.Color),
		TimeZone:            getStringPtr(graphOptions.TimeZone),
		SelfSufficient:      getStringPtr(graphOptions.SelfSufficient),
		IsSecret:            getBoolPtr(graphOptions.IsSecret),
		PublishOptionalData: getBoolPtr(graphOptions.PublishOptionalData),
	}
}

// NewCmdGraphGetAll creates a get all graph command.
func NewCmdGraphGetAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all",
		Short: "Get Graph definitions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			definitions, err := pixelaClient.Graph().GetAll()
			if err != nil {
				return fmt.Errorf("graph get all failed: %w", err)
			}
			if definitions.IsSuccess == false {
				s, err := marshalResult(&definitions.Result)
				if err != nil {
					return fmt.Errorf("marshal graph get all result failed: %w", err)
				}
				cmd.Printf("%s\n", s)
				return ErrNeglect
			}

			defs := make([]graphDefinition, len(definitions.Graphs))
			for i, v := range definitions.Graphs {
				defs[i] = gToG(&v)
			}

			b, err := json.Marshal(&graphDefinitions{Graphs: defs})
			if err != nil {
				return fmt.Errorf("marshal graph get all definitions failed: %w", err)
			}
			cmd.Printf("%s\n", string(b))

			return nil
		},
	}

	return cmd
}

type graphDefinitions struct {
	Graphs []graphDefinition `json:"graphs"`
}

type graphDefinition struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Unit                string   `json:"unit"`
	Type                string   `json:"type"`
	Color               string   `json:"color"`
	TimeZone            string   `json:"timezone"`
	PurgeCacheURLs      []string `json:"purgeCacheURLs"`
	SelfSufficient      string   `json:"selfSufficient"`
	IsSecret            bool     `json:"isSecret"`
	PublishOptionalData bool     `json:"publishOptionalData"`
}

func gToG(g *pixela.GraphDefinition) graphDefinition {
	return graphDefinition{
		ID:                  g.ID,
		Name:                g.Name,
		Unit:                g.Unit,
		Type:                g.Type,
		Color:               g.Color,
		TimeZone:            g.TimeZone,
		PurgeCacheURLs:      g.PurgeCacheURLs,
		SelfSufficient:      g.SelfSufficient,
		IsSecret:            g.IsSecret,
		PublishOptionalData: g.PublishOptionalData,
	}
}

// NewCmdGraphGet creates a get graph command.
func NewCmdGraphGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Graph definition",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createGraphGetInput()
			result, err := pixelaClient.Graph().Get(input)
			if err != nil {
				return fmt.Errorf("graph get failed: %w", err)
			}
			if result.IsSuccess == false {
				s, err := marshalResult(&result.Result)
				if err != nil {
					return fmt.Errorf("marshal graph get result failed: %w", err)
				}
				cmd.Printf("%s\n", s)
				return ErrNeglect
			}

			g := gToG(result)
			b, err := json.Marshal(g)
			if err != nil {
				return fmt.Errorf("marshal graph get definition failed: %w", err)
			}
			cmd.Printf("%s\n", string(b))

			return nil
		},
	}

	cmd.Flags().StringVar(&graphOptions.ID, "id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("id")

	return cmd
}

func createGraphGetInput() *pixela.GraphGetInput {
	return &pixela.GraphGetInput{
		ID: getStringPtr(graphOptions.ID),
	}
}

// NewCmdGraphGetSVG creates a get graph svg command.
func NewCmdGraphGetSVG() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "svg",
		Short: "Get the Graph in SVG format diagram",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createGraphGetSVGInput()
			result, err := pixelaClient.Graph().GetSVG(input)
			if err != nil {
				e := err.Error()
				s := e[strings.LastIndex(e, `{"message"`):]
				cmd.Printf("%s\n", s)
				return ErrNeglect
			}
			cmd.Printf("%s\n", result)

			return nil
		},
	}

	cmd.Flags().StringVar(&graphOptions.ID, "id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&graphOptions.Date, "date", "", "Create a pixelation graph dating back to the past with that day as the start date")
	cmd.Flags().StringVar(&graphOptions.Mode, "mode", "", "The Graph display mode")
	cmd.Flags().StringVar(&graphOptions.Appearance, "appearance", "", "The graph appearance mode")

	return cmd
}

func createGraphGetSVGInput() *pixela.GraphGetSVGInput {
	return &pixela.GraphGetSVGInput{
		ID:         getStringPtr(graphOptions.ID),
		Date:       getStringPtr(graphOptions.Date),
		Mode:       getStringPtr(graphOptions.Mode),
		Appearance: getStringPtr(graphOptions.Appearance),
	}
}

// NewCmdGraphURL creates a graph URL command.
func NewCmdGraphURL() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Get Graph detail URL",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createGraphURLInput()
			url := pixelaClient.Graph().URL(input)
			cmd.Printf("%s\n", url)

			return nil
		},
	}

	cmd.Flags().StringVar(&graphOptions.ID, "id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&graphOptions.Mode, "mode", "", "The graph html page mode")

	return cmd
}

func createGraphURLInput() *pixela.GraphURLInput {
	return &pixela.GraphURLInput{
		ID:   getStringPtr(graphOptions.ID),
		Mode: getStringPtr(graphOptions.Mode),
	}
}

// NewCmdGraphStats creates a graphs stats command.
func NewCmdGraphStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Get various statistics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createGraphStatsInput()
			stats, err := pixelaClient.Graph().Stats(input)
			if err != nil {
				return fmt.Errorf("graph stats failed: %w", err)
			}
			if stats.IsSuccess == false {
				s, err := marshalResult(&stats.Result)
				if err != nil {
					return fmt.Errorf("marshal graph stats result failed: %w", err)
				}
				cmd.Printf("%s\n", s)
				return ErrNeglect
			}

			b, err := json.Marshal(&graphStats{
				TotalPixelsCount: stats.TotalPixelsCount,
				MaxQuantity:      stats.MaxQuantity,
				MinQuantity:      stats.MinQuantity,
				TotalQuantity:    stats.TotalQuantity,
				AvgQuantity:      stats.AvgQuantity,
				TodaysQuantity:   stats.TodaysQuantity,
			})
			if err != nil {
				return fmt.Errorf("marshal graph stats failed: %w", err)
			}
			cmd.Printf("%s\n", string(b))

			return nil
		},
	}

	cmd.Flags().StringVar(&graphOptions.ID, "id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("id")

	return cmd
}

func createGraphStatsInput() *pixela.GraphStatsInput {
	return &pixela.GraphStatsInput{
		ID: getStringPtr(graphOptions.ID),
	}
}

type graphStats struct {
	TotalPixelsCount int     `json:"totalPixelsCount"`
	MaxQuantity      int     `json:"maxQuantity"`
	MinQuantity      int     `json:"minQuantity"`
	TotalQuantity    int     `json:"totalQuantity"`
	AvgQuantity      float64 `json:"avgQuantity"`
	TodaysQuantity   int     `json:"todaysQuantity"`
}

// NewCmdGraphUpdate creates a update graphs command.
func NewCmdGraphUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Graph",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createGraphUpdateInput()
			result, err := pixelaClient.Graph().Update(input)
			if err != nil {
				return fmt.Errorf("graph update failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal graph update result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&graphOptions.ID, "id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&graphOptions.Name, "name", "", "The name of the pixelation graph")
	cmd.Flags().StringVar(&graphOptions.Unit, "unit", "", "A Unit of the quantity recorded in the pixelation graph")
	cmd.Flags().StringVar(&graphOptions.Color, "color", "", "Defines the display color of the pixel in the pixelation graph")
	cmd.Flags().StringVar(&graphOptions.TimeZone, "timezone", "", "The timezone for handling this graph")
	cmd.Flags().StringSliceVar(&graphOptions.PurgeCacheURLs, "purge-cache-urls", []string{}, "URL to send the purge request to purge the cache when the graph is updated")
	cmd.Flags().StringVar(&graphOptions.SelfSufficient, "self-sufficient", "", "See: https://docs.pixe.la/entry/put-graph")
	cmd.Flags().BoolVar(&graphOptions.IsSecret, "secret", false, "The Graph not displayed on the graph list page")
	cmd.Flags().BoolVar(&graphOptions.IsPublish, "publish", false, "The Graph displayed on the graph list page")
	cmd.Flags().BoolVar(&graphOptions.PublishOptionalData, "publish-optional-data", false, "Each pixel's optionalData will be added to the generated SVG data")
	cmd.Flags().BoolVar(&graphOptions.HideOptionalData, "hide-optional-data", false, "Each pixel's optionalData will not be added to the generated SVG data")

	return cmd
}

func createGraphUpdateInput() *pixela.GraphUpdateInput {
	var secret *bool
	if graphOptions.IsPublish {
		secret = pixela.Bool(false)
	}
	if graphOptions.IsSecret {
		secret = pixela.Bool(true)
	}

	var publishOptionalData *bool
	if graphOptions.PublishOptionalData {
		publishOptionalData = pixela.Bool(true)
	}
	if graphOptions.HideOptionalData {
		publishOptionalData = pixela.Bool(false)
	}

	return &pixela.GraphUpdateInput{
		ID:                  getStringPtr(graphOptions.ID),
		Name:                getStringPtr(graphOptions.Name),
		Unit:                getStringPtr(graphOptions.Unit),
		Color:               getStringPtr(graphOptions.Color),
		TimeZone:            getStringPtr(graphOptions.TimeZone),
		PurgeCacheURLs:      graphOptions.PurgeCacheURLs,
		SelfSufficient:      getStringPtr(graphOptions.SelfSufficient),
		IsSecret:            secret,
		PublishOptionalData: publishOptionalData,
	}
}

// NewCmdGraphDelete creates a delete graphs command.
func NewCmdGraphDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a Graph",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if graphOptions.DeleteMe == false {
				cmd.Println("Specify the '--delete-me' flag to confirm the deletion.")
				return nil
			}

			input := createGraphDeleteInput()
			result, err := pixelaClient.Graph().Delete(input)
			if err != nil {
				return fmt.Errorf("graph delete failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal graph delete result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&graphOptions.ID, "id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("id")

	// グラフの削除は非常に危険なので `--delete-me` フラグが指定したときだけ削除する
	cmd.Flags().BoolVarP(&graphOptions.DeleteMe, "delete-me", "", false, "Delete your Graph")

	return cmd
}

func createGraphDeleteInput() *pixela.GraphDeleteInput {
	return &pixela.GraphDeleteInput{
		ID: getStringPtr(graphOptions.ID),
	}
}

// NewCmdGraphGetPixelDates creates a get pixel dates command.
func NewCmdGraphGetPixelDates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pixels",
		Short: "Get a Date list of Pixel registered",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createGraphGetPixelDatesInput()
			dates, err := pixelaClient.Graph().GetPixelDates(input)
			if err != nil {
				return fmt.Errorf("graph get pixel dates failed: %w", err)
			}

			if dates.IsSuccess == false {
				s, err := marshalResult(&dates.Result)
				if err != nil {
					return fmt.Errorf("marshal graph get pixel dates result failed: %w", err)
				}
				cmd.Printf("%s\n", s)
				return ErrNeglect
			}

			b, err := marshalPixels(dates.Pixels, graphOptions.WithBody)
			if err != nil {
				return fmt.Errorf("marshal graph get pixel dates failed: %w", err)
			}
			cmd.Printf("%s\n", string(b))

			return nil
		},
	}

	cmd.Flags().StringVar(&graphOptions.ID, "id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&graphOptions.From, "from", "", "The start position of the period")
	cmd.Flags().StringVar(&graphOptions.To, "to", "", "The end position of the period")
	cmd.Flags().BoolVar(&graphOptions.WithBody, "with-body", false, "Get all the information the Pixel has")

	return cmd
}

func createGraphGetPixelDatesInput() *pixela.GraphGetPixelDatesInput {
	return &pixela.GraphGetPixelDatesInput{
		ID:       getStringPtr(graphOptions.ID),
		From:     getStringPtr(graphOptions.From),
		To:       getStringPtr(graphOptions.To),
		WithBody: getBoolPtr(graphOptions.WithBody),
	}
}

func marshalPixels(datePixels interface{}, withBody bool) ([]byte, error) {
	if withBody {
		p, ok := datePixels.([]pixela.PixelWithBody)
		if !ok {
			return []byte{}, fmt.Errorf("type assertion failed: %T", datePixels)
		}
		b, err := json.Marshal(&pixelsWithBody{Pixels: p})
		if err != nil {
			return []byte{}, fmt.Errorf("marshal pixels with body failed: %w", err)
		}
		return b, nil
	}

	p, ok := datePixels.([]string)
	if !ok {
		return []byte{}, fmt.Errorf("type assertion failed: %T", datePixels)
	}
	b, err := json.Marshal(&pixels{Pixels: p})
	if err != nil {
		return []byte{}, fmt.Errorf("marshal pixels with body failed: %w", err)
	}
	return b, nil
}

type pixelsWithBody struct {
	Pixels []pixela.PixelWithBody `json:"pixels"`
}

type pixels struct {
	Pixels []string `json:"pixels"`
}

// NewCmdGraphStopwatch creates a graph stopwatch command.
func NewCmdGraphStopwatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stopwatch",
		Short: "Start and end the measurement of the time",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createGraphStopwatchInput()
			result, err := pixelaClient.Graph().Stopwatch(input)
			if err != nil {
				return fmt.Errorf("graph stopwatch failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal graph stopwatch result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&graphOptions.ID, "id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("id")

	return cmd
}

func createGraphStopwatchInput() *pixela.GraphStopwatchInput {
	return &pixela.GraphStopwatchInput{
		ID: getStringPtr(graphOptions.ID),
	}
}
