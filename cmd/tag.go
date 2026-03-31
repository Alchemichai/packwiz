package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:     "tag [modname] [tags...]",
	Short:   "Add one or more tags to a mod; equivalent to 'packwiz tag add'",
	Aliases: []string{"tags"},
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		modifyTags(args[0], args[1:], true)
	},
}

var tagAddCmd = &cobra.Command{
	Use:   "add [modname] [tags...]",
	Short: "Add one or more tags to a mod",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		modifyTags(args[0], args[1:], true)
	},
}

var tagRemoveCmd = &cobra.Command{
	Use:     "remove [modname] [tags...]",
	Short:   "Remove one or more tags to a mod",
	Aliases: []string{"rm"},
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		modifyTags(args[0], args[1:], false)
	},
}

var tagListCmd = &cobra.Command{
	Use:     "list [modname]",
	Short:   "List tags on a mod, or all tags in the pack if no mod is specified",
	Aliases: []string{"ls"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pack, err := core.LoadPack()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		index, err := pack.LoadIndex()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(args) == 0 {
			//No mod specified, list all tags across mods
			mods, err := index.LoadAllMods()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			var tagList []string

			for _, mod := range mods {
				tagList = append(tagList, mod.Tags...)
			}

			slices.Sort(tagList)
			tagList = slices.Compact(tagList)
			fmt.Println(strings.Join(tagList, "\n"))
		} else {
			// Mod specified, list tags on that mod
			modPath, ok := index.FindMod(args[0])
			if !ok {
				fmt.Printf("Can't find mod %q\n", args[0])
				os.Exit(1)
			}
			mod, err := core.LoadMod(modPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if len(mod.Tags) == 0 {
				fmt.Printf("%s has no tags\n", mod.Name)
				return
			}
			fmt.Printf(strings.Join(mod.Tags, "\n"))
		}
	},
}

// TODO: Refactor into get/add/set/remove.
// Should this be using this using modName, slugs, or IDs?
func modifyTags(modName string, tags []string, adding bool) {
	pack, err := core.LoadPack()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	index, err := pack.LoadIndex()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	modPath, ok := index.FindMod(modName)
	if !ok {
		fmt.Println("Can't find this file; please ensure you have run packwiz refresh and use the name of the .pw.toml file (defaults to the project slug)")
		os.Exit(1)
	}
	mod, err := core.LoadMod(modPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if adding {
		for _, tag := range tags {
			if !slices.Contains(mod.Tags, tag) {
				mod.Tags = append(mod.Tags, tag)
				fmt.Printf("Added tag %q to %q", tag, mod.Name)
			} else {
				fmt.Printf("Mod %q already contains tag %q, skipping...", mod.Name, tag)
			}
		}
	} else {
		for _, tag := range tags {
			if slices.Contains(mod.Tags, tag) {
				mod.Tags = slices.DeleteFunc(mod.Tags, func(t string) bool {
					return t == tag
				})
				fmt.Printf("Removed tag %q from %q", tag, mod.Name)
			} else {
				fmt.Printf("Mod %q does not contain tag %q, skipping...", mod.Name, tag)
			}
		}
	}

	format, hash, err := mod.Write()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = index.RefreshFileWithHash(modPath, format, hash, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = index.Write()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pack.UpdateIndexHash()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pack.Write()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func removeTags(modName string, tags []string) {
	for _, tag := range tags {
		if slices.Contains(tags, tag) {
			tags = slices.DeleteFunc(tags, func(t string) bool {
				return t == tag
			})
			fmt.Printf("Removed tag %q from %q", tag, modName)
		} else {
			fmt.Printf("Mod %q does not contain tag %q, skipping...", modName, tag)
		}
	}
}

func addTags(modName string, tags []string) {
	for _, tag := range tags {
		if !slices.Contains(tags, tag) {
			tags = append(tags, tag)
			fmt.Printf("Added tag %q to %q", tag, modName)
		} else {
			fmt.Printf("Mod %q already contains tag %q, skipping...", modName, tag)
		}
	}
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagRemoveCmd)
	tagCmd.AddCommand(tagListCmd)
}
