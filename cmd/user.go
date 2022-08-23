package cmd

import (
	"fmt"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userOptions = &struct {
	AgreeTermsOfService bool
	NotMinor            bool
	ThanksCode          string
	NewToken            string
	DeleteMe            bool
}{}

// NewCmdUser creates a user command.
func NewCmdUser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "User",
		Args:  cobra.NoArgs,
		RunE:  showHelp,
	}

	cmd.AddCommand(NewCmdUserCreate())
	cmd.AddCommand(NewCmdUserUpdate())
	cmd.AddCommand(NewCmdUserDelete())

	return cmd
}

// NewCmdUserCreate creates a create user command.
func NewCmdUserCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Pixela user",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("agree_terms_of_service", cmd.Flags().Lookup("agree-terms-of-service")); err != nil {
				return fmt.Errorf("bind flag failed: %w", err)
			}
			if err := viper.BindPFlag("not_minor", cmd.Flags().Lookup("not-minor")); err != nil {
				return fmt.Errorf("bind flag failed: %w", err)
			}
			if err := viper.BindPFlag("thanks_code", cmd.Flags().Lookup("thanks-code")); err != nil {
				return fmt.Errorf("bind flag failed: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createUserCreateInput()
			result, err := pixelaClient.User().Create(input)
			if err != nil {
				return fmt.Errorf("user create failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal user create result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if !result.IsSuccess {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&userOptions.AgreeTermsOfService, "agree-terms-of-service", "a", false, "Agree to the terms of service")
	cmd.Flags().BoolVarP(&userOptions.NotMinor, "not-minor", "m", false, "You are not a minor or if you are a minor and you have the parental consent of using this service")
	cmd.Flags().StringVarP(&userOptions.ThanksCode, "thanks-code", "c", "", "Like a registration code obtained when you register for Patreon support")

	return cmd
}

func createUserCreateInput() *pixela.UserCreateInput {
	return &pixela.UserCreateInput{
		AgreeTermsOfService: getBoolFlag("agree_terms_of_service"),
		NotMinor:            getBoolFlag("not_minor"),
		ThanksCode:          getStringFlag("thanks_code"),
	}
}

// NewCmdUserUpdate creates a update user command.
func NewCmdUserUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates user token",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("new_token", cmd.Flags().Lookup("new-token")); err != nil {
				return fmt.Errorf("bind flag failed: %w", err)
			}
			if err := viper.BindPFlag("thanks_code", cmd.Flags().Lookup("thanks-code")); err != nil {
				return fmt.Errorf("bind flag failed: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			input := createUserUpdateInput()
			result, err := pixelaClient.User().Update(input)
			if err != nil {
				return fmt.Errorf("user update failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal user update result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if !result.IsSuccess {
				return ErrNeglect
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&userOptions.NewToken, "new-token", "n", "", "A new authentication token for update")
	cmd.Flags().StringVarP(&userOptions.ThanksCode, "thanks-code", "c", "", "Like a registration code obtained when you register for Patreon support")

	return cmd
}

func createUserUpdateInput() *pixela.UserUpdateInput {
	return &pixela.UserUpdateInput{
		NewToken:   getStringFlag("new_token"),
		ThanksCode: getStringFlag("thanks_code"),
	}
}

// NewCmdUserDelete creates a update user command.
func NewCmdUserDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a Pixela user",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !userOptions.DeleteMe {
				cmd.Println("Specify the '--delete-me' flag to confirm the deletion.")
				return nil
			}

			result, err := pixelaClient.User().Delete()
			if err != nil {
				return fmt.Errorf("user delete failed: %w", err)
			}
			s, err := marshalResult(result)
			if err != nil {
				return fmt.Errorf("marshal user delete result failed: %w", err)
			}
			cmd.Printf("%s\n", s)

			if !result.IsSuccess {
				return ErrNeglect
			}
			return nil
		},
	}

	// ユーザーの削除は非常に危険なので `--delete-me` フラグが指定したときだけ削除する
	// 環境変数は読み込まないように viper にはバインドしないでおく
	cmd.Flags().BoolVarP(&userOptions.DeleteMe, "delete-me", "", false, "Delete your Pixela account")

	return cmd
}
