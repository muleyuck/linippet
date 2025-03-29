/*
Copyright © 2025 muleyuck <takuty.008.awenite.1121@gmail.com>
*/
package cmd

import (
	_ "embed"
	"fmt"

	"github.com/muleyuck/linippet/scripts"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize linppet",
	Long:  "set environment and key bind to initialize linippet",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must be specified Shell name. [example: linippet init bash]")
		}
		shellName := args[0]
		switch shellName {
		case "zsh":
			fmt.Printf("%s", scripts.InitializeZShellScript)
			return nil
		case "bash":
			fmt.Printf("%s", scripts.InitializeBashScript)
			return nil
		}
		return fmt.Errorf("%s is Unsupported Shell", shellName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
