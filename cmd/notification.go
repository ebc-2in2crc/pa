package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/spf13/cobra"
)

var notificationOptions = &struct {
	ID        string
	Name      string
	Target    string
	Condition string
	Threshold string
	RemindBy  string
	ChannelID string
	GraphID   string
}{}

// NewCmdNotification creates a notification command.
func NewCmdNotification() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notification",
		Short: "Notification",
		Args:  cobra.NoArgs,
		RunE:  showHelp,
	}

	cmd.AddCommand(NewCmdNotificationCreate())
	cmd.AddCommand(NewCmdNotificationGetAll())
	cmd.AddCommand(NewCmdNotificationUpdate())
	cmd.AddCommand(NewCmdNotificationDelete())

	return cmd
}

// NewCmdNotificationCreate creates a create notification command.
func NewCmdNotificationCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Notification",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createNotificationCreateInput()
			result, err := pixelaClient.Notification().Create(input)
			if err != nil {
				return fmt.Errorf("notification create failed: %w", err)
			}
			s, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal notification create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&notificationOptions.ID, "id", "", "ID for identifying the channel")
	cmd.Flags().StringVar(&notificationOptions.Name, "name", "", "The name of the channel")
	cmd.Flags().StringVar(&notificationOptions.Target, "target", "", "the target to be notified")
	cmd.Flags().StringVar(&notificationOptions.Condition, "condition", "", "the condition used to judge whether to notify or not")
	cmd.Flags().StringVar(&notificationOptions.Threshold, "threshold", "", "The threshold value for deciding whether to notify or not")
	cmd.Flags().StringVar(&notificationOptions.RemindBy, "remind-by", "", "Time of day")
	cmd.Flags().StringVar(&notificationOptions.ChannelID, "channel-id", "", "The ID of the channel to be notified")
	cmd.Flags().StringVar(&notificationOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")

	return cmd
}

func createNotificationCreateInput() *pixela.NotificationCreateInput {
	return &pixela.NotificationCreateInput{
		ID:        getStringPtr(notificationOptions.ID),
		Name:      getStringPtr(notificationOptions.Name),
		Target:    getStringPtr(notificationOptions.Target),
		Condition: getStringPtr(notificationOptions.Condition),
		Threshold: getStringPtr(notificationOptions.Threshold),
		RemindBy:  getStringPtr(notificationOptions.RemindBy),
		ChannelID: getStringPtr(notificationOptions.ChannelID),
		GraphID:   getStringPtr(notificationOptions.GraphID),
	}
}

// NewCmdNotificationGetAll creates a get all notification command.
func NewCmdNotificationGetAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Notification",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createNotificationGetAllInput()
			nds, err := pixelaClient.Notification().GetAll(input)
			if err != nil {
				return fmt.Errorf("notification get all failed: %w", err)
			}

			if nds.IsSuccess == false {
				s, err := marshalResult(&nds.Result)
				if err != nil {
					return fmt.Errorf("marshal notification get all result failed: %w", err)
				}
				cmd.Printf("%s\n", s)
				return ErrNeglect
			}

			b, err := json.Marshal(&notificationDefinitions{Notifications: nds.Notifications})
			if err != nil {
				return fmt.Errorf("marshal notification get all failed: %w", err)
			}

			s := string(b)
			s = strings.ReplaceAll(s, "\\u003c", "<")
			s = strings.ReplaceAll(s, "\\u003e", ">")
			cmd.Printf("%s\n", s)

			return nil
		},
	}

	cmd.Flags().StringVar(&notificationOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")

	return cmd
}

func createNotificationGetAllInput() *pixela.NotificationGetAllInput {
	return &pixela.NotificationGetAllInput{
		GraphID: getStringPtr(notificationOptions.GraphID),
	}
}

type notificationDefinitions struct {
	Notifications []pixela.NotificationDefinition `json:"notifications"`
}

// NewCmdNotificationUpdate creates a update notification command.
func NewCmdNotificationUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Notification",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createNotificationUpdateInput()
			result, err := pixelaClient.Notification().Update(input)
			if err != nil {
				return fmt.Errorf("notification update failed: %w", err)
			}
			s, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal notification update result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&notificationOptions.ID, "id", "", "ID for identifying the channel")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&notificationOptions.Name, "name", "", "The name of the channel")
	cmd.Flags().StringVar(&notificationOptions.Target, "target", "", "the target to be notified")
	cmd.Flags().StringVar(&notificationOptions.Condition, "condition", "", "the condition used to judge whether to notify or not")
	cmd.Flags().StringVar(&notificationOptions.Threshold, "threshold", "", "The threshold value for deciding whether to notify or not")
	cmd.Flags().StringVar(&notificationOptions.RemindBy, "remind-by", "", "Time of day")
	cmd.Flags().StringVar(&notificationOptions.ChannelID, "channel-id", "", "The ID of the channel to be notified")
	cmd.Flags().StringVar(&notificationOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")

	return cmd
}

func createNotificationUpdateInput() *pixela.NotificationUpdateInput {
	return &pixela.NotificationUpdateInput{
		ID:        getStringPtr(notificationOptions.ID),
		Name:      getStringPtr(notificationOptions.Name),
		Target:    getStringPtr(notificationOptions.Target),
		Condition: getStringPtr(notificationOptions.Condition),
		Threshold: getStringPtr(notificationOptions.Threshold),
		RemindBy:  getStringPtr(notificationOptions.RemindBy),
		ChannelID: getStringPtr(notificationOptions.ChannelID),
		GraphID:   getStringPtr(notificationOptions.GraphID),
	}
}

// NewCmdNotificationDelete creates a delete notification command.
func NewCmdNotificationDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Notification",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createNotificationDeleteInput()
			result, err := pixelaClient.Notification().Delete(input)
			if err != nil {
				return fmt.Errorf("notification delete failed: %w", err)
			}
			s, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal notification delete result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&notificationOptions.ID, "id", "", "ID for identifying the channel")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&notificationOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.MarkFlagRequired("graph-id")

	return cmd
}

func createNotificationDeleteInput() *pixela.NotificationDeleteInput {
	return &pixela.NotificationDeleteInput{
		ID:      getStringPtr(notificationOptions.ID),
		GraphID: getStringPtr(notificationOptions.GraphID),
	}
}
