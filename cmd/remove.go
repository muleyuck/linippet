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

// removeCmd represents the remove command
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
