package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/packwiz/packwiz/cmdshared"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an external file from the modpack; equivalent to manually removing the file and running packwiz refresh",
	Aliases: []string{"delete", "uninstall", "rm"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeMod(args[0])
	},
}

type removeContext struct {
	pack     core.Pack
	index    core.Index
	mods     []*core.Mod         // all mods, loaded once
	toRemove map[string]core.Mod // keyed by slug for deduplication
}

func removeMod(modSlug string) {
	errCount := 0

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
	modPath, ok := index.FindMod(modSlug)
	if !ok {
		fmt.Println("Can't find this file; please ensure you have run packwiz refresh and use the name of the .pw.toml file (defaults to the project slug)")
		os.Exit(1)
	}
	mods, err := index.LoadAllMods()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mod, err := core.LoadMod(modPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := &removeContext{
		pack:     pack,
		index:    index,
		mods:     mods,
		toRemove: map[string]core.Mod{mod.Slug: mod},
	}
	//TODO: refine non-interactive use, perhaps non-recursive by default unless --recursive is used.
	if !viper.GetBool("non-interactive") {
		removeModRecursive(mod, ctx)
	}

	for _, m := range ctx.toRemove {
		mPath := m.GetFilePath()
		fmt.Printf("Deleting metadata file \"%s\"...\n", mPath)
		err = os.Remove(mPath)
		if err != nil {
			fmt.Println(err)
			errCount++
		}
		fmt.Println("Removing from index...")
		err = index.RemoveFile(mPath)
		if err != nil {
			fmt.Println(err)
			errCount++
		}
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
	if errCount > 0 {
		fmt.Println("Mods were removed with some errors! Please see above.")
		os.Exit(1)
	} else {
		fmt.Printf("%d mods removed successfully!\n", len(ctx.toRemove))
	}
}

func removeModRecursive(mod core.Mod, ctx *removeContext) {
	if _, exists := ctx.toRemove[mod.Slug]; exists {
		return
	}
	//add target mod to list of mods to remove
	ctx.toRemove[mod.Slug] = mod

	// check for orphaned dependencies
	if mod.Relations != nil && len(mod.Relations.Dependencies) > 0 {
		var orphanedMods []core.Mod

		//assume all auto-installed dependencies are orphaned to start, load mod for each dependency slug
		for _, mSlug := range mod.Relations.Dependencies {
			mPath, ok := ctx.index.FindMod(mSlug)
			if !ok {
				fmt.Printf("Couldn't find \"%s\" dependency with slug \"%s\" in pack. Ignoring...\n", mod.Name, mSlug)
				continue
			}
			depMod, err := core.LoadMod(mPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if depMod.AutoInstalled {
				orphanedMods = append(orphanedMods, depMod)
			}
		}

		for _, otherMod := range ctx.mods {
			if otherMod.Relations != nil && len(otherMod.Relations.Dependencies) > 0 && otherMod.Slug != mod.Slug {
				for _, otherDep := range otherMod.Relations.Dependencies {
					//If another mod depends on one of the removed mod's dependencies,
					//it isn't orphaned and is removed from orphaned deps.
					orphanedMods = slices.DeleteFunc(orphanedMods, func(dep core.Mod) bool {
						return dep.Slug == otherDep
					})
				}
			}
		}

		if len(orphanedMods) > 0 {
			fmt.Printf("Removing \"%s\" will orphan the following dependencies:\n", mod.Name)
			for _, m := range orphanedMods {
				fmt.Println(m.Name)
			}
			if cmdshared.PromptYesNo("Would you like to remove these as well? [Y/n]") {
				for _, m := range orphanedMods {
					ctx.toRemove[m.Slug] = m
				}
			}
		}
	}

	//Check for affected dependents
	var dependentMods []core.Mod
	for _, otherMod := range ctx.mods {
		if otherMod.Relations != nil && len(otherMod.Relations.Dependencies) > 0 && otherMod.Slug != mod.Slug {
			for _, otherDep := range otherMod.Relations.Dependencies {
				//if another mod depends on the removed mod, it may be affected by its removal.
				if otherDep == mod.Slug {
					dependentMods = append(dependentMods, *otherMod)
					break
				}
			}
		}
	}
	//TODO: This does not account for the dependents and dependencies of
	if len(dependentMods) > 0 {
		fmt.Printf("The following mods are dependent on \"%s\" and may not work upon its removal:\n", mod.Name)
		for _, m := range dependentMods {
			fmt.Println(m.Name)
		}
		if cmdshared.PromptYesNo("Would you like to remove these as well? [Y/n]") {
			for _, m := range dependentMods {
				ctx.toRemove[m.Slug] = m
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
