package cmd

import (
	"github.com/spf13/cobra"
)

var relationsCmd = &cobra.Command{
	Use:     "relations [modname]",
	Short:   "List a mod's relations; equivalent to 'packwiz relations list'",
	Aliases: []string{"relation"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: this. Any way to directly reference the list command from here
		// without creating a middle-man function?
	},
}

var relationListCmd = &cobra.Command{
	Use:     "list [modname]",
	Short:   "List a mod's relations",
	Aliases: []string{"ls"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: this.
	},
}

// display formatting functions
func displayRelations(modName string) {
	//Is there a way to display a relation tree?
}

// basic helper functions
func getAllRelations() {

}

// func getModRelations(modName string) core.ModRelations {
// 	pack, err := core.LoadPack()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	index, err := pack.LoadIndex()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	//get single mod Relations.
// 	modPath, ok := index.FindMod(modName)
// 	if !ok {
// 		fmt.Println("Can't find this file; please ensure you have run packwiz refresh and use the name of the .pw.toml file (defaults to the project slug)")
// 		os.Exit(1)
// 	}
// 	mod, err := core.LoadMod(modPath)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

//		return mod.Relations
//	}
func getDependencies() {}

func addRelation(modName string, relationSlug string) {

}

func removeRelation() {

}

func setRelations() {

}

func init() {
	rootCmd.AddCommand(relationsCmd)
	relationsCmd.AddCommand(relationListCmd)
}
