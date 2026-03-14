package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all the mods in the modpack",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// Load pack
		pack, err := core.LoadPack()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Load index
		index, err := pack.LoadIndex()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Load mods
		mods, err := index.LoadAllMods()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Filter mods by side
		if viper.IsSet("list.side") {
			side := viper.GetString("list.side")
			if side != core.UniversalSide && side != core.ServerSide && side != core.ClientSide {
				fmt.Printf("Invalid side %q, must be one of client, server, or both (default)\n", side)
				os.Exit(1)
			}

			i := 0
			for _, mod := range mods {
				if mod.Side == side || mod.Side == core.EmptySide || mod.Side == core.UniversalSide || side == core.UniversalSide {
					mods[i] = mod
					i++
				}
			}
			mods = mods[:i]
		}

		//Filter mods by tag
		if viper.IsSet("list.tag") {
			targetTag := viper.GetString("list.tag")

			i := 0
			for _, mod := range mods {
				for _, tag := range mod.Tags {
					if targetTag == tag {
						mods[i] = mod
						i++
						break
					}
				}
			}
			if i == 0 {
				fmt.Printf("Could not find a mod with tag %q", targetTag)
				os.Exit(1)
			} else {
				mods = mods[:i]
			}
		}

		sort.Slice(mods, func(i, j int) bool {
			return strings.ToLower(mods[i].Name) < strings.ToLower(mods[j].Name)
		})

		// Print mods
		if viper.GetBool("list.version") {
			for _, mod := range mods {
				fmt.Printf("%s (%s)\n", mod.Name, mod.FileName)
			}
		} else if viper.GetBool("list.meta") {
			//print toml metaFile names
			for _, mod := range mods {
				metaFile := strings.TrimSuffix(filepath.Base(mod.GetFilePath()), core.MetaExtension)
				fmt.Printf("%s (%s)\n", mod.Name, metaFile)
			}
		} else {
			for _, mod := range mods {
				fmt.Println(mod.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("version", "v", false, "Print name and version")
	_ = viper.BindPFlag("list.version", listCmd.Flags().Lookup("version"))
	listCmd.Flags().BoolP("meta", "m", false, "Print name and metadata file slug")
	_ = viper.BindPFlag("list.meta", listCmd.Flags().Lookup("meta"))
	listCmd.Flags().StringP("side", "s", "", "Filter mods by side (e.g., client or server)")
	_ = viper.BindPFlag("list.side", listCmd.Flags().Lookup("side"))
	listCmd.Flags().StringP("tag", "t", "", "Filter mods by tag")
	_ = viper.BindPFlag("list.tag", listCmd.Flags().Lookup("tag"))

}
