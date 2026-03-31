package settings

import (
	"fmt"
	"os"
	"strings"

	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var defaultPlatformCmd = &cobra.Command{
	Use:       "default-platform [platform1] [platform2]",
	Short:     "Set default platform(s) used with \"packwiz install\" command.\nA mod will be searched for on the first provided platform and then the next if an exact match is not found.",
	Aliases:   []string{"dp", "default-platforms", "platforms", "platform"},
	ValidArgs: []string{"modrinth", "curseforge", "url"},
	Args:      cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		//TODO: some of the below is likely necessary for all settings commands, certainly everything
		//affecting "options". Perhaps make a boilerplate optionsInit with the below code to be used across commands.
		modpack, err := core.LoadPack()
		if err != nil {
			// Check if it's a no such file or directory error
			if os.IsNotExist(err) {
				fmt.Println("No pack.toml file found, run 'packwiz init' to create one!")
				os.Exit(1)
			}
			fmt.Printf("Error loading pack: %s\n", err)
			os.Exit(1)
		}
		// Check if they have no options whatsoever
		if modpack.Options == nil {
			// Initialize the options
			modpack.Options = make(map[string]interface{})
		}

		modpack.Options["default-platforms"] = args

		err = modpack.Write()
		if err != nil {
			fmt.Printf("Error writing pack: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Default platforms set are now: \"%s\"", strings.Join(args, ", "))
	},
}

func init() {
	settingsCmd.AddCommand(defaultPlatformCmd)
}
