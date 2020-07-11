package cmd

import (
	"encoding/json"
	"fmt"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/spf13/cobra"
)

var channelOptions = &struct {
	ID               string
	Name             string
	Type             string
	SlackUserName    string
	SlackChannelName string
	SlackURL         string
}{}

// NewCmdChannel creates a new channel command.
func NewCmdChannel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "Channel",
		Args:  cobra.NoArgs,
		RunE:  showHelp,
	}

	cmd.AddCommand(NewCmdChannelCreate())
	cmd.AddCommand(NewCmdChannelGetAll())
	cmd.AddCommand(NewCmdChannelUpdate())
	cmd.AddCommand(NewCmdChannelDelete())

	return cmd
}

// NewCmdChannelCreate creates a create channel command.
func NewCmdChannelCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Channel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createChannelCreateInput()
			result, err := pixelaClient.Channel().Create(input)
			if err != nil {
				return fmt.Errorf("channel create failed: %w", err)
			}
			s, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal channel create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&channelOptions.ID, "id", "", "ID for identifying the channel")
	cmd.Flags().StringVar(&channelOptions.Name, "name", "", "The name of the channel")
	cmd.Flags().StringVar(&channelOptions.Type, "type", "", "The type for notification")
	cmd.Flags().StringVar(&channelOptions.SlackUserName, "slack-username", "", "See: https://docs.pixe.la/entry/post-channel 'userName'")
	cmd.Flags().StringVar(&channelOptions.SlackChannelName, "slack-channel-name", "", "See: https://docs.pixe.la/entry/post-channel 'channelName")
	cmd.Flags().StringVar(&channelOptions.SlackURL, "slack-url", "", "See: https://docs.pixe.la/entry/post-channel 'url")

	return cmd
}

func createChannelCreateInput() *pixela.ChannelCreateInput {
	return &pixela.ChannelCreateInput{
		ID:   getStringPtr(channelOptions.ID),
		Name: getStringPtr(channelOptions.Name),
		Type: getStringPtr(channelOptions.Type),
		SlackDetail: &pixela.SlackDetail{
			UserName:    getStringPtr(channelOptions.SlackUserName),
			ChannelName: getStringPtr(channelOptions.SlackChannelName),
			URL:         getStringPtr(channelOptions.SlackURL),
		},
	}
}

// NewCmdChannelGetAll creates a get all channel command.
func NewCmdChannelGetAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Channel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cds, err := pixelaClient.Channel().GetAll()
			if err != nil {
				return fmt.Errorf("channel get all failed: %w", err)
			}

			if cds.IsSuccess == false {
				s, err := marshalResult(&cds.Result)
				if err != nil {
					return fmt.Errorf("marshal channel get all result failed: %w", err)
				}
				cmd.Printf("%s\n", s)
				return ErrNeglect
			}

			b, err := json.Marshal(&channelDefinitions{Channels: cds.Channels})
			if err != nil {
				return fmt.Errorf("marshal channel get all failed: %w", err)
			}
			cmd.Printf("%s\n", string(b))

			return nil
		},
	}

	return cmd
}

type channelDefinitions struct {
	Channels []pixela.ChannelDefinition `json:"channels"`
}

// NewCmdChannelUpdate creates a update channel command.
func NewCmdChannelUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Channel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createChannelUpdateInput()
			result, err := pixelaClient.Channel().Update(input)
			if err != nil {
				return fmt.Errorf("channel update failed: %w", err)
			}
			s, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal channel update result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&channelOptions.ID, "id", "", "ID for identifying the channel")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&channelOptions.Name, "name", "", "The name of the channel")
	cmd.Flags().StringVar(&channelOptions.Type, "type", "", "The type for notification")
	cmd.Flags().StringVar(&channelOptions.SlackUserName, "slack-username", "", "See: https://docs.pixe.la/entry/post-channel 'userName'")
	cmd.Flags().StringVar(&channelOptions.SlackChannelName, "slack-channel-name", "", "See: https://docs.pixe.la/entry/post-channel 'channelName")
	cmd.Flags().StringVar(&channelOptions.SlackURL, "slack-url", "", "See: https://docs.pixe.la/entry/post-channel 'url")

	return cmd
}

func createChannelUpdateInput() *pixela.ChannelUpdateInput {
	return &pixela.ChannelUpdateInput{
		ID:   getStringPtr(channelOptions.ID),
		Name: getStringPtr(channelOptions.Name),
		Type: getStringPtr(channelOptions.Type),
		SlackDetail: &pixela.SlackDetail{
			UserName:    getStringPtr(channelOptions.SlackUserName),
			ChannelName: getStringPtr(channelOptions.SlackChannelName),
			URL:         getStringPtr(channelOptions.SlackURL),
		},
	}
}

// NewCmdChannelDelete creates a delete channel command.
func NewCmdChannelDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Channel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createChannelDeleteInput()
			result, err := pixelaClient.Channel().Delete(input)
			if err != nil {
				return fmt.Errorf("channel delete failed: %w", err)
			}
			s, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal channel delete result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&channelOptions.ID, "id", "", "ID for identifying the channel")
	cmd.MarkFlagRequired("id")

	return cmd
}

func createChannelDeleteInput() *pixela.ChannelDeleteInput {
	return &pixela.ChannelDeleteInput{
		ID: getStringPtr(channelOptions.ID),
	}
}
