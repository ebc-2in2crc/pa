package cmd

import (
	"encoding/json"
	"fmt"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/spf13/cobra"
)

var pixelOptions = &struct {
	GraphID      string
	Date         string
	Quantity     string
	OptionalData string
}{}

// NewCmdPixel creates a pixel command.
func NewCmdPixel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pixel",
		Short: "Pixel",
		Args:  cobra.NoArgs,
		RunE:  showHelp,
	}

	cmd.AddCommand(NewCmdPixelCreate())
	cmd.AddCommand(NewCmdPixelIncrement())
	cmd.AddCommand(NewCmdPixelDecrement())
	cmd.AddCommand(NewCmdPixelGet())
	cmd.AddCommand(NewCmdPixelUpdate())
	cmd.AddCommand(NewCmdPixelDelete())

	return cmd
}

// NewCmdPixelCreate creates a create pixel command.
func NewCmdPixelCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Pixel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createPixelCreateInput()
			result, err := pixelaClient.Pixel().Create(input)
			if err != nil {

				return fmt.Errorf("pixel create failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal pixel create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&pixelOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")
	cmd.Flags().StringVar(&pixelOptions.Date, "date", "", "The date on which the quantity is to be recorded")
	cmd.Flags().StringVar(&pixelOptions.Quantity, "quantity", "", "The quantity to be registered on the specified date")
	cmd.Flags().StringVar(&pixelOptions.OptionalData, "optional-data", "", "Additional information other than quantity")

	return cmd
}

func createPixelCreateInput() *pixela.PixelCreateInput {
	return &pixela.PixelCreateInput{
		GraphID:      getStringPtr(pixelOptions.GraphID),
		Date:         getStringPtr(pixelOptions.Date),
		Quantity:     getStringPtr(pixelOptions.Quantity),
		OptionalData: getStringPtr(pixelOptions.OptionalData),
	}
}

// NewCmdPixelIncrement creates a increment pixel command.
func NewCmdPixelIncrement() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "increment",
		Short: "Increment quantity 'Pixel' of the day",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createPixelIncrementInput()
			result, err := pixelaClient.Pixel().Increment(input)
			if err != nil {
				return fmt.Errorf("pixel increment failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal pixel increment result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&pixelOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")

	return cmd
}

func createPixelIncrementInput() *pixela.PixelIncrementInput {
	return &pixela.PixelIncrementInput{
		GraphID: getStringPtr(pixelOptions.GraphID),
	}
}

// NewCmdPixelDecrement creates a decrement pixel command.
func NewCmdPixelDecrement() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrement",
		Short: "Decrement quantity 'Pixel' of the day",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createPixelDecrementInput()
			result, err := pixelaClient.Pixel().Decrement(input)
			if err != nil {
				return fmt.Errorf("pixel decrement failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal pixel decrement result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&pixelOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")

	return cmd
}

func createPixelDecrementInput() *pixela.PixelDecrementInput {
	return &pixela.PixelDecrementInput{
		GraphID: getStringPtr(pixelOptions.GraphID),
	}
}

// NewCmdPixelGet creates a get pixel command.
func NewCmdPixelGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get registered quantity as 'Pixel'",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createPixelGetInput()
			q, err := pixelaClient.Pixel().Get(input)
			if err != nil {
				return fmt.Errorf("pixel get failed: %w", err)
			}

			if q.IsSuccess == false {
				s, err := marshalResult(&q.Result)
				if err != nil {
					return fmt.Errorf("marshal pixel get result failed: %w", err)
				}
				cmd.Printf("%s\n", s)
				return ErrNeglect
			}

			b, err := json.Marshal(&quantity{Quantity: q.Quantity, OptionalData: q.OptionalData})
			if err != nil {
				return fmt.Errorf("marshal pixel get failed: %w", err)
			}
			cmd.Printf("%s\n", string(b))

			return nil
		},
	}

	cmd.Flags().StringVar(&pixelOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")
	cmd.Flags().StringVar(&pixelOptions.Date, "date", "", "The date on which the quantity is to be recorded")
	cmd.MarkFlagRequired("date")

	return cmd
}

func createPixelGetInput() *pixela.PixelGetInput {
	return &pixela.PixelGetInput{
		GraphID: getStringPtr(pixelOptions.GraphID),
		Date:    getStringPtr(pixelOptions.Date),
	}
}

type quantity struct {
	Quantity     string `json:"quantity"`
	OptionalData string `json:"optionalData"`
}

// NewCmdPixelUpdate creates a update pixel command.
func NewCmdPixelUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Pixel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createPixelUpdateInput()
			result, err := pixelaClient.Pixel().Update(input)
			if err != nil {
				return fmt.Errorf("pixel update failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal update create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&pixelOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")
	cmd.Flags().StringVar(&pixelOptions.Date, "date", "", "The date on which the quantity is to be recorded")
	cmd.MarkFlagRequired("date")
	cmd.Flags().StringVar(&pixelOptions.Quantity, "quantity", "", "The quantity to be registered on the specified date")
	cmd.Flags().StringVar(&pixelOptions.OptionalData, "optional-data", "", "Additional information other than quantity")

	return cmd
}

func createPixelUpdateInput() *pixela.PixelUpdateInput {
	return &pixela.PixelUpdateInput{
		GraphID:      getStringPtr(pixelOptions.GraphID),
		Date:         getStringPtr(pixelOptions.Date),
		Quantity:     getStringPtr(pixelOptions.Quantity),
		OptionalData: getStringPtr(pixelOptions.OptionalData),
	}
}

// NewCmdPixelDelete creates a delete pixel command.
func NewCmdPixelDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Pixel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createPixelDeleteInput()
			result, err := pixelaClient.Pixel().Delete(input)
			if err != nil {
				return fmt.Errorf("pixel delete failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal delete create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&pixelOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")
	cmd.Flags().StringVar(&pixelOptions.Date, "date", "", "The date on which the quantity is to be recorded")
	cmd.MarkFlagRequired("date")

	return cmd
}

func createPixelDeleteInput() *pixela.PixelDeleteInput {
	return &pixela.PixelDeleteInput{
		GraphID: getStringPtr(pixelOptions.GraphID),
		Date:    getStringPtr(pixelOptions.Date),
	}
}
