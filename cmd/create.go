/*
Copyright © 2025 muleyuck <takuty.008.awenite.1121@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/muleyuck/linippet/internal/linippet"
	"github.com/muleyuck/linippet/internal/tui"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create new a snippet",
	Long:  `Create new a snippet command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t := tui.NewCreateTui()
		t.SetAction()
		if err := t.StartApp(); err != nil {
			panic(err)
		}
		if !t.Submit {
			return nil
		}
		if len(t.Result) <= 0 {
			fmt.Println("Cannot create blank snippet.")
			return nil
		}
		if err := linippet.AddLinippet(t.Result); err != nil {
			return err
		}
		fmt.Println("Success to create snippet!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
