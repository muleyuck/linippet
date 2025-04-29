/*
Copyright © 2025 muleyuck <takuty.008.awenite.1121@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/muleyuck/linippet/internal/linippet"
	"github.com/muleyuck/linippet/internal/tui"
	"github.com/muleyuck/linippet/scripts"
	"github.com/spf13/cobra"
)

var (
	versionFlag bool
	listFlag    bool
)

var rootCmd = &cobra.Command{
	Use:   "linippet",
	Short: "Choose your snippet and output stdout",
	Long:  `linippet is submit a snippet you have registered.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		versionFlag, _ := cmd.Flags().GetBool("version")
		listFlag, _ := cmd.Flags().GetBool("list")
		if versionFlag {
			fmt.Printf("linippet %s", scripts.AppVersion)
		} else if listFlag {
			linippets, err := linippet.ReadLinippets()
			if err != nil {
				return err
			}
			if len(linippets) <= 0 {
				fmt.Println("linippet: There are no snippets")
				return nil
			}
			for i, linippet := range linippets {
				fmt.Printf("%d : %s\n", i+1, linippet.Snippet)
			}
		} else {
			t := tui.NewRootTui()
			t.LazyLoadLinippet()
			t.SetAction()
			if err := t.StartApp(); err != nil {
				panic(err)
			}
			fmt.Println(t.Result)
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "version")
	rootCmd.Flags().BoolVar(&listFlag, "list", false, "show snippet list")
}
