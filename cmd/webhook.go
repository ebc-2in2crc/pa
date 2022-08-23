package cmd

import (
	"encoding/json"
	"fmt"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/spf13/cobra"
)

var webhookOptions = &struct {
	GraphID     string
	Type        string
	WebhookHash string
}{}

// NewCmdWebhook creates a webhook command.
func NewCmdWebhook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Webhook",
		Args:  cobra.NoArgs,
		RunE:  showHelp,
	}

	cmd.AddCommand(NewCmdWebhookCreate())
	cmd.AddCommand(NewCmdWebhookGetAll())
	cmd.AddCommand(NewCmdWebhookInvoke())
	cmd.AddCommand(NewCmdWebhookDelete())

	return cmd
}

// NewCmdWebhookCreate creates a create webhook command.
func NewCmdWebhookCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Webhook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createWebhookCreateInput()
			result, err := pixelaClient.Webhook().Create(input)
			if err != nil {
				return fmt.Errorf("webhook create failed: %w", err)
			}
			s, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal webhook create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if !result.IsSuccess {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&webhookOptions.GraphID, "graph-id", "", "ID for identifying the pixelation graph")
	cmd.Flags().StringVar(&webhookOptions.Type, "type", "", "The behavior when this Webhook is invoked")

	return cmd
}

func createWebhookCreateInput() *pixela.WebhookCreateInput {
	return &pixela.WebhookCreateInput{
		GraphID: getStringPtr(webhookOptions.GraphID),
		Type:    getStringPtr(webhookOptions.Type),
	}
}

// NewCmdWebhookGetAll creates a get all webhook command.
func NewCmdWebhookGetAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Webhook definitions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			whs, err := pixelaClient.Webhook().GetAll()
			if err != nil {
				return fmt.Errorf("webhook get all failed: %w", err)
			}

			if !whs.IsSuccess {
				s, err := marshalResult(&whs.Result)
				if err != nil {
					return fmt.Errorf("marshal webhook get all result failed: %w", err)
				}
				cmd.Printf("%s\n", s)
				return ErrNeglect
			}

			b, err := json.Marshal(&webhookDefinitions{Webhooks: whs.Webhooks})
			if err != nil {
				return fmt.Errorf("marshal webhook get all failed: %w", err)
			}
			cmd.Printf("%s\n", string(b))

			return nil
		},
	}

	return cmd
}

type webhookDefinitions struct {
	Webhooks []pixela.WebhookDefinition `json:"webhooks"`
}

// NewCmdWebhookInvoke creates a invoke webhook command.
func NewCmdWebhookInvoke() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoke",
		Short: "Invoke Webhook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createWebhookInvokeInput()
			result, err := pixelaClient.Webhook().Invoke(input)
			if err != nil {
				return fmt.Errorf("webhook invoke failed: %w", err)
			}
			s, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal invoke create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if !result.IsSuccess {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&webhookOptions.WebhookHash, "hash", "", "Webhook hash")
	_ = cmd.MarkFlagRequired("hash")

	return cmd
}

func createWebhookInvokeInput() *pixela.WebhookInvokeInput {
	return &pixela.WebhookInvokeInput{
		WebhookHash: getStringPtr(webhookOptions.WebhookHash),
	}
}

// NewCmdWebhookDelete creates a delete webhook command.
func NewCmdWebhookDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Webhook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createWebhookDeleteInput()
			result, err := pixelaClient.Webhook().Delete(input)
			if err != nil {
				return fmt.Errorf("webhook delete failed: %w", err)
			}
			s, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("marshal delete create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if !result.IsSuccess {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&webhookOptions.WebhookHash, "hash", "", "Webhook hash")
	_ = cmd.MarkFlagRequired("hash")

	return cmd
}

func createWebhookDeleteInput() *pixela.WebhookDeleteInput {
	return &pixela.WebhookDeleteInput{
		WebhookHash: getStringPtr(webhookOptions.WebhookHash),
	}
}
