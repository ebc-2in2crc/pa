package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var globalOptions = &struct {
	username string
	token    string
}{}

var rootCmd *cobra.Command
var cfgFile string
var version = "v0.0.1"

// NewCmdRoot creates a root command.
func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pa",
		Short:         "The Pixela Command Line Interface is a unified tool to manage your Pixela services",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cobra.OnInitialize(initConfig)
	cmd.Version = version

	viper.AutomaticEnv()
	viper.SetEnvPrefix("pa")
	cmd.PersistentFlags().StringVarP(&globalOptions.username, "username", "u", "", "Pixela user name")
	viper.BindPFlag("username", cmd.PersistentFlags().Lookup("username"))
	cmd.PersistentFlags().StringVarP(&globalOptions.token, "token", "t", "", "Pixela user token")
	viper.BindPFlag("token", cmd.PersistentFlags().Lookup("token"))

	addSubCommand(cmd)

	return cmd
}

func addSubCommand(cmd *cobra.Command) {
	cmd.AddCommand(NewCmdUser())
	cmd.AddCommand(NewCmdChannel())
	cmd.AddCommand(NewCmdGraph())
	cmd.AddCommand(NewCmdPixel())
	cmd.AddCommand(NewCmdNotification())
	cmd.AddCommand(NewCmdWebhook())
	cmd.AddCommand(NewCmdCompletion())
}

// Execute executes root command.
func Execute() {
	rootCmd = NewCmdRoot()
	rootCmd.SetOut(os.Stdout)

	err := rootCmd.Execute()
	if err != nil {
		if errors.Is(err, ErrNeglect) == false {
			rootCmd.PrintErr(err)
		}
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".pa")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getUsername() string {
	return viper.GetString("username")
}

func getToken() string {
	return viper.GetString("token")
}
