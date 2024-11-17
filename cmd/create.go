/*
Copyright © 2024 muleyuck <takuty.008.awenite.1121@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/muleyuck/linippet/internal/file"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultCommand := ""
		// TODO: recieve input value from tui

		dataPath, err := file.CheckDataPath()
		if err != nil {
			return err
		}

		if err := file.WriteJsonFile(dataPath, defaultCommand); err != nil {
			return err
		}
		fmt.Println("Create snippet success!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
