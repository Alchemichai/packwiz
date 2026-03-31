package core

import (
	"fmt"
	"os"
)

// func getAllTags(){

// }

func getModTags(modName string) *[]string {
	pack, err := LoadPack()
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
	mod, err := LoadMod(modPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &mod.Tags
}

// TODO: this
// func filterByTags(modList []*core.mod, tags []string) []core.mod{

// }
