/*
Copyright © 2024 muleyuck <takuty.008.awenite.1121@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/muleyuck/linippet/internal/linippet"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "show snippets list",
	Long:  `Your snippets will be output stdout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		linippets, err := linippet.ReadLinippets()
		if err != nil {
			return err
		}
		if len(linippets) <= 0 {
			fmt.Println("There are no snippets")
			return nil
		}
		for i, linippet := range linippets {
			fmt.Printf("%d : %s\n", i, linippet.Snippet)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
