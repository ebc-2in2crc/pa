package cmd

import (
	"errors"
	"fmt"
	"os"

	pixela "github.com/ebc-2in2crc/pixela4go"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var globalOptions = &struct {
	username   string
	token      string
	retryCount int
}{}

var rootCmd *cobra.Command
var cfgFile string
var version = "dev"

// NewCmdRoot creates a root command.
func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pa",
		Short:         "The Pixela Command Line Interface is a unified tool to manage your Pixela services",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			pixela.RetryCount = getRetry()
		},
	}

	cobra.OnInitialize(initConfig)
	cmd.Version = version

	viper.AutomaticEnv()
	viper.SetEnvPrefix("pa")
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pa)")
	cmd.PersistentFlags().StringVarP(&globalOptions.username, "username", "u", "", "Pixela user name")
	_ = viper.BindPFlag("username", cmd.PersistentFlags().Lookup("username"))
	cmd.PersistentFlags().StringVarP(&globalOptions.token, "token", "t", "", "Pixela user token")
	_ = viper.BindPFlag("token", cmd.PersistentFlags().Lookup("token"))
	cmd.PersistentFlags().IntVarP(&globalOptions.retryCount, "retry", "r", 0, "Specify the number of retries when the API call is rejected")
	_ = viper.BindPFlag("retry", cmd.PersistentFlags().Lookup("retry"))

	addSubCommand(cmd)

	return cmd
}

func addSubCommand(cmd *cobra.Command) {
	cmd.AddCommand(NewCmdUser())
	cmd.AddCommand(NewCmdUserProfile())
	cmd.AddCommand(NewCmdGraph())
	cmd.AddCommand(NewCmdPixel())
	cmd.AddCommand(NewCmdWebhook())
	cmd.AddCommand(NewCmdCompletion())
}

// Execute executes root command.
func Execute() {
	rootCmd = NewCmdRoot()
	rootCmd.SetOut(os.Stdout)

	err := rootCmd.Execute()
	if err != nil {
		if !errors.Is(err, ErrNeglect) {
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

		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(wd)
		viper.AddConfigPath(home)
		viper.SetConfigName(".pa")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func getUsername() string {
	return viper.GetString("username")
}

func getToken() string {
	return viper.GetString("token")
}

func getRetry() int {
	return viper.GetInt("retry")
}
