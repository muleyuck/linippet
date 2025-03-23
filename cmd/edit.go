/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/muleyuck/linippet/internal/linippet"
	"github.com/muleyuck/linippet/internal/tui"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit a snippet.",
	Long:  "Edit snippet which be chosen from your snippets list",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := tui.NewEditTui()
		t.LazyLoadLinippet()
		t.SetAction()
		if err := t.StartApp(); err != nil {
			panic(err)
		}
		if !t.Submit {
			return nil
		}
		if len(t.Result) <= 0 {
			fmt.Println("Cannot save blank snippet.")
			return nil
		}
		if err := linippet.UpdateLinippet(t.SelectId, t.Result); err != nil {
			return err
		}
		fmt.Println("Success to edit snippet!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
