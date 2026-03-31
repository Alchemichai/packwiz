package cmd

import (
	"fmt"
	"os"

	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: implement multi-default installation
var installCmd = &cobra.Command{
	Use:     "install [mod slug/search term]",
	Short:   "attempt to install using default mod platforms in order set by \"packwiz settings default-platform\"",
	Aliases: []string{"i", "get", "add", "+"},
	Run: func(cmd *cobra.Command, args []string) {
		_, err := core.LoadPack()
		if err != nil {
			// Check if it's a no such file or directory error
			if os.IsNotExist(err) {
				fmt.Println("No pack.toml file found, run 'packwiz init' to create one!")
				os.Exit(1)
			}
			fmt.Printf("Error loading pack: %s\n", err)
			os.Exit(1)
		}

		if len(viper.GetStringSlice("default-platforms")) == 0 {
			fmt.Println("No default platform set. Please set one with \"packwiz settings default-platform\"")
			os.Exit(1)
		}
		rootCmd.SetArgs(append([]string{viper.GetStringSlice("default-platforms")[0], "install"}, args...))
		rootCmd.Execute()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
