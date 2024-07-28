/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Takuty-a11y/linippet/utils"
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
	Run: func(cmd *cobra.Command, args []string) {
		addSnippet()
	},
}

func addSnippet() {
	snp := interactive()
	path, err := utils.GetSnippetFilePath()
	if err != nil {
		utils.Fatal(err)
	}
	err = utils.WriteFile(path, snp+"\n")
	if err != nil {
		utils.Fatal(err)
	}
}

func interactive() string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "What is CommandLine?:"+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func init() {
	rootCmd.AddCommand(createCmd)
}
