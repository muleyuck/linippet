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

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove a snippet.",
	Long:  "Remove a snippet which be chosen from your snippets list",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := tui.NewRemoveTui()
		t.LazyLoadLinippet()
		t.SetAction()
		if err := t.StartApp(); err != nil {
			panic(err)
		}
		if !t.Submit {
			return nil
		}
		if err := linippet.RemoveLinippet(t.SelectId); err != nil {
			return err
		}
		fmt.Println("Success to remove snippet!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
