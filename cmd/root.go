/*
Copyright © 2024 muleyuck <takuty.008.awenite.1121@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Takuty-a11y/linippet/internal/file"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "linippet",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		dataPath, err := file.CheckDataPath()
		if err != nil {
			panic(err)
		}
		linippets, err := file.ReadJsonFile(dataPath)
		if err != nil {
			panic(err)
		}
		// TODO: Fuzzy Search
		linippet := linippets[0]

		re := regexp.MustCompile(`\${{(\w+)}}`)
		matchArgs := re.FindAllStringSubmatch(linippet.Snippet, -1)
		// without change when no args
		if len(matchArgs) <= 0 {
			fmt.Println(linippet.Snippet)
			return nil
		}
		// arr := make([]string, len(matchArgs))
		// for _, matchArg := range matchArgs {
		// 	arr = append(arr, matchArg[1])
		// }
		// TODO: input args by tui
		fmt.Println(linippet.Snippet)
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
