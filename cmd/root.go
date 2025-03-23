/*
Copyright © 2024 muleyuck <takuty.008.awenite.1121@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/muleyuck/linippet/internal/tui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "linippet",
	Short: "Choose your snippet and output stdout",
	Long:  `linippet is submit a snippet you have registered.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		t := tui.NewRootTui()
		t.LazyLoadLinippet()
		t.SetAction()
		if err := t.StartApp(); err != nil {
			panic(err)
		}
		fmt.Println(t.GetTrimmedResult())
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-cli-sample.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
