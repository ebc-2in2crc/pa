package cmd

import (
	"fmt"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/spf13/cobra"
)

var userProfileOptions = &struct {
	DisplayName       string
	GravatarIconEmail string
	Title             string
	Timezone          string
	AboutURL          string
	ContributeURLs    []string
	PinnedGraphID     string
}{}

// NewCmdUserProfile creates a user profile command.
func NewCmdUserProfile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Profile",
		Args:  cobra.NoArgs,
		RunE:  showHelp,
	}

	cmd.AddCommand(NewCmdUserProfileUpdate())
	cmd.AddCommand(NewCmdUserProfileURL())

	return cmd
}

// NewCmdUserProfileUpdate creates a update user profile command.
func NewCmdUserProfileUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates User Profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createUserProfileUpdateInput()
			result, err := pixelaClient.UserProfile().Update(input)
			if err != nil {
				return fmt.Errorf("user profile update failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal user profile update result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if result.IsSuccess == false {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&userProfileOptions.DisplayName, "display-name", "", "The user's name for the display")
	cmd.Flags().StringVar(&userProfileOptions.GravatarIconEmail, "gravatar-icon-email", "", "The email address registered as an icon in Gravatar")
	cmd.Flags().StringVar(&userProfileOptions.Title, "title", "", "The title of the user")
	cmd.Flags().StringVar(&userProfileOptions.Timezone, "timezone", "", "Specify the user's time zone")
	cmd.Flags().StringVar(&userProfileOptions.AboutURL, "about-url", "", "Users can only show one external link")
	cmd.Flags().StringSliceVar(&userProfileOptions.ContributeURLs, "contribute-urls", []string{}, "The contribute URLs")
	cmd.Flags().StringVar(&userProfileOptions.PinnedGraphID, "pinned-graph-id", "", " Pin one of their own graphs")

	return cmd
}

func createUserProfileUpdateInput() *pixela.UserProfileUpdateInput {
	return &pixela.UserProfileUpdateInput{
		DisplayName:       getStringPtr(userProfileOptions.DisplayName),
		GravatarIconEmail: getStringPtr(userProfileOptions.GravatarIconEmail),
		Title:             getStringPtr(userProfileOptions.Title),
		Timezone:          getStringPtr(userProfileOptions.Timezone),
		AboutURL:          getStringPtr(userProfileOptions.AboutURL),
		ContributeURLs:    userProfileOptions.ContributeURLs,
		PinnedGraphID:     getStringPtr(userProfileOptions.PinnedGraphID),
	}
}

// NewCmdUserProfileURL creates a user profile URL command.
func NewCmdUserProfileURL() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get User Profile page URL",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := pixelaClient.UserProfile().URL()
			cmd.Printf("%s\n", url)

			return nil
		},
	}

	return cmd
}
